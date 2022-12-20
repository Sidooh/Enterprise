package presenter

import (
	"enterprise.sidooh/pkg/clients"
	"enterprise.sidooh/pkg/entities"
)

type Account struct {
	Id           uint               `json:"id"`
	Phone        string             `json:"phone"`
	Name         string             `json:"name"`
	EnterpriseId uint               `json:"enterprise_id"`
	Teams        []*entities.Team   `json:"teams,omitempty"`
	Vouchers     []*clients.Voucher `json:"vouchers,omitempty"`
}
