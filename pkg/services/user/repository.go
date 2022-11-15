package user

import (
	"enterprise.sidooh/api/presenter"
	"enterprise.sidooh/pkg/datastore"
)

//Repository interface allows us to access the CRUD Operations here.
type Repository interface {
	ReadUser(id int) (*presenter.User, error)
}
type repository struct {
}

func (r *repository) ReadUser(id int) (*presenter.User, error) {
	var user presenter.User
	result := datastore.DB.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

//NewRepo is the single instance repo that is being created.
func NewRepo() Repository {
	return &repository{}
}
