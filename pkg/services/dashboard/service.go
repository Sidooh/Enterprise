package dashboard

import (
	"enterprise.sidooh/pkg"
	"enterprise.sidooh/pkg/clients"
	"enterprise.sidooh/pkg/entities"
	"enterprise.sidooh/utils"
)

type Service interface {
	GetRecentVoucherTransactionsForEnterprise(accountId int, limit int) (*[]clients.VoucherTransaction, error)
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

	voucherTransactions, err := s.paymentsApi.FetchVoucherTransactions(int(enterprise.AccountId), utils.VOUCHER_TRANSACTIONS_LIMIT)
	if err != nil {
		return nil, pkg.ErrServerError
	}

	totalCount, totalAmount := 0, 0
	for _, vT := range *voucherTransactions {
		if vT.Type == "CREDIT" {
			totalCount++
			totalAmount += vT.Amount
		}
	}

	//	TODO: Implement vouchers disbursed logic
	result := &clients.DashboardStatistics{
		FloatBalance:            response.Balance,
		AccountsCount:           int(s.repository.CountAccounts()),
		DisbursedVouchersAmount: totalAmount,
		DisbursedVouchersCount:  totalCount,
	}

	return result, nil
}

func (s *service) GetRecentVoucherTransactionsForEnterprise(accountId int, limit int) (*[]clients.VoucherTransaction, error) {
	response, err := s.paymentsApi.FetchVoucherTransactions(accountId, limit)
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
