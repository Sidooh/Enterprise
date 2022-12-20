package account

import (
	"enterprise.sidooh/api/presenter"
	"enterprise.sidooh/pkg"
	"enterprise.sidooh/pkg/clients"
	"enterprise.sidooh/pkg/entities"
	"golang.org/x/exp/slices"
)

type Service interface {
	FetchAccounts() (*[]entities.Account, error)
	GetAccount(id int) (*presenter.Account, error)
	CreateAccount(account *entities.Account) (*entities.Account, error)

	FetchAccountsForEnterprise(enterpriseId int) (*[]entities.Account, error)
	GetAccountForEnterprise(enterpriseId int, id int) (*presenter.Account, error)
	CreateBulkAccounts(accounts []entities.Account) (*[]entities.Account, map[string]string)
}

type service struct {
	accountsApi       *clients.ApiClient
	paymentsApi       *clients.ApiClient
	accountRepository Repository
}

func (s *service) FetchAccounts() (*[]entities.Account, error) {
	return s.accountRepository.ReadAccounts()
}

func (s *service) GetAccount(id int) (*presenter.Account, error) {
	account, err := s.accountRepository.ReadAccount(id)

	response, err := s.paymentsApi.FetchVouchers(int(account.AccountId))
	if err != nil {
		return nil, pkg.ErrServerError
	}

	return &presenter.Account{
		Id:           account.Id,
		EnterpriseId: account.EnterpriseId,
		Name:         account.Name,
		Phone:        account.Phone,
		Vouchers:     response,
	}, err
}

func (s *service) CreateAccount(account *entities.Account) (*entities.Account, error) {
	_, err := s.accountRepository.ReadEnterpriseAccountByPhone(int(account.EnterpriseId), account.Phone)
	if err == nil {
		return nil, pkg.ErrInvalidAccount
	}

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

func (s *service) GetAccountForEnterprise(enterpriseId int, id int) (*presenter.Account, error) {
	account, err := s.accountRepository.ReadEnterpriseAccount(enterpriseId, id)
	if err != nil {
		return nil, pkg.ErrServerError
	}

	response, err := s.paymentsApi.FetchVouchers(int(account.AccountId))
	if err != nil {
		return nil, pkg.ErrServerError
	}

	return &presenter.Account{
		Id:           account.Id,
		EnterpriseId: account.EnterpriseId,
		Name:         account.Name,
		Phone:        account.Phone,
		Teams:        account.Teams,
		Vouchers:     response,
	}, err
}

func (s *service) CreateBulkAccounts(accounts []entities.Account) (*[]entities.Account, map[string]string) {
	var phones []string
	for _, account := range accounts {
		phones = append(phones, account.Phone)
	}
	existingAccounts, err := s.accountRepository.ReadEnterpriseAccountsByPhone(int(accounts[0].EnterpriseId), phones)
	if err != nil {
		panic("something went wrong")
	}

	var existingPhones []string
	for _, account := range *existingAccounts {
		if slices.Contains(phones, account.Phone) {
			existingPhones = append(existingPhones, account.Phone)
		}
	}

	var newAccounts []entities.Account
	exceptions := map[string]string{}

	for _, account := range accounts {
		if slices.Contains(existingPhones, account.Phone) {
			exceptions[account.Phone] = "account already exists"
			continue
		}

		response, err := s.accountsApi.GetOrCreateAccount(account.Phone)
		if err != nil {
			exceptions[account.Phone] = "could not create account, try again"
			continue
		}

		account.AccountId = uint(response.Id)
		account.Phone = response.Phone

		a, err := s.accountRepository.CreateAccount(&account)
		if err != nil {
			exceptions[account.Phone] = err.Error()
			continue
		} else {
			newAccounts = append(newAccounts, *a)
		}
	}

	return &newAccounts, exceptions
}

func NewService(account Repository) Service {
	return &service{
		accountRepository: account,
		accountsApi:       clients.GetAccountClient(),
		paymentsApi:       clients.GetPaymentClient(),
	}
}
