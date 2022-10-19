package auth

import (
	"bytes"
	"encoding/json"
	"enterprise.sidooh/api/presenter"
	"enterprise.sidooh/pkg/client"
	"enterprise.sidooh/pkg/datastore"
	"errors"
	"github.com/Permify/permify-gorm/options"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type AuthResponse struct {
	Token string `json:"access_token"`
}

type Service interface {
	Register(data presenter.Registration) (*[]presenter.Enterprise, error)
	Login(data presenter.Login) (*presenter.LoginResponse, error)
}

type service struct {
	api        *client.ApiClient
	repository Repository
}

func (s *service) Register(data presenter.Registration) (*[]presenter.Enterprise, error) {
	//TODO implement me
	panic("implement me")
}

func (s *service) Login(data presenter.Login) (*presenter.LoginResponse, error) {
	// TODO: Check the email exists in Ent Accs, fetch with Enterprise to make it more efficient
	account, err := s.repository.GetAccountByEmailWithEnterprise(data.Email)
	if err != nil {
		log.Error(err)
		return nil, errors.New("invalid credentials")
	}

	jsonData, err := json.Marshal(data)
	dataBytes := bytes.NewBuffer(jsonData)

	var authResponse = new(AuthResponse)

	err = s.api.NewRequest(http.MethodPost, "/users/signin", dataBytes).Send(authResponse)
	if err != nil {
		return nil, err
	}

	roles, totalCount, err := datastore.Permify.GetRolesOfUser(account.Id, options.RoleOption{
		WithPermissions: true, // preload role's permissions
	})
	if err != nil {
		return nil, err
	}
	if totalCount == 0 {
		return nil, errors.New("invalid credentials")
	}

	accountData := &presenter.Account{
		Id:    account.Id,
		Phone: account.Phone,
		Name:  account.Name,
		Email: account.Email,
		Enterprise: presenter.Enterprise{
			Id:      account.Enterprise.Id,
			Name:    account.Enterprise.Name,
			Country: account.Enterprise.Country,
			Address: account.Enterprise.Address,
			Phone:   account.Enterprise.Phone,
			Email:   account.Enterprise.Email,
		},
		Roles:       roles.Names(),
		Permissions: roles.Permissions().Names(),
	}

	response := &presenter.LoginResponse{Token: authResponse.Token, Admin: accountData}

	return response, err
}

func NewService(r Repository, apiClient *client.ApiClient) Service {
	return &service{repository: r, api: apiClient}
}
