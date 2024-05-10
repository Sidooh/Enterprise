package handlers

import (
	"enterprise.sidooh/api/middleware"
	"enterprise.sidooh/api/presenter"
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

type CreateAccountsRequest struct {
	Accounts []CreateAccountRequest `json:"accounts" validate:"dive"`
}

func GetAccount(service account.Service) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id, err := ctx.ParamsInt("id")
		if err != nil {
			ctx.Status(http.StatusBadRequest)
			return ctx.JSON(utils.ValidationErrorResponse(errors.New("invalid id parameter")))
		}

		fetched := new(presenter.Account)

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

		phone, err := utils.GetPhoneByCountry("KE", request.Phone)
		if err != nil {
			return ctx.Status(http.StatusUnprocessableEntity).JSON(err)
		}

		fetched := new(entities.Account)

		// TODO: Use permissions for this part - determine who can add accounts
		if utils.IsAdmin(ctx) {
			enterpriseId := utils.GetEnterpriseId(ctx)
			fetched, err = service.CreateAccount(&entities.Account{
				Phone:        phone,
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

func CreateBulkAccounts(service account.Service) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var request CreateAccountsRequest
		if err := middleware.BindAndValidateRequest(ctx, &request); err != nil {
			return ctx.Status(http.StatusUnprocessableEntity).JSON(err)
		}

		fetched := new([]entities.Account)
		err := map[string]string{}

		// TODO: Use permissions for this part - determine who can add accounts
		if utils.IsAdmin(ctx) {
			enterpriseId := utils.GetEnterpriseId(ctx)

			var accounts []entities.Account

			for _, accountRequest := range request.Accounts {
				phone, err := utils.GetPhoneByCountry("KE", accountRequest.Phone)
				if err != nil {
					return ctx.Status(http.StatusUnprocessableEntity).JSON(err)
				}

				accounts = append(accounts, entities.Account{
					Phone:        phone,
					Name:         accountRequest.Name,
					EnterpriseId: uint(enterpriseId),
				})
			}

			fetched, err = service.CreateBulkAccounts(accounts)
		} else {
			return utils.HandleUnauthorized(ctx)
		}

		r := struct {
			Success interface{} `json:"success,omitempty"`
			Failed  interface{} `json:"failed,omitempty"`
		}{
			fetched, err,
		}

		return utils.HandleSuccessResponse(ctx, r)
	}
}
