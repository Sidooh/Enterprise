package entities

import "time"

type Enterprise struct {
	ModelID

	Name            string     `json:"name" gorm:"unique"`
	Country         string     `json:"country"`
	Address         string     `json:"address"`
	Phone           string     `json:"phone" gorm:"unique"`
	PhoneVerifiedAt *time.Time `gorm:"type:timestamp null" json:"-"`
	Email           string     `json:"email" gorm:"unique"`
	EmailVerifiedAt *time.Time `gorm:"type:timestamp null" json:"-"`

	AccountId      uint `json:"account_id" gorm:"uniqueIndex"`
	FloatAccountId uint `json:"float_account_id" gorm:"uniqueIndex"`

	ModelTimeStamps
}
