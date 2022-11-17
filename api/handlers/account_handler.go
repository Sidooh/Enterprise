package handlers

import (
	"enterprise.sidooh/api/presenter"
	"enterprise.sidooh/pkg/services/account"
	"enterprise.sidooh/utils"
	"errors"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

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
		fetched := new([]presenter.Account)
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
