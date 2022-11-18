package entities

type Account struct {
	ModelID

	Phone string `json:"phone" gorm:"unique"`
	Name  string `json:"name"`

	AccountId    uint `json:"account_id" gorm:"index"`
	EnterpriseId uint `json:"enterprise_id" gorm:"index"`

	// TODO: account is unique in enterprise by phone / accountId

	ModelTimeStamps
}

type AccountWithEnterprise struct {
	Account

	Enterprise Enterprise
}

func (*AccountWithEnterprise) TableName() string {
	return "accounts"
}
