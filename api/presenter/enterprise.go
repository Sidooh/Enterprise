package presenter

type Enterprise struct {
	Id      uint   `json:"id"`
	Name    string `json:"name"`
	Country string `json:"country,omitempty"`
	Address string `json:"address,omitempty"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
}

type EnterpriseWithUser struct {
	Enterprise
	User User
}
