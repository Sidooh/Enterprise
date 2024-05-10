package entities

type Team struct {
	ModelID

	Name string `json:"name"`

	EnterpriseId uint `json:"enterprise_id" gorm:"index"`

	Accounts []*Account `gorm:"many2many:team_accounts;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"accounts"`

	ModelTimeStamps
}
