package account

import (
	"enterprise.sidooh/pkg"
	"enterprise.sidooh/pkg/clients"
	"enterprise.sidooh/pkg/entities"
)

type Service interface {
	FetchAccounts() (*[]entities.Account, error)
	GetAccount(id int) (*entities.Account, error)
	CreateAccount(account *entities.Account) (*entities.Account, error)

	FetchAccountsForEnterprise(enterpriseId int) (*[]entities.Account, error)
	GetAccountForEnterprise(enterpriseId int, id int) (*entities.Account, error)
}

type service struct {
	accountsApi       *clients.ApiClient
	accountRepository Repository
}

func (s *service) FetchAccounts() (*[]entities.Account, error) {
	return s.accountRepository.ReadAccounts()
}

func (s *service) GetAccount(id int) (*entities.Account, error) {
	return s.accountRepository.ReadAccount(id)
}

func (s *service) CreateAccount(account *entities.Account) (*entities.Account, error) {
	_, err := s.accountRepository.ReadEnterpriseAccountByPhone(int(account.EnterpriseId), account.Phone)
	if err == nil {
		return nil, pkg.ErrInvalidAccount
	}

	// TODO: Refactor these common api calls
	response, err := s.accountsApi.GetOrCreateAccount(account.Phone)
	if err != nil {
		return nil, pkg.ErrServerError
	}

	account.AccountId = uint(response.Id)
	account.Phone = response.Phone

	return s.accountRepository.CreateAccount(account)
}

func (s *service) FetchAccountsForEnterprise(enterpriseId int) (*[]entities.Account, error) {
	return s.accountRepository.ReadEnterpriseAccounts(enterpriseId)
}

func (s *service) GetAccountForEnterprise(enterpriseId int, id int) (*entities.Account, error) {
	return s.accountRepository.ReadEnterpriseAccount(enterpriseId, id)
}

func NewService(account Repository) Service {
	accountsApi := clients.InitAccountClient()
	return &service{accountRepository: account, accountsApi: accountsApi}
}
