package presenter

type User struct {
	Id    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserWithRelations struct {
	Id          uint       `json:"id"`
	Name        string     `json:"name"`
	Email       string     `json:"email"`
	Enterprise  Enterprise `json:"enterprise,omitempty"`
	Roles       []string   `json:"roles"`
	Permissions []string   `json:"permissions,omitempty"`
}
