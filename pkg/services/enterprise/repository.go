package enterprise

import (
	"enterprise.sidooh/api/presenter"
	"enterprise.sidooh/pkg/datastore"
)

//Repository interface allows us to access the CRUD Operations here.
type Repository interface {
	//CreateEnterprise(enterprise *entities.Enterprise) (*entities.Enterprise, error)
	ReadEnterprises() (*[]presenter.Enterprise, error)
	ReadEnterprise(id int) (*presenter.Enterprise, error)
	//UpdateEnterprise(enterprise *entities.Enterprise) (*entities.Enterprise, error)
	//DeleteEnterprise(Id uint) error
}
type repository struct {
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

//NewRepo is the single instance repo that is being created.
func NewRepo() Repository {
	return &repository{}
}
