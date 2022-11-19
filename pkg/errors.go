package pkg

import "errors"

var (
	ErrInvalidEnterprise = errors.New("enterprise details are invalid")

	ErrInvalidUser = errors.New("user details are invalid")

	ErrInvalidAccount = errors.New("account details are invalid")

	ErrUnauthorized = errors.New("unauthorized")

	ErrServerError = errors.New("something went wrong")
)
