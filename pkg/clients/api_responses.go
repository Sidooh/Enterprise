package clients

import "time"

type Account struct {
	Id     int    `json:"id"`
	Phone  string `json:"phone"`
	Active bool   `json:"active"`
}

type FloatAccount struct {
	Id            int    `json:"id"`
	Balance       int    `json:"balance"`
	AccountId     int    `json:"account_id"`
	FloatableId   int    `json:"floatable_id"`
	FloatableType string `json:"floatable_type"`
}

type FloatAccountTransaction struct {
	Id             int       `json:"id"`
	Type           string    `json:"type"`
	Amount         int       `json:"amount"`
	Description    string    `json:"description"`
	FloatAccountId int       `json:"float_account_id"`
	CreatedAt      time.Time `json:"created_at"`
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

type VoucherTransaction struct {
	Id          int       `json:"id"`
	Type        string    `json:"type"`
	Amount      int       `json:"amount"`
	Description string    `json:"description"`
	VoucherId   int       `json:"voucher_id"`
	CreatedAt   time.Time `json:"created_at"`
	Voucher     Voucher   `json:"voucher,omitempty"`
}

type Voucher struct {
	Id        int `json:"id"`
	AccountId int `json:"account_id"`
	Balance   int `json:"balance"`
}

type DashboardStatistics struct {
	FloatBalance      int `json:"float_balance"`
	AccountsCount     int `json:"accounts_count"`
	VouchersDisbursed int `json:"vouchers_disbursed"`
}
