package presenter

import (
	"github.com/gofiber/fiber/v2"
)

type Enterprise struct {
	Id uint `json:"id"`
}

func EnterprisesSuccessResponse(data *[]Enterprise) *fiber.Map {
	return &fiber.Map{
		"status": true,
		"data":   data,
	}
}

func EnterpriseSuccessResponse(data *Enterprise) *fiber.Map {
	//enterprise := Enterprise{Id: data.Id}

	return &fiber.Map{
		"status": true,
		"data":   data,
	}
}

func EnterpriseErrorResponse(err error) *fiber.Map {
	return &fiber.Map{
		"status": false,
		"error":  err.Error(),
	}
}
