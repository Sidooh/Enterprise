package entities

type TeamAccount struct {
	ModelID

	TeamId    uint `json:"team_id" gorm:"index"`
	AccountId uint `json:"account_id" gorm:"index"`

	ModelTimeStamps
}
