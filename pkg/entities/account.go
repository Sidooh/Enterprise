package entities

type Account struct {
	ModelID

	Phone string `json:"phone" gorm:"unique"`
	Name  string `json:"name" gorm:"uniqueIndex"`
	Email string `json:"email" gorm:"uniqueIndex"`

	AccountId    uint `json:"account_id" gorm:"uniqueIndex"`
	EnterpriseId uint `json:"enterprise_id" gorm:"uniqueIndex"`

	ModelTimeStamps
}

type AccountWithEnterprise struct {
	Account

	Enterprise Enterprise
}

func (*AccountWithEnterprise) TableName() string {
	return "accounts"
}
