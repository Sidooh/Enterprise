package handlers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"sidooh-enterprise-gateway/api/presenter"
	"sidooh-enterprise-gateway/pkg/enterprise"
)

func GetEnterprise(service enterprise.Service) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		//var requestBody entities.Enterprise
		//err := ctx.BodyParser(&requestBody)

		id, err := ctx.ParamsInt("id")
		if err != nil {
			ctx.Status(http.StatusBadRequest)
			return ctx.JSON(presenter.EnterpriseErrorResponse(errors.New("invalid id parameter")))
		}

		fetched, err := service.GetEnterprise(id)
		if err != nil {
			ctx.Status(http.StatusInternalServerError)
			return ctx.JSON(presenter.EnterpriseErrorResponse(errors.New("something went wrong")))
		}

		return ctx.JSON(presenter.EnterpriseSuccessResponse(fetched))
	}
}
