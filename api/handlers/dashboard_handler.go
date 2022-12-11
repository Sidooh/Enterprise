package handlers

import (
	"enterprise.sidooh/pkg/clients"
	"enterprise.sidooh/pkg/services/dashboard"
	"enterprise.sidooh/utils"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func GetStatistics(service dashboard.Service) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		fetched := new(clients.DashboardStatistics)
		err := *new(error)

		if utils.IsAdmin(ctx) {
			enterprise := utils.GetEnterprise(ctx)
			fetched, err = service.GetDashboardStatistics(enterprise)
		} else {
			return utils.HandleUnauthorized(ctx)
		}

		if err != nil {
			return utils.HandleErrorResponse(ctx, err)
		}

		return utils.HandleSuccessResponse(ctx, fetched)
	}
}

func GetRecentVoucherTransactionsForEnterprise(service dashboard.Service) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		fetched := new([]clients.VoucherTransaction)
		err := *new(error)

		limit, err := strconv.ParseInt(ctx.Query("limit"), 0, 8)
		if err != nil {
			return utils.HandleErrorResponse(ctx, err)
		}

		if utils.IsAdmin(ctx) {
			enterprise := utils.GetEnterprise(ctx)
			fetched, err = service.GetRecentVoucherTransactionsForEnterprise(enterprise, int(limit))
		} else {
			return utils.HandleUnauthorized(ctx)
		}

		if err != nil {
			return utils.HandleErrorResponse(ctx, err)
		}

		return utils.HandleSuccessResponse(ctx, fetched)
	}
}

func GetRecentFloatTransactionsForEnterprise(service dashboard.Service) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		fetched := new([]clients.FloatAccountTransaction)
		err := *new(error)

		limit, err := strconv.ParseInt(ctx.Query("limit"), 0, 8)
		if err != nil {
			limit = 0
			//return utils.HandleErrorResponse(ctx, err)
		}

		if utils.IsAdmin(ctx) {
			enterprise := utils.GetEnterprise(ctx)
			fetched, err = service.GetRecentFloatAccountTransactionsForEnterprise(enterprise, int(limit))
		} else {
			return utils.HandleUnauthorized(ctx)
		}

		if err != nil {
			return utils.HandleErrorResponse(ctx, err)
		}

		return utils.HandleSuccessResponse(ctx, fetched)
	}
}
