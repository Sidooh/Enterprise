package user

import (
	"enterprise.sidooh/pkg/datastore"
	"enterprise.sidooh/pkg/entities"
	"fmt"
)

//Repository interface allows us to access the CRUD Operations here.
type Repository interface {
	CreateUser(user *entities.User) (*entities.User, error)
	ReadUser(id int) (*entities.User, error)
	ReadUserByEmailOrPhone(email string, phone string) (*entities.User, error)

	ReadUserByEmailWithEnterprise(email string) (*entities.UserWithEnterprise, error)
	ReadUserByIdWithEnterprise(id int) (*entities.UserWithEnterprise, error)
}
type repository struct {
}

func (r *repository) CreateUser(user *entities.User) (u *entities.User, err error) {
	err = datastore.DB.Create(&user).Error
	return user, err
}

func (r *repository) ReadUser(id int) (user *entities.User, err error) {
	err = datastore.DB.First(&user, id).Error
	return
}

func (r *repository) ReadUserByEmailOrPhone(email string, phone string) (user *entities.User, err error) {
	err = datastore.DB.Where("email", email).Or("phone LIKE ?", fmt.Sprintf("%%%s%%", phone)).
		First(&user).Error
	return
}

func (r *repository) ReadUserByEmailWithEnterprise(email string) (user *entities.UserWithEnterprise, err error) {
	err = datastore.DB.Where("users.email = ?", email).Joins("Enterprise").First(&user).Error
	return
}

func (r *repository) ReadUserByIdWithEnterprise(id int) (user *entities.UserWithEnterprise, err error) {
	err = datastore.DB.Where("users.id = ?", id).Joins("Enterprise").First(&user).Error
	return
}

//NewRepo is the single instance repo that is being created.
func NewRepo() Repository {
	return &repository{}
}
