package presenter

import (
	"github.com/gofiber/fiber/v2"
)

type Enterprise struct {
	Id uint `json:"id"`
}

func EnterpriseErrorResponse(err error) *fiber.Map {
	return &fiber.Map{
		"status": false,
		"error":  err.Error(),
	}
}
