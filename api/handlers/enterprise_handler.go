package handlers

import (
	"enterprise.sidooh/pkg/services/enterprise"
	"enterprise.sidooh/utils"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

func GetEnterprise(service enterprise.Service) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id, err := ctx.ParamsInt("id")
		if err != nil {
			ctx.Status(http.StatusBadRequest)
			return ctx.JSON(utils.ValidationErrorResponse(errors.New("invalid id parameter")))
		}

		fetched, err := service.GetEnterprise(id)
		if err != nil {
			return utils.HandleErrorResponse(ctx, err)
		}

		return utils.HandleSuccessResponse(ctx, fetched)
	}
}

func GetEnterprises(service enterprise.Service) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		claimData := ctx.Locals("jwtClaims")
		fmt.Println(claimData)

		fetched, err := service.FetchEnterprises()
		if err != nil {
			return utils.HandleErrorResponse(ctx, err)
		}

		return utils.HandleSuccessResponse(ctx, fetched)
	}
}
