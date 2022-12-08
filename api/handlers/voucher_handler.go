package handlers

import (
	"enterprise.sidooh/api/middleware"
	"enterprise.sidooh/pkg/clients"
	"enterprise.sidooh/pkg/services/voucher"
	"enterprise.sidooh/utils"
	"errors"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type CreateVoucherTypeRequest struct {
	Name string `json:"name" validate:"required"`
}

type DisburseVoucherTypeRequest struct {
	AccountId int `json:"account_id" validate:"required,numeric"`
	Amount    int `json:"amount" validate:"required,numeric,min=100,max=10000"`
}

func GetVoucherType(service voucher.Service) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id, err := ctx.ParamsInt("id")
		if err != nil {
			ctx.Status(http.StatusBadRequest)
			return ctx.JSON(utils.ValidationErrorResponse(errors.New("invalid id parameter")))
		}

		fetched := new(clients.VoucherType)

		/*if utils.IsSuperAdmin(ctx) {
			fetched, err = service.GetVoucherType(id)
		} else */if utils.IsAdmin(ctx) {
			enterpriseId := utils.GetEnterpriseId(ctx)
			fetched, err = service.GetVoucherTypeForEnterprise(enterpriseId, id)
		} else {
			return utils.HandleUnauthorized(ctx)
		}

		if err != nil {
			return utils.HandleErrorResponse(ctx, err)
		}

		return utils.HandleSuccessResponse(ctx, fetched)
	}
}

func GetVoucherTypes(service voucher.Service) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		fetched := new([]clients.VoucherType)
		err := *new(error)

		/*if utils.IsSuperAdmin(ctx) {
			fetched, err = service.FetchVoucherTypes()
		} else*/if utils.IsAdmin(ctx) {
			enterprise := utils.GetEnterprise(ctx)
			fetched, err = service.FetchVoucherTypesForEnterprise(int(enterprise.AccountId))
		} else {
			return utils.HandleUnauthorized(ctx)
		}

		if err != nil {
			return utils.HandleErrorResponse(ctx, err)
		}

		return utils.HandleSuccessResponse(ctx, fetched)
	}
}

func CreateVoucherType(service voucher.Service) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var request CreateVoucherTypeRequest
		if err := middleware.BindAndValidateRequest(ctx, &request); err != nil {
			return ctx.Status(http.StatusUnprocessableEntity).JSON(err)
		}

		fetched := new(clients.VoucherType)
		err := *new(error)

		// TODO: Use permissions for this part - determine who can add voucherTypes
		if utils.IsAdmin(ctx) {
			enterprise := utils.GetEnterprise(ctx)
			fetched, err = service.CreateVoucherType(int(enterprise.AccountId), request.Name)
		} else {
			return utils.HandleUnauthorized(ctx)
		}

		if err != nil {
			return utils.HandleErrorResponse(ctx, err)
		}

		return utils.HandleSuccessResponse(ctx, fetched)
	}
}

func DisburseVoucherType(service voucher.Service) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id, err := ctx.ParamsInt("id")
		if err != nil {
			ctx.Status(http.StatusBadRequest)
			return ctx.JSON(utils.ValidationErrorResponse(errors.New("invalid id parameter")))
		}

		var request DisburseVoucherTypeRequest
		if err = middleware.BindAndValidateRequest(ctx, &request); err != nil {
			return ctx.Status(http.StatusUnprocessableEntity).JSON(err)
		}

		fetched := new(clients.VoucherType)

		// TODO: Use permissions for this part - determine who can add voucherTypes
		if utils.IsAdmin(ctx) {
			enterprise := utils.GetEnterprise(ctx)
			fetched, err = service.DisburseVoucherType(enterprise, id, request.AccountId, request.Amount)
		} else {
			return utils.HandleUnauthorized(ctx)
		}

		if err != nil {
			return utils.HandleErrorResponse(ctx, err)
		}

		return utils.HandleSuccessResponse(ctx, fetched)
	}
}
