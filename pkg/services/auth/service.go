package auth

import (
	"enterprise.sidooh/api/presenter"
	"github.com/gofiber/fiber/v2"
)

type Service interface {
	Register(data presenter.Registration) (*[]presenter.Enterprise, error)
	Login(data presenter.Login) (*presenter.Enterprise, error)
}

type service struct {
	apiClient  *fiber.Client
	repository Repository
}

func (s service) Register(data presenter.Registration) (*[]presenter.Enterprise, error) {
	//TODO implement me
	panic("implement me")
}

func (s service) Login(data presenter.Login) (*presenter.Enterprise, error) {
	//TODO implement me
	panic("implement me")
}

func NewService(r Repository) Service {
	return &service{repository: r}
}
