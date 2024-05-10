package handlers

import (
	"enterprise.sidooh/api/middleware"
	"enterprise.sidooh/pkg/clients"
	"enterprise.sidooh/pkg/services/float"
	"enterprise.sidooh/utils"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type CreditFloatAccountRequest struct {
	Amount int `json:"amount" validate:"required,numeric,min=1000"`
	Phone  int `json:"phone" validate:"required,numeric,min=9"`
}

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

func GetFloatAccountTransactions(service float.Service) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		fetched := new([]clients.FloatAccountTransaction)
		err := *new(error)

		/*if utils.IsSuperAdmin(ctx) {
			fetched, err = service.GetFloatAccount(id)
		} else */if utils.IsAdmin(ctx) {
			enterprise := utils.GetEnterprise(ctx)
			fetched, err = service.GetFloatAccountTransactionsForEnterprise(enterprise)
		} else {
			return utils.HandleUnauthorized(ctx)
		}

		if err != nil {
			return utils.HandleErrorResponse(ctx, err)
		}

		return utils.HandleSuccessResponse(ctx, fetched)
	}
}

func CreditFloatAccount(service float.Service) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var request CreditFloatAccountRequest
		if err := middleware.BindAndValidateRequest(ctx, &request); err != nil {
			return ctx.Status(http.StatusUnprocessableEntity).JSON(err)
		}

		fetched := new(clients.FloatAccount)
		err := *new(error)

		// TODO: Use permissions for this part - determine who can add accounts
		if utils.IsAdmin(ctx) {
			enterprise := utils.GetEnterprise(ctx)
			fetched, err = service.CreditFloatAccountForEnterprise(enterprise, request.Amount, request.Phone)
		} else {
			return utils.HandleUnauthorized(ctx)
		}

		if err != nil {
			return utils.HandleErrorResponse(ctx, err)
		}

		return utils.HandleSuccessResponse(ctx, fetched)
	}
}
