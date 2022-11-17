package entities

type Account struct {
	ModelID

	Phone string `json:"phone" gorm:"unique"`
	Name  string `json:"name"`
	Email string `json:"email" gorm:"unique"`

	AccountId    uint `json:"account_id" gorm:"index"`
	EnterpriseId uint `json:"enterprise_id" gorm:"index"`

	ModelTimeStamps
}

type AccountWithEnterprise struct {
	Account

	Enterprise Enterprise
}

func (*AccountWithEnterprise) TableName() string {
	return "accounts"
}
