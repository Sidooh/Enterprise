package account

import (
	"enterprise.sidooh/api/presenter"
	"enterprise.sidooh/pkg/datastore"
	"enterprise.sidooh/pkg/entities"
	"fmt"
)

//Repository interface allows us to access the CRUD Operations here.
type Repository interface {
	//CreateAccount(account *entities.Account) (*entities.Account, error)
	ReadAccounts() (*[]presenter.Account, error)
	ReadAccount(id int) (*presenter.Account, error)
	//ReadAccountByEmailOrPhone(email string, phone string) (*presenter.Account, error)
	//UpdateAccount(account *entities.Account, column string, value interface{}) (*entities.Account, error)
	//DeleteAccount(Id uint) error

	ReadEnterpriseAccounts(enterpriseId int) (*[]presenter.Account, error)
	ReadEnterpriseAccount(enterpriseId int, id int) (*presenter.Account, error)
}
type repository struct {
}

func (r *repository) CreateAccount(account *entities.Account) (*entities.Account, error) {
	result := datastore.DB.Create(&account)
	if result.Error != nil {
		return nil, result.Error
	}

	return account, nil
}

func (r *repository) ReadAccounts() (*[]presenter.Account, error) {
	var accounts []presenter.Account
	result := datastore.DB.Find(&accounts)
	if result.Error != nil {
		return nil, result.Error
	}

	return &accounts, nil
}

func (r *repository) ReadAccount(id int) (*presenter.Account, error) {
	var account presenter.Account
	result := datastore.DB.First(&account, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &account, nil
}

func (r *repository) ReadAccountByEmailOrPhone(email string, phone string) (*presenter.Account, error) {
	var account presenter.Account
	result := datastore.DB.Where("email", email).Or("phone LIKE ?", fmt.Sprintf("%%%s%%", phone)).First(&account)
	if result.Error != nil {
		return nil, result.Error
	}

	return &account, nil
}

func (r *repository) UpdateAccount(account *entities.Account, column string, value interface{}) (*entities.Account, error) {
	result := datastore.DB.Model(&account).Update(column, value)
	if result.Error != nil {
		return nil, result.Error
	}

	return account, nil
}

//======================================================================================================================
//	Enterprise Limited Data Manipulation
//======================================================================================================================

func (r *repository) ReadEnterpriseAccounts(enterpriseId int) (*[]presenter.Account, error) {
	var accounts []presenter.Account
	result := datastore.DB.Where("enterprise_id", enterpriseId).Find(&accounts)
	if result.Error != nil {
		return nil, result.Error
	}

	return &accounts, nil
}

func (r *repository) ReadEnterpriseAccount(enterpriseId int, id int) (*presenter.Account, error) {
	var account presenter.Account
	result := datastore.DB.Where("enterprise_id", enterpriseId).First(&account, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &account, nil
}

//NewRepo is the single instance repo that is being created.
func NewRepo() Repository {
	return &repository{}
}
