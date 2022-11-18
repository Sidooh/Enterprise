package handlers

import (
	"enterprise.sidooh/api/middleware"
	"enterprise.sidooh/pkg/entities"
	"enterprise.sidooh/pkg/services/account"
	"enterprise.sidooh/utils"
	"errors"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type CreateAccountRequest struct {
	Name  string `json:"name" validate:"required"`
	Phone string `json:"phone" validate:"required,numeric,min=9"`
}

func GetAccount(service account.Service) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id, err := ctx.ParamsInt("id")
		if err != nil {
			ctx.Status(http.StatusBadRequest)
			return ctx.JSON(utils.ValidationErrorResponse(errors.New("invalid id parameter")))
		}

		fetched := new(entities.Account)

		if utils.IsSuperAdmin(ctx) {
			fetched, err = service.GetAccount(id)
		} else if utils.IsAdmin(ctx) {
			enterpriseId := utils.GetEnterpriseId(ctx)
			fetched, err = service.GetAccountForEnterprise(enterpriseId, id)
		} else {
			return utils.HandleUnauthorized(ctx)
		}

		if err != nil {
			return utils.HandleErrorResponse(ctx, err)
		}

		return utils.HandleSuccessResponse(ctx, fetched)
	}
}

func GetAccounts(service account.Service) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		fetched := new([]entities.Account)
		err := *new(error)

		if utils.IsSuperAdmin(ctx) {
			fetched, err = service.FetchAccounts()
		} else if utils.IsAdmin(ctx) {
			enterpriseId := utils.GetEnterpriseId(ctx)
			fetched, err = service.FetchAccountsForEnterprise(enterpriseId)
		} else {
			return utils.HandleUnauthorized(ctx)
		}

		if err != nil {
			return utils.HandleErrorResponse(ctx, err)
		}

		return utils.HandleSuccessResponse(ctx, fetched)
	}
}

func CreateAccount(service account.Service) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var request CreateAccountRequest
		if err := middleware.BindAndValidateRequest(ctx, &request); err != nil {
			return ctx.Status(http.StatusUnprocessableEntity).JSON(err)
		}

		fetched := new(entities.Account)
		err := *new(error)

		// TODO: Use permissions for this part - determine who can add accounts
		if utils.IsAdmin(ctx) {
			enterpriseId := utils.GetEnterpriseId(ctx)
			fetched, err = service.CreateAccount(&entities.Account{
				Phone:        request.Phone,
				Name:         request.Name,
				EnterpriseId: uint(enterpriseId),
			})
		} else {
			return utils.HandleUnauthorized(ctx)
		}

		if err != nil {
			return utils.HandleErrorResponse(ctx, err)
		}

		return utils.HandleSuccessResponse(ctx, fetched)
	}
}
