package auth

import (
	"enterprise.sidooh/api/presenter"
	"enterprise.sidooh/pkg"
	"enterprise.sidooh/pkg/cache"
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
	enterpriseExists, err := s.enterpriseRepository.ReadEnterpriseByEmailOrPhone(data.Email, data.Phone)
	if enterpriseExists != nil {
		return nil, pkg.ErrInvalidEnterprise
	}

	// Check user does not exist by email
	userExists, err := s.userRepository.ReadUserByEmailOrPhone(data.Email, data.Phone)
	if userExists.Id != 0 {
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
	// TODO: Refactor to payments client
	updatedEnterprise := enterprise

	floatAccount, err := s.paymentsApi.CreateFloatAccount(int(enterprise.Id), account.Id)
	if err != nil {
		return nil, pkg.ErrServerError
	}
	updatedEnterprise, err = s.enterpriseRepository.UpdateEnterprise(enterprise, "float_account_id", floatAccount.Id)
	if err != nil {
		return nil, err
	}

	// return data
	return &presenter.EnterpriseWithUser{
		Enterprise: presenter.Enterprise{
			Id:      updatedEnterprise.Id,
			Name:    updatedEnterprise.Name,
			Country: updatedEnterprise.Country,
			Address: updatedEnterprise.Address,
			Phone:   updatedEnterprise.Phone,
			Email:   updatedEnterprise.Email,
		},
		User: presenter.User{
			Id:    user.Id,
			Name:  user.Name,
			Email: user.Email,
		},
	}, nil
}

func (s *service) Login(data presenter.Login) (*presenter.LoginResponse, error) {
	user, err := s.userRepository.ReadUserByEmailWithEnterprise(data.Email)
	if err != nil {
		log.Error(err)
		return nil, pkg.ErrUnauthorized
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

	otp := utils.RandomInt(100000, 999999)

	// Send OTP to phone number
	message := fmt.Sprintf("S-%v is your verification code.", otp)
	switch channel {
	case utils.MAIL:
		err = s.notifyApi.SendMail("DEFAULT", user.Email, message)
	default:
		err = s.notifyApi.SendSMS("DEFAULT", user.Phone, message)
	}

	cache.Cache.Set(fmt.Sprintf("otp_%s", user.Phone), otp, 5*time.Minute)

	return err
}

func (s *service) ValidateOTP(id int, otp int) (*presenter.LoginResponse, error) {
	user, err := s.userRepository.ReadUserByIdWithEnterprise(id)
	if err != nil {
		log.Error(err)
		return nil, pkg.ErrUnauthorized
	}

	savedOtp := cache.Cache.Get(fmt.Sprintf("otp_%s", user.Phone))
	if savedOtp == nil || int((*savedOtp).(int64)) != otp {
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
