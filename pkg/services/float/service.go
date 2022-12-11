package float

import (
	"enterprise.sidooh/pkg"
	"enterprise.sidooh/pkg/clients"
	"enterprise.sidooh/pkg/entities"
	"enterprise.sidooh/pkg/services/enterprise"
)

type Service interface {
	GetFloatAccountForEnterprise(enterprise entities.Enterprise) (*clients.FloatAccount, error)
	GetFloatAccountTransactionsForEnterprise(enterprise entities.Enterprise) (*[]clients.FloatAccountTransaction, error)
	CreditFloatAccountForEnterprise(enterprise entities.Enterprise, amount, phone int) (*clients.FloatAccount, error)
}

type service struct {
	enterpriseRepository enterprise.Repository
	paymentsApi          *clients.ApiClient
}

func (s *service) GetFloatAccountForEnterprise(enterprise entities.Enterprise) (*clients.FloatAccount, error) {
	response, err := s.paymentsApi.FetchFloatAccount(int(enterprise.FloatAccountId))
	if err != nil {
		return nil, pkg.ErrServerError
	}

	return response, nil
}

func (s *service) GetFloatAccountTransactionsForEnterprise(enterprise entities.Enterprise) (*[]clients.FloatAccountTransaction, error) {
	response, err := s.paymentsApi.FetchFloatAccountTransactions(int(enterprise.FloatAccountId), 0)
	if err != nil {
		return nil, pkg.ErrServerError
	}

	return response, nil
}

func (s *service) CreditFloatAccountForEnterprise(enterprise entities.Enterprise, amount, phone int) (*clients.FloatAccount, error) {
	_, err := s.paymentsApi.CreditFloatAccount(int(enterprise.AccountId), int(enterprise.FloatAccountId), amount, phone)
	if err != nil {
		return nil, pkg.ErrServerError
	}

	return s.GetFloatAccountForEnterprise(enterprise)
}

func NewService(enterpriseRepository enterprise.Repository) Service {
	return &service{paymentsApi: clients.GetPaymentClient(), enterpriseRepository: enterpriseRepository}
}
