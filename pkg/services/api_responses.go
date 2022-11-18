package services

import "enterprise.sidooh/pkg/clients"

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

type AccountApiResponse struct {
	clients.ApiResponse

	Data Account `json:"data"`
}

type FloatAccountApiResponse struct {
	clients.ApiResponse

	Data FloatAccount `json:"data"`
}
