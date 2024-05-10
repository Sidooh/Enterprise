package auth

import (
	"enterprise.sidooh/api/presenter"
	"enterprise.sidooh/pkg"
	"enterprise.sidooh/pkg/clients"
	"enterprise.sidooh/pkg/datastore"
	"enterprise.sidooh/pkg/entities"
	"enterprise.sidooh/pkg/services/enterprise"
	"enterprise.sidooh/pkg/services/user"
	"enterprise.sidooh/utils"
	"fmt"
	"github.com/Permify/permify-gorm/options"
	log "github.com/sirupsen/logrus"
	"time"
)

type Service interface {
	User(id int) (*presenter.UserWithRelations, error)
	Register(data presenter.Registration) (*presenter.EnterpriseWithUser, error)
	Verify(id int, phoneOtp int, emailOtp int) (*presenter.EnterpriseWithUser, error)
	Login(data presenter.Login) (*presenter.LoginResponse, error)

	GenerateOTP(id int, channel string) error
	ValidateOTP(id int, otp int) (*presenter.LoginResponse, error)
}

type service struct {
	accountsApi          *clients.ApiClient
	notifyApi            *clients.ApiClient
	paymentsApi          *clients.ApiClient
	enterpriseRepository enterprise.Repository
	userRepository       user.Repository
}

func (s *service) User(id int) (*presenter.UserWithRelations, error) {
	user, err := s.userRepository.ReadUserByIdWithEnterprise(id)
	if err != nil {
		log.Error(err)
		return nil, pkg.ErrUnauthorized
	}

	return getUserData(*user)
}

func (s *service) Register(data presenter.Registration) (*presenter.EnterpriseWithUser, error) {
	// Check enterprise does not exist, by email, phone or name
	_, err := s.enterpriseRepository.ReadEnterpriseByEmailOrPhone(data.Email, data.Phone)
	if err == nil {
		return nil, pkg.ErrInvalidEnterprise
	}

	// Check user does not exist by email
	_, err = s.userRepository.ReadUserByEmailOrPhone(data.Email, data.Phone)
	if err == nil {
		return nil, pkg.ErrInvalidUser
	}

	// 1. Create/get Sidooh account
	account, err := s.accountsApi.GetOrCreateAccount(data.Phone)
	if err != nil {
		return nil, pkg.ErrServerError
	}

	// 2. Create Enterprise
	enterprise, err := s.enterpriseRepository.CreateEnterprise(&entities.Enterprise{
		Name:      data.Name,
		Phone:     account.Phone,
		Email:     data.Email,
		AccountId: uint(account.Id),
	})
	if err != nil {
		return nil, err
	}

	password, _ := utils.HashPassword(data.Password)

	// 3. Create User
	user, err := s.userRepository.CreateUser(&entities.User{
		Phone:        account.Phone,
		Name:         data.AdminName,
		Email:        data.Email,
		Password:     password,
		EnterpriseId: enterprise.Id,
	})
	if err != nil {
		return nil, err
	}

	_ = datastore.Permify.AddRolesToUser(user.Id, "ADMIN")

	// 4. Create Float account
	updatedEnterprise := enterprise

	floatAccount, err := s.paymentsApi.CreateFloatAccount(int(enterprise.Id), account.Id)
	if err != nil {
		return nil, pkg.ErrServerError
	}
	updatedEnterprise, err = s.enterpriseRepository.UpdateEnterprise(enterprise, "float_account_id", floatAccount.Id)
	if err != nil {
		return nil, err
	}

	// send otps for verification
	go func() {
		_ = s.GenerateOTP(int(user.Id), "SMS")
		_ = s.GenerateOTP(int(user.Id), "MAIL")
	}()

	// return data
	return getEnterpriseData(entities.UserWithEnterprise{
		User:       *user,
		Enterprise: *updatedEnterprise,
	}), nil
}

func (s *service) Verify(id int, phoneOtp int, emailOtp int) (*presenter.EnterpriseWithUser, error) {
	user, err := s.userRepository.ReadUserByIdWithEnterprise(id)
	if err != nil {
		log.Error(err)
		return nil, pkg.ErrUnauthorized
	}

	if utils.CheckOTP(user.Phone, phoneOtp) && utils.CheckOTP(user.Email, emailOtp) {
		updateEnterprise, err := s.enterpriseRepository.UpdateEnterprise(&user.Enterprise, "phone_verified_at", time.Now())
		if err != nil {
			return nil, err
		}
		updateEnterprise, err = s.enterpriseRepository.UpdateEnterprise(&user.Enterprise, "email_verified_at", time.Now())
		if err != nil {
			return nil, err
		}

		user.Enterprise = *updateEnterprise

		return getEnterpriseData(*user), nil
	}

	return nil, pkg.ErrUnauthorized
}

func (s *service) Login(data presenter.Login) (*presenter.LoginResponse, error) {
	user, err := s.userRepository.ReadUserByEmailWithEnterprise(data.Email)
	if err != nil {
		log.Error(err)
		return nil, pkg.ErrUnauthorized
	}

	if user.Enterprise.PhoneVerifiedAt == nil || user.Enterprise.EmailVerifiedAt == nil {
		return nil, pkg.ErrInvalidEnterprise
	}

	res := utils.VerifyPassword(user.Password, data.Password)

	if !res {
		return nil, pkg.ErrUnauthorized
	}

	userData, err := getUserData(*user)

	response := &presenter.LoginResponse{User: userData}

	return response, err
}

func (s *service) GenerateOTP(id int, channel string) error {
	user, err := s.userRepository.ReadUser(id)
	if err != nil {
		log.Error(err)
		return pkg.ErrUnauthorized
	}

	otp := utils.RandomInt(6)

	// Send OTP to phone number
	message := fmt.Sprintf("S-%v is your verification code.", otp)
	destination := user.Phone
	f := s.notifyApi.SendSMS

	switch channel {
	case utils.MAIL:
		destination = user.Email
		f = s.notifyApi.SendMail
	}

	err = f("DEFAULT", destination, message)

	utils.SetOTP(destination, int(otp))

	return err
}

func (s *service) ValidateOTP(id int, otp int) (*presenter.LoginResponse, error) {
	user, err := s.userRepository.ReadUserByIdWithEnterprise(id)
	if err != nil {
		log.Error(err)
		return nil, pkg.ErrUnauthorized
	}

	if !utils.CheckOTP(user.Phone, otp) && !utils.CheckOTP(user.Email, otp) {
		return nil, pkg.ErrUnauthorized
	}

	userData, err := getUserData(*user)
	response := &presenter.LoginResponse{User: userData}

	return response, err
}

func getUserData(user entities.UserWithEnterprise) (*presenter.UserWithRelations, error) {
	roles, totalCount, err := datastore.Permify.GetRolesOfUser(user.Id, options.RoleOption{
		WithPermissions: true, // preload role's permissions
	})
	if err != nil {
		return nil, err
	}
	if totalCount == 0 {
		return nil, pkg.ErrUnauthorized
	}

	userData := &presenter.UserWithRelations{
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
		Enterprise: presenter.Enterprise{
			Id:      user.Enterprise.Id,
			Name:    user.Enterprise.Name,
			Country: user.Enterprise.Country,
			Address: user.Enterprise.Address,
			Phone:   user.Enterprise.Phone,
			Email:   user.Enterprise.Email,
		},
		Roles:       roles.Names(),
		Permissions: roles.Permissions().Names(),
	}

	return userData, nil
}

func getEnterpriseData(user entities.UserWithEnterprise) *presenter.EnterpriseWithUser {
	return &presenter.EnterpriseWithUser{
		Enterprise: presenter.Enterprise{
			Id:      user.Enterprise.Id,
			Name:    user.Enterprise.Name,
			Country: user.Enterprise.Country,
			Address: user.Enterprise.Address,
			Phone:   user.Enterprise.Phone,
			Email:   user.Enterprise.Email,
		},
		User: presenter.User{
			Id:    user.Id,
			Name:  user.Name,
			Email: user.Email,
		},
	}
}

func NewService(userRepository user.Repository) Service {
	enterpriseRepository := enterprise.NewRepo()

	return &service{
		enterpriseRepository: enterpriseRepository,
		userRepository:       userRepository,

		accountsApi: clients.GetAccountClient(),
		paymentsApi: clients.GetPaymentClient(),
		notifyApi:   clients.GetNotifyClient(),
	}
}
