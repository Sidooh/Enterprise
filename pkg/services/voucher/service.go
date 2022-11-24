package voucher

import (
	"enterprise.sidooh/pkg"
	"enterprise.sidooh/pkg/clients"
	"enterprise.sidooh/pkg/entities"
	"enterprise.sidooh/pkg/services/account"
	"fmt"
	"golang.org/x/exp/slices"
)

type Service interface {
	// TODO: Naming convention, determine which to use

	CreateVoucherType(enterpriseId int, name string) (*clients.VoucherType, error)
	FetchVoucherTypesForEnterprise(enterpriseId int) (*[]clients.VoucherType, error)
	GetVoucherTypeForEnterprise(enterpriseId, id int) (*clients.VoucherType, error)
	DisburseVoucherType(enterprise entities.Enterprise, voucherTypeId, accountId, amount int) (*clients.VoucherType, error)
}

type service struct {
	accountRepository account.Repository
	paymentsApi       *clients.ApiClient
}

func (s *service) CreateVoucherType(enterpriseId int, name string) (*clients.VoucherType, error) {
	response, err := s.paymentsApi.CreateVoucherType(enterpriseId, name)
	if err != nil {
		return nil, pkg.ErrServerError
	}

	return response, nil
}

func (s *service) FetchVoucherTypesForEnterprise(enterpriseId int) (*[]clients.VoucherType, error) {
	response, err := s.paymentsApi.FetchVoucherTypes(enterpriseId)
	if err != nil {
		return nil, pkg.ErrServerError
	}

	return response, nil
}

func (s *service) GetVoucherTypeForEnterprise(enterpriseId int, id int) (*clients.VoucherType, error) {
	response, err := s.paymentsApi.FetchVoucherType(enterpriseId, id)
	if err != nil {
		return nil, pkg.ErrServerError
	}

	return response, nil
}

func (s *service) DisburseVoucherType(enterprise entities.Enterprise, voucherTypeId, accountId, amount int) (*clients.VoucherType, error) {
	account, err := s.accountRepository.ReadAccount(accountId)
	if err != nil {
		return nil, pkg.ErrInvalidAccount
	}

	voucherType, err := s.GetVoucherTypeForEnterprise(int(enterprise.Id), voucherTypeId)
	if err != nil {
		return nil, pkg.ErrServerError
	}

	index := slices.IndexFunc(voucherType.Vouchers, func(v clients.Voucher) bool {
		return v.AccountId == int(account.AccountId)
	})

	voucherId := 0

	fmt.Println(index)
	if index < 0 {
		voucher, err := s.paymentsApi.CreateVoucher(int(enterprise.AccountId), accountId, voucherTypeId)
		if err != nil {
			return nil, err
		}
		voucherId = voucher.Id
	} else {
		voucherId = voucherType.Vouchers[index].Id
	}

	response, err := s.paymentsApi.DisburseVoucher(int(enterprise.Id), int(enterprise.FloatAccountId), voucherId, amount)
	if err != nil {
		return nil, pkg.ErrServerError
	}

	return response, err
}

func NewService() Service {
	accountRepository := account.NewRepo()
	return &service{paymentsApi: clients.GetPaymentClient(), accountRepository: accountRepository}
}
