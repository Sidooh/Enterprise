package auth

import (
	jwt2 "enterprise.sidooh/api/middleware/jwt"
	"enterprise.sidooh/api/presenter"
	"enterprise.sidooh/pkg/client"
	"enterprise.sidooh/pkg/datastore"
	"enterprise.sidooh/utils"
	"errors"
	"github.com/Permify/permify-gorm/options"
	"github.com/golang-jwt/jwt/v4"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"time"
)

type AuthResponse struct {
	Token string `json:"access_token"`
}

type Service interface {
	User(id int) (*presenter.User, error)
	Register(data presenter.Registration) (*[]presenter.Enterprise, error)
	Login(data presenter.Login) (*presenter.LoginResponse, error)
}

type service struct {
	api        *client.ApiClient
	repository Repository
}

func (s *service) User(id int) (*presenter.User, error) {
	user, err := s.repository.GetUserByIdWithEnterprise(id)
	if err != nil {
		log.Error(err)
		return nil, errors.New("unauthorized")
	}

	roles, totalCount, err := datastore.Permify.GetRolesOfUser(user.Id, options.RoleOption{
		WithPermissions: true, // preload role's permissions
	})
	if err != nil {
		return nil, err
	}
	if totalCount == 0 {
		return nil, errors.New("unauthorized")
	}

	userData := &presenter.User{
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

	return userData, err
}

func (s *service) Register(data presenter.Registration) (*[]presenter.Enterprise, error) {
	//TODO implement me
	panic("implement me")

	// 1. Create/get Sidooh account

	// 2. Create Enterprise

	// 3. Create User

	// 4. Create Float account

}

func (s *service) Login(data presenter.Login) (*presenter.LoginResponse, error) {
	user, err := s.repository.GetUserByEmailWithEnterprise(data.Email)
	if err != nil {
		log.Error(err)
		return nil, errors.New("invalid credentials")
	}

	res := utils.Compare(user.Password, data.Password)

	if !res {
		return nil, errors.New("invalid credentials")
	}

	roles, totalCount, err := datastore.Permify.GetRolesOfUser(user.Id, options.RoleOption{
		WithPermissions: true, // preload role's permissions
	})
	if err != nil {
		return nil, err
	}
	if totalCount == 0 {
		return nil, errors.New("invalid credentials")
	}

	validity := time.Duration(viper.GetInt("ACCESS_TOKEN_VALIDITY")) * time.Minute
	token, err := jwt2.Encode(&jwt.MapClaims{
		"name":  user.Name,
		"email": user.Email,
		"id":    user.Id,
	}, validity)

	userData := &presenter.User{
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

	response := &presenter.LoginResponse{Token: token, User: userData}

	return response, err
}

func NewService(r Repository, apiClient *client.ApiClient) Service {
	return &service{repository: r, api: apiClient}
}
