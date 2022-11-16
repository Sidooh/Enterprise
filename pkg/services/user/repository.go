package user

import (
	"enterprise.sidooh/api/presenter"
	"enterprise.sidooh/pkg/datastore"
	"enterprise.sidooh/pkg/entities"
	"fmt"
)

//Repository interface allows us to access the CRUD Operations here.
type Repository interface {
	CreateUser(user *entities.User) (*entities.User, error)
	ReadUser(id int) (*presenter.User, error)
	ReadUserByEmailOrPhone(email string, phone string) (*presenter.User, error)
}
type repository struct {
}

func (r *repository) CreateUser(user *entities.User) (*entities.User, error) {
	result := datastore.DB.Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (r *repository) ReadUser(id int) (*presenter.User, error) {
	var user presenter.User
	result := datastore.DB.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (r *repository) ReadUserByEmailOrPhone(email string, phone string) (*presenter.User, error) {
	var user presenter.User
	result := datastore.DB.Where("email", email).Or("phone LIKE ?", fmt.Sprintf("%%%s%%", phone)).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

//NewRepo is the single instance repo that is being created.
func NewRepo() Repository {
	return &repository{}
}
