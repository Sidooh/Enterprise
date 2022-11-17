package account

import (
	"enterprise.sidooh/api/presenter"
	"github.com/gofiber/fiber/v2"
)

type Service interface {
	FetchAccounts() (*[]presenter.Account, error)
	GetAccount(id int) (*presenter.Account, error)
	//CreateAccount(account *entities.Account) (*entities.Account, error)

	FetchAccountsForEnterprise(enterpriseId int) (*[]presenter.Account, error)
	GetAccountForEnterprise(enterpriseId int, id int) (*presenter.Account, error)
}

type service struct {
	apiClient  *fiber.Client
	repository Repository
}

func (s *service) FetchAccounts() (*[]presenter.Account, error) {
	return s.repository.ReadAccounts()
}

func (s *service) GetAccount(id int) (*presenter.Account, error) {
	return s.repository.ReadAccount(id)
}

//func (s *service) CreateAccount(account *entities.Account) (*entities.Account, error) {
//	return s.repository.CreateAccount(account)
//}

func (s *service) FetchAccountsForEnterprise(enterpriseId int) (*[]presenter.Account, error) {
	return s.repository.ReadEnterpriseAccounts(enterpriseId)
}

func (s *service) GetAccountForEnterprise(enterpriseId int, id int) (*presenter.Account, error) {
	return s.repository.ReadEnterpriseAccount(enterpriseId, id)
}

func NewService(r Repository) Service {
	return &service{repository: r}
}
