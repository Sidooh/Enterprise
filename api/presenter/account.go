package presenter

type Account struct {
	Id    uint   `json:"id"`
	Phone string `json:"phone"`
	Name  string `json:"name"`
	Email string `json:"email"`
	//Enterprise Enterprise `json:"enterprise,omitempty"`
}
