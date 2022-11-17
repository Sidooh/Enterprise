package account

import (
	"bytes"
	"encoding/json"
	"enterprise.sidooh/api/presenter"
	"enterprise.sidooh/pkg"
	"enterprise.sidooh/pkg/client"
	"enterprise.sidooh/pkg/entities"
	"enterprise.sidooh/pkg/services"
	"net/http"
)

type Service interface {
	FetchAccounts() (*[]presenter.Account, error)
	GetAccount(id int) (*presenter.Account, error)
	CreateAccount(account *entities.Account) (*presenter.Account, error)

	FetchAccountsForEnterprise(enterpriseId int) (*[]presenter.Account, error)
	GetAccountForEnterprise(enterpriseId int, id int) (*presenter.Account, error)
}

type service struct {
	accountsApi       *client.ApiClient
	accountRepository Repository
}

func (s *service) FetchAccounts() (*[]presenter.Account, error) {
	return s.accountRepository.ReadAccounts()
}

func (s *service) GetAccount(id int) (*presenter.Account, error) {
	return s.accountRepository.ReadAccount(id)
}

func (s *service) CreateAccount(account *entities.Account) (*presenter.Account, error) {
	accountExists, err := s.accountRepository.ReadEnterpriseAccountByPhone(int(account.EnterpriseId), account.Phone)
	if accountExists != nil {
		return nil, pkg.ErrInvalidAccount
	}

	// TODO: Refactor these common api calls
	var apiResponse = new(services.AccountApiResponse)

	jsonData, err := json.Marshal(map[string]string{"phone": account.Phone})
	dataBytes := bytes.NewBuffer(jsonData)

	err = s.accountsApi.NewRequest(http.MethodPost, "/accounts", dataBytes).Send(apiResponse)
	if err != nil {
		if err.Error() == "phone is already taken" {
			err = s.accountsApi.NewRequest(http.MethodGet, "/accounts/phone/"+account.Phone, nil).Send(apiResponse)
		} else {
			return nil, err
		}
	}

	response := apiResponse.Data

	account.AccountId = uint(response.Id)
	account.Phone = response.Phone

	model, err := s.accountRepository.CreateAccount(account)

	return &presenter.Account{
		Id:    model.Id,
		Phone: model.Phone,
		Name:  model.Name,
	}, err
}

func (s *service) FetchAccountsForEnterprise(enterpriseId int) (*[]presenter.Account, error) {
	return s.accountRepository.ReadEnterpriseAccounts(enterpriseId)
}

func (s *service) GetAccountForEnterprise(enterpriseId int, id int) (*presenter.Account, error) {
	return s.accountRepository.ReadEnterpriseAccount(enterpriseId, id)
}

func NewService(account Repository) Service {
	accountsApi := client.InitAccountClient()
	return &service{accountRepository: account, accountsApi: accountsApi}
}
