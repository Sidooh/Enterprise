package auth

import (
	"bytes"
	"encoding/json"
	"enterprise.sidooh/api/presenter"
	"enterprise.sidooh/pkg"
	"enterprise.sidooh/pkg/clients"
	"enterprise.sidooh/pkg/datastore"
	"enterprise.sidooh/pkg/entities"
	"enterprise.sidooh/pkg/services"
	"enterprise.sidooh/pkg/services/enterprise"
	"enterprise.sidooh/pkg/services/user"
	"enterprise.sidooh/utils"
	"github.com/Permify/permify-gorm/options"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type Service interface {
	User(id int) (*presenter.UserWithRelations, error)
	Register(data presenter.Registration) (*presenter.EnterpriseWithUser, error)
	Login(data presenter.Login) (*presenter.LoginResponse, error)
}

type service struct {
	accountsApi          *clients.ApiClient
	paymentsApi          *clients.ApiClient
	authRepository       Repository
	enterpriseRepository enterprise.Repository
	userRepository       user.Repository
}

func (s *service) User(id int) (*presenter.UserWithRelations, error) {
	user, err := s.authRepository.GetUserByIdWithEnterprise(id)
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
	if userExists != nil {
		return nil, pkg.ErrInvalidUser
	}

	// 1. Create/get Sidooh account
	var apiResponse = new(services.AccountApiResponse)

	jsonData, err := json.Marshal(map[string]string{"phone": data.Phone})
	dataBytes := bytes.NewBuffer(jsonData)

	err = s.accountsApi.NewRequest(http.MethodPost, "/accounts", dataBytes).Send(apiResponse)
	if err != nil {
		if err.Error() == "phone is already taken" {
			err = s.accountsApi.NewRequest(http.MethodGet, "/accounts/phone/"+data.Phone, nil).Send(apiResponse)
		} else {
			return nil, err
		}
	}

	account := apiResponse.Data

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

	var response = new(services.FloatAccountApiResponse)

	jsonData, err = json.Marshal(map[string]string{
		"initiator":  "ENTERPRISE",
		"reference":  strconv.Itoa(int(enterprise.Id)),
		"account_id": strconv.Itoa(account.Id),
	})
	dataBytes = bytes.NewBuffer(jsonData)

	err = s.paymentsApi.NewRequest(http.MethodPost, "/float-accounts", dataBytes).Send(response)
	if err == nil {
		floatAccount := response.Data

		updatedEnterprise, err = s.enterpriseRepository.UpdateEnterprise(enterprise, "float_account_id", floatAccount.Id)
		if err != nil {
			return nil, err
		}
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
	user, err := s.authRepository.GetUserByEmailWithEnterprise(data.Email)
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

func NewService(auth Repository, apiClient *clients.ApiClient) Service {
	enterpriseRepository := enterprise.NewRepo()
	userRepository := user.NewRepo()
	paymentsApi := clients.InitPaymentClient()

	return &service{
		authRepository:       auth,
		accountsApi:          apiClient,
		enterpriseRepository: enterpriseRepository,
		userRepository:       userRepository,
		paymentsApi:          paymentsApi,
	}
}
