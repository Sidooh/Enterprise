package presenter

type Registration struct {
	Name string `json:"name"`
	//Country       string `json:"country"`
	//Address       string `json:"address"`
	Phone     string `json:"phone"`
	AdminName string `json:"admin_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}
