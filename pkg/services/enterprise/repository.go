package enterprise

import (
	"enterprise.sidooh/api/presenter"
	"enterprise.sidooh/pkg/datastore"
	"enterprise.sidooh/pkg/entities"
	"fmt"
)

//Repository interface allows us to access the CRUD Operations here.
type Repository interface {
	CreateEnterprise(enterprise *entities.Enterprise) (*entities.Enterprise, error)
	ReadEnterprises() (*[]presenter.Enterprise, error)
	ReadEnterprise(id int) (*presenter.Enterprise, error)
	ReadEnterpriseByEmailOrPhone(email string, phone string) (*presenter.Enterprise, error)
	UpdateEnterprise(enterprise *entities.Enterprise, column string, value interface{}) (*entities.Enterprise, error)
}
type repository struct {
}

func (r *repository) CreateEnterprise(enterprise *entities.Enterprise) (*entities.Enterprise, error) {
	result := datastore.DB.Create(&enterprise)
	if result.Error != nil {
		return nil, result.Error
	}

	return enterprise, nil
}

func (r *repository) ReadEnterprises() (*[]presenter.Enterprise, error) {
	var enterprises []presenter.Enterprise
	result := datastore.DB.Find(&enterprises)
	if result.Error != nil {
		return nil, result.Error
	}

	return &enterprises, nil
}

func (r *repository) ReadEnterprise(id int) (*presenter.Enterprise, error) {
	var enterprise presenter.Enterprise
	result := datastore.DB.First(&enterprise, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &enterprise, nil
}

func (r *repository) ReadEnterpriseByEmailOrPhone(email string, phone string) (*presenter.Enterprise, error) {
	var enterprise presenter.Enterprise
	result := datastore.DB.Where("email", email).Or("phone LIKE ?", fmt.Sprintf("%%%s%%", phone)).First(&enterprise)
	if result.Error != nil {
		return nil, result.Error
	}

	return &enterprise, nil
}

func (r *repository) UpdateEnterprise(enterprise *entities.Enterprise, column string, value interface{}) (*entities.Enterprise, error) {
	result := datastore.DB.Model(&enterprise).Update(column, value)
	if result.Error != nil {
		return nil, result.Error
	}

	return enterprise, nil
}

//NewRepo is the single instance repo that is being created.
func NewRepo() Repository {
	return &repository{}
}
