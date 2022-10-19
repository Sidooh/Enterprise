package entities

import "github.com/Permify/permify-gorm/models"

type Account struct {
	ModelID

	Phone string `json:"phone" gorm:"unique"`
	Name  string `json:"name" gorm:"uniqueIndex"`
	Email string `json:"email" gorm:"uniqueIndex"`

	AccountId    uint `json:"account_id" gorm:"uniqueIndex"`
	EnterpriseId uint `json:"enterprise_id" gorm:"uniqueIndex"`

	// permify
	Roles []models.Role `gorm:"many2many:user_roles;OnUpdate:CASCADE,OnDelete:CASCADE;joinForeignKey:UserId"`

	ModelTimeStamps
}

type AccountWithEnterprise struct {
	Account

	Enterprise Enterprise
}

func (*AccountWithEnterprise) TableName() string {
	return "accounts"
}
