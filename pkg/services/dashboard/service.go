package dashboard

import (
	"enterprise.sidooh/pkg"
	"enterprise.sidooh/pkg/clients"
	"enterprise.sidooh/pkg/entities"
)

type Service interface {
	GetRecentVoucherTransactionsForEnterprise(enterprise entities.Enterprise, limit int) (*[]clients.VoucherTransaction, error)
	GetRecentFloatAccountTransactionsForEnterprise(enterprise entities.Enterprise, limit int) (*[]clients.FloatAccountTransaction, error)
	GetDashboardStatistics(enterprise entities.Enterprise) (*clients.DashboardStatistics, error)
}

type service struct {
	paymentsApi *clients.ApiClient
	repository  Repository
}

func (s *service) GetDashboardStatistics(enterprise entities.Enterprise) (*clients.DashboardStatistics, error) {
	response, err := s.paymentsApi.FetchFloatAccount(int(enterprise.FloatAccountId))
	if err != nil {
		return nil, pkg.ErrServerError
	}

	result := &clients.DashboardStatistics{
		FloatBalance:      response.Balance,
		AccountsCount:     int(s.repository.CountAccounts()),
		VouchersDisbursed: 21,
	}

	return result, nil
}

func (s *service) GetRecentVoucherTransactionsForEnterprise(enterprise entities.Enterprise, limit int) (*[]clients.VoucherTransaction, error) {
	response, err := s.paymentsApi.FetchVoucherTransactions(int(enterprise.AccountId), limit)
	if err != nil {
		return nil, pkg.ErrServerError
	}

	return response, nil
}

func (s *service) GetRecentFloatAccountTransactionsForEnterprise(enterprise entities.Enterprise, limit int) (*[]clients.FloatAccountTransaction, error) {
	response, err := s.paymentsApi.FetchFloatAccountTransactions(int(enterprise.FloatAccountId), limit)
	if err != nil {
		return nil, pkg.ErrServerError
	}

	return response, nil
}

func NewService(r Repository) Service {
	return &service{paymentsApi: clients.GetPaymentClient(), repository: r}
}
