package services

import "enterprise.sidooh/pkg/client"

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
	client.ApiResponse

	Data Account `json:"data"`
}

type FloatAccountApiResponse struct {
	client.ApiResponse

	Data FloatAccount `json:"data"`
}
