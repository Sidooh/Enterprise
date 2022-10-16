package entities

type Enterprise struct {
	ModelID

	Name    string `json:"name" gorm:"uniqueIndex"`
	Country string `json:"country"`
	Address string `json:"address"`
	Phone   string `json:"phone" gorm:"unique"`
	Email   string `json:"email" gorm:"unique"`

	AccountId uint `json:"account_id" gorm:"uniqueIndex"`

	ModelTimeStamps
}
