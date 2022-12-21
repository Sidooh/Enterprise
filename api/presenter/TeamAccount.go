package presenter

type TeamAccount struct {
	Id        uint `json:"id"`
	TeamId    int  `json:"team_id"`
	AccountId int  `json:"account_id"`
	Team      Team `json:"team,omitempty"`
	Account   Team `json:"account,omitempty"`
}
