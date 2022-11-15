package user

import (
	"enterprise.sidooh/api/presenter"
	"github.com/gofiber/fiber/v2"
)

type Service interface {
	GetUser(id int) (*presenter.User, error)
}

type service struct {
	apiClient  *fiber.Client
	repository Repository
}

func (s *service) GetUser(id int) (*presenter.User, error) {
	return s.repository.ReadUser(id)
}

func NewService(r Repository) Service {
	return &service{repository: r}
}
