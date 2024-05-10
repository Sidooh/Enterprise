package team

import (
	"enterprise.sidooh/pkg/datastore"
	"enterprise.sidooh/pkg/entities"
	"fmt"
)

// Repository interface allows us to access the CRUD Operations here.
type Repository interface {
	ReadTeams() (*[]entities.Team, error)
	ReadTeam(id int) (*entities.Team, error)
	CreateTeam(team *entities.Team) (*entities.Team, error)
	AddTeamAccount(team *entities.Team, accountId uint) error
	//UpdateTeam(team *entities.Team, column string, value interface{}) (*entities.Team, error)
	//DeleteTeam(Id uint) error

	ReadEnterpriseTeams(enterpriseId int) (*[]entities.Team, error)
	ReadEnterpriseTeam(enterpriseId int, id int) (*entities.Team, error)
}
type repository struct {
}

func (r *repository) CreateTeam(team *entities.Team) (t *entities.Team, err error) {
	err = datastore.DB.Create(&team).Error
	return team, err
}

func (r *repository) AddTeamAccount(team *entities.Team, accountId uint) error {
	return datastore.DB.Model(&team).Association("Accounts").Append(&entities.Account{
		ModelID: entities.ModelID{Id: accountId},
	})
}

func (r *repository) ReadTeams() (teams *[]entities.Team, err error) {
	err = datastore.DB.Find(&teams).Error
	return
}

func (r *repository) ReadTeam(id int) (team *entities.Team, err error) {
	err = datastore.DB.First(&team, id).Error
	return
}

func (r *repository) UpdateTeam(team *entities.Team, column string, value interface{}) (t *entities.Team, err error) {
	err = datastore.DB.Model(&team).Update(column, value).Error
	return
}

//======================================================================================================================
//	Enterprise Limited Data Manipulation
//======================================================================================================================

func (r *repository) ReadEnterpriseTeams(enterpriseId int) (teams *[]entities.Team, err error) {
	err = datastore.DB.Preload("Accounts").Where("enterprise_id", enterpriseId).Find(&teams).Error
	return
}

func (r *repository) ReadEnterpriseTeam(enterpriseId int, id int) (team *entities.Team, err error) {
	err = datastore.DB.Preload("Accounts").Where("enterprise_id", enterpriseId).First(&team, id).Error
	return
}

func (r *repository) ReadEnterpriseTeamByName(enterpriseId int, name string) (team *entities.Team, err error) {
	err = datastore.DB.Where("enterprise_id", enterpriseId).
		Where("name LIKE ?", fmt.Sprintf("%%%s%%", name)).First(&team).Error
	return
}

// NewRepo is the single instance repo that is being created.
func NewRepo() Repository {
	return &repository{}
}
