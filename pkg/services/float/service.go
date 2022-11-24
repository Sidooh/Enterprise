package float

import (
	"enterprise.sidooh/pkg"
	"enterprise.sidooh/pkg/clients"
	"enterprise.sidooh/pkg/entities"
	"enterprise.sidooh/pkg/services/enterprise"
)

type Service interface {
	GetFloatAccountForEnterprise(enterprise entities.Enterprise) (*clients.FloatAccount, error)
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

func NewService(enterpriseRepository enterprise.Repository) Service {
	return &service{paymentsApi: clients.GetPaymentClient(), enterpriseRepository: enterpriseRepository}
}
