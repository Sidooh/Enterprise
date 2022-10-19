package presenter

type Enterprise struct {
	Id      uint   `json:"id"`
	Name    string `json:"name"`
	Country string `json:"country"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
}
