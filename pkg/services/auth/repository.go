package auth

import (
	"enterprise.sidooh/pkg/datastore"
	"enterprise.sidooh/pkg/entities"
)

type Repository interface {
	GetAccountByEmailWithEnterprise(email string) (*entities.AccountWithEnterprise, error)
}
type repository struct {
}

func (r *repository) GetAccountByEmailWithEnterprise(email string) (*entities.AccountWithEnterprise, error) {
	var account entities.AccountWithEnterprise
	result := datastore.DB.Where("accounts.email = ?", email).Joins("Enterprise").First(&account)
	if result.Error != nil {
		return nil, result.Error
	}

	return &account, nil
}

func NewRepo() Repository {
	return &repository{}
}
