package clients

type Account struct {
	Id     int    `json:"id"`
	Phone  string `json:"phone"`
	Active bool   `json:"active"`
}

type FloatAccount struct {
	Id            int    `json:"id"`
	AccountId     int    `json:"account_id"`
	FloatableId   int    `json:"floatable_id"`
	FloatableType string `json:"floatable_type"`
}

type VoucherType struct {
	Id          int         `json:"id"`
	Name        string      `json:"name"`
	IsLocked    int         `json:"is_locked"`
	LimitAmount int         `json:"limit_amount"`
	ExpiresAt   string      `json:"expires_at,omitempty"`
	Settings    interface{} `json:"settings"`
	AccountId   int         `json:"account_id"`
	Vouchers    []Voucher   `json:"vouchers,omitempty"`
}

type Voucher struct {
	Id        int `json:"id"`
	AccountId int `json:"account_id"`
	Balance   int `json:"balance"`
}
