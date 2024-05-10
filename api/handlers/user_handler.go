package handlers

import (
	"enterprise.sidooh/pkg/services/user"
	"enterprise.sidooh/utils"
	"errors"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

func GetUser(service user.Service) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id, err := ctx.ParamsInt("id")
		if err != nil {
			ctx.Status(http.StatusBadRequest)
			return ctx.JSON(utils.ValidationErrorResponse(errors.New("invalid id parameter")))
		}

		fetched, err := service.GetUser(id)
		if err != nil {
			return utils.HandleErrorResponse(ctx, err)
		}

		return utils.HandleSuccessResponse(ctx, fetched)
	}
}
