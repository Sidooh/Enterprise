package presenter

import "enterprise.sidooh/pkg/entities"

type Account struct {
	Id           uint   `json:"id"`
	Phone        string `json:"phone"`
	Name         string `json:"name"`
	EnterpriseId string `json:"enterprise_id"`
	//Enterprise Enterprise `json:"enterprise,omitempty"`
	Team entities.Team `json:"team,omitempty"`
}
