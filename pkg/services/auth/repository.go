package auth

import (
	"enterprise.sidooh/pkg/datastore"
	"enterprise.sidooh/pkg/entities"
)

type Repository interface {
	GetUserByEmailWithEnterprise(email string) (*entities.UserWithEnterprise, error)
	GetUserByIdWithEnterprise(id int) (*entities.UserWithEnterprise, error)
}
type repository struct {
}

func (r *repository) GetUserByEmailWithEnterprise(email string) (*entities.UserWithEnterprise, error) {
	var user entities.UserWithEnterprise
	result := datastore.DB.Where("users.email = ?", email).Joins("Enterprise").First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (r *repository) GetUserByIdWithEnterprise(id int) (*entities.UserWithEnterprise, error) {
	var user entities.UserWithEnterprise
	result := datastore.DB.Where("users.id = ?", id).Joins("Enterprise").First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func NewRepo() Repository {
	return &repository{}
}
