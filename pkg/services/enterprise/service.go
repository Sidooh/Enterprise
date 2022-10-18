package enterprise

import (
	"enterprise.sidooh/api/presenter"
	"github.com/gofiber/fiber/v2"
)

type Service interface {
	FetchEnterprises() (*[]presenter.Enterprise, error)
	GetEnterprise(id int) (*presenter.Enterprise, error)
}

type service struct {
	apiClient  *fiber.Client
	repository Repository
}

func (s *service) FetchEnterprises() (*[]presenter.Enterprise, error) {
	return s.repository.ReadEnterprises()
}

func (s *service) GetEnterprise(id int) (*presenter.Enterprise, error) {
	return s.repository.ReadEnterprise(id)
}

func NewService(r Repository) Service {
	return &service{repository: r}
}
