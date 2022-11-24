package handlers

import (
	"enterprise.sidooh/pkg/clients"
	"enterprise.sidooh/pkg/services/float"
	"enterprise.sidooh/utils"
	"github.com/gofiber/fiber/v2"
)

func GetFloatAccount(service float.Service) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		fetched := new(clients.FloatAccount)
		err := *new(error)

		/*if utils.IsSuperAdmin(ctx) {
			fetched, err = service.GetFloatAccount(id)
		} else */if utils.IsAdmin(ctx) {
			enterprise := utils.GetEnterprise(ctx)
			fetched, err = service.GetFloatAccountForEnterprise(enterprise)
		} else {
			return utils.HandleUnauthorized(ctx)
		}

		if err != nil {
			return utils.HandleErrorResponse(ctx, err)
		}

		return utils.HandleSuccessResponse(ctx, fetched)
	}
}
