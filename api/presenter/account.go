package presenter

type Account struct {
	Id           uint   `json:"id"`
	Phone        string `json:"phone"`
	Name         string `json:"name"`
	EnterpriseId string `json:"enterprise_id"`
	//Enterprise Enterprise `json:"enterprise,omitempty"`
}
