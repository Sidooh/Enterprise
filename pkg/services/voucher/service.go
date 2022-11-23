package voucher

import (
	"enterprise.sidooh/pkg"
	"enterprise.sidooh/pkg/clients"
)

type Service interface {
	//FetchVoucherTypes() (*[]entities.Account, error)
	//GetVoucherType(id int) (*entities.Account, error)
	CreateVoucherType(enterpriseId int, name string) (*clients.VoucherType, error)

	FetchVoucherTypesForEnterprise(enterpriseId int) (*[]clients.VoucherType, error)
	GetVoucherTypeForEnterprise(enterpriseId int, id int) (*clients.VoucherType, error)
}

type service struct {
	paymentsApi *clients.ApiClient
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

func NewService() Service {
	return &service{paymentsApi: clients.GetPaymentClient()}
}
