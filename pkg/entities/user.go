package entities

import "time"

type User struct {
	ModelID

	Phone           string     `json:"phone" gorm:"unique"`
	Name            string     `json:"name"`
	Email           string     `json:"email" gorm:"uniqueIndex"`
	EmailVerifiedAt *time.Time `gorm:"type:timestamp" json:"-"`
	Password        string     `json:"-"`

	EnterpriseId uint `json:"enterprise_id" gorm:"index"`

	ModelTimeStamps
}

type UserWithEnterprise struct {
	User

	Enterprise Enterprise
}

func (*UserWithEnterprise) TableName() string {
	return "users"
}
