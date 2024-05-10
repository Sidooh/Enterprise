package presenter

type Team struct {
	Id           uint      `json:"id"`
	Name         string    `json:"name"`
	EnterpriseId int       `json:"enterprise_id"`
	Accounts     []Account `json:"accounts,omitempty"`
}
