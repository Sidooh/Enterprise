package handlers

import (
	"enterprise.sidooh/api/middleware"
	"enterprise.sidooh/api/presenter"
	"enterprise.sidooh/pkg/entities"
	"enterprise.sidooh/pkg/services/team"
	"enterprise.sidooh/utils"
	"errors"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type CreateTeamRequest struct {
	Name string `json:"name" validate:"required"`
}

// TODO: Refactor all handlers to do the data transformation from entities to presenters
//  Or should we do it in service?

func GetTeam(service team.Service) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id, err := ctx.ParamsInt("id")
		if err != nil {
			ctx.Status(http.StatusBadRequest)
			return ctx.JSON(utils.ValidationErrorResponse(errors.New("invalid id parameter")))
		}

		fetched := new(entities.Team)

		if utils.IsSuperAdmin(ctx) {
			fetched, err = service.GetTeam(id)
		} else if utils.IsAdmin(ctx) {
			enterpriseId := utils.GetEnterpriseId(ctx)
			fetched, err = service.GetTeamForEnterprise(enterpriseId, id)
		} else {
			return utils.HandleUnauthorized(ctx)
		}

		if err != nil {
			return utils.HandleErrorResponse(ctx, err)
		}

		return utils.HandleSuccessResponse(ctx, fetched)
	}
}

func GetTeams(service team.Service) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		fetched := new([]entities.Team)
		err := *new(error)

		if utils.IsSuperAdmin(ctx) {
			fetched, err = service.FetchTeams()
		} else if utils.IsAdmin(ctx) {
			enterpriseId := utils.GetEnterpriseId(ctx)
			fetched, err = service.FetchTeamsForEnterprise(enterpriseId)
		} else {
			return utils.HandleUnauthorized(ctx)
		}

		if err != nil {
			return utils.HandleErrorResponse(ctx, err)
		}

		return utils.HandleSuccessResponse(ctx, fetched)
	}
}

func CreateTeam(service team.Service) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var request CreateTeamRequest
		if err := middleware.BindAndValidateRequest(ctx, &request); err != nil {
			return ctx.Status(http.StatusUnprocessableEntity).JSON(err)
		}

		fetched := new(presenter.Team)
		err := *new(error)

		// TODO: Use permissions for this part - determine who can add teams
		if utils.IsAdmin(ctx) {
			enterpriseId := utils.GetEnterpriseId(ctx)
			fetched, err = service.CreateTeam(&entities.Team{
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
