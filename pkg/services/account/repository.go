package account

import (
	"enterprise.sidooh/pkg/datastore"
	"enterprise.sidooh/pkg/entities"
	"fmt"
)

// Repository interface allows us to access the CRUD Operations here.
type Repository interface {
	ReadAccounts() (*[]entities.Account, error)
	ReadAccount(id int) (*entities.Account, error)
	CreateAccount(account *entities.Account) (*entities.Account, error)

	ReadEnterpriseAccounts(enterpriseId int) (*[]entities.Account, error)
	ReadEnterpriseAccount(enterpriseId int, id int) (*entities.Account, error)
	ReadEnterpriseAccountByPhone(enterpriseId int, phone string) (*entities.Account, error)
	ReadEnterpriseAccountsByPhone(enterpriseId int, phones []string) (*[]entities.Account, error)
}
type repository struct {
}

func (r *repository) CreateAccount(account *entities.Account) (a *entities.Account, err error) {
	err = datastore.DB.Create(&account).Error
	return account, err
}

func (r *repository) ReadAccounts() (accounts *[]entities.Account, err error) {
	err = datastore.DB.Find(&accounts).Error
	return
}

func (r *repository) ReadAccount(id int) (account *entities.Account, err error) {
	err = datastore.DB.Preload("Teams").First(&account, id).Error
	return
}

func (r *repository) UpdateAccount(account *entities.Account, column string, value interface{}) (a *entities.Account, err error) {
	err = datastore.DB.Model(&account).Update(column, value).Error
	return account, err
}

//======================================================================================================================
//	Enterprise Limited Data Manipulation
//======================================================================================================================

func (r *repository) ReadEnterpriseAccounts(enterpriseId int) (accounts *[]entities.Account, err error) {
	err = datastore.DB.Where("enterprise_id", enterpriseId).Find(&accounts).Error
	return
}

func (r *repository) ReadEnterpriseAccount(enterpriseId int, id int) (account *entities.Account, err error) {
	err = datastore.DB.Preload("Teams").Where("enterprise_id", enterpriseId).First(&account, id).Error
	return
}

func (r *repository) ReadEnterpriseAccountByPhone(enterpriseId int, phone string) (account *entities.Account, err error) {
	err = datastore.DB.Where("enterprise_id", enterpriseId).
		Where("phone LIKE ?", fmt.Sprintf("%%%s%%", phone)).First(&account).Error
	return
}

//func (r *repository) ReadEnterpriseAccountsByPhone(enterpriseId int, phones []string) (account *[]entities.Account, err error) {
//	regexp := ""
//	for _, phone := range phones {
//		regexp += phone + "|"
//	}
//
//	err = datastore.DB.Where("enterprise_id", enterpriseId).
//		Where("phone REGEXP ?", regexp).First(&account).Error
//	return
//}

func (r *repository) ReadEnterpriseAccountsByPhone(enterpriseId int, phones []string) (account *[]entities.Account, err error) {
	err = datastore.DB.Where("enterprise_id", enterpriseId).
		Where("phone IN ?", phones).Find(&account).Error
	return
}

// NewRepo is the single instance repo that is being created.
func NewRepo() Repository {
	return &repository{}
}
