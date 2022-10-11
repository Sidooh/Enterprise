package presenter

import (
	"enterprise-gateway.sidooh/pkg/entities"
	"github.com/gofiber/fiber/v2"
)

type Enterprise struct {
	Id int `json:"id"`
}

func EnterpriseSuccessResponse(data *entities.Enterprise) *fiber.Map {
	enterprise := Enterprise{Id: data.Id}

	return &fiber.Map{
		"status": true,
		"data":   enterprise,
	}
}

func EnterpriseErrorResponse(err error) *fiber.Map {
	return &fiber.Map{
		"status": false,
		"error":  err.Error(),
	}
}
