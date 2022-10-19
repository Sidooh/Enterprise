package presenter

type Registration struct {
	Name          string `json:"name"`
	Country       string `json:"country"`
	Address       string `json:"address"`
	Phone         string `json:"phone"`
	Email         string `json:"email"`
	AdminName     string `json:"admin_name"`
	AdminPhone    string `json:"admin_phone"`
	AdminEmail    string `json:"admin_email"`
	AdminPassword string `json:"admin_password"`
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string   `json:"token"`
	Admin *Account `json:"admin"`
}
