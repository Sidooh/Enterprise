package auth

import "enterprise.sidooh/api/presenter"

type Repository interface {
	Register() (*presenter.Account, error)
	Login() (*presenter.Account, error)
}
type repository struct {
}

func (r repository) Register() (*presenter.Account, error) {
	//TODO implement me
	panic("implement me")
}

func (r repository) Login() (*presenter.Account, error) {
	//TODO implement me
	panic("implement me")
}

func NewRepo() Repository {
	return &repository{}
}
