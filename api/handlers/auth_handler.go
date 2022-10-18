package handlers

import (
	"enterprise.sidooh/api/middleware"
	"enterprise.sidooh/api/presenter"
	"enterprise.sidooh/pkg/services/auth"
	"enterprise.sidooh/utils"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type LoginRequest struct {
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required,min=8,max=64"`
}

type RegisterRequest struct {
	Name          string `json:"name" validate:"required"`
	Country       string `json:"country" validate:"required"`
	Address       string `json:"address" validate:"required"`
	Phone         string `json:"phone" validate:"required,numeric,min=9"`
	Email         string `json:"email" validate:"required,email"`
	AdminName     string `json:"admin_name" validate:"required"`
	AdminPhone    string `json:"admin_phone" validate:"required,numeric,min=9"`
	AdminEmail    string `json:"admin_email" validate:"required,email"`
	AdminPassword string `json:"admin_password" validate:"required,min=8,max=64"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

func GetAuthAccount(service auth.Service) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		return utils.HandleSuccessResponse(ctx, ctx.Get("Authorization"))
	}
}

func Register(service auth.Service) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var request RegisterRequest
		if err := middleware.BindAndValidateRequest(ctx, &request); err != nil {
			return ctx.Status(http.StatusUnprocessableEntity).JSON(err)
		}

		register, err := service.Register(presenter.Registration{
			Name:          "",
			Country:       "",
			Address:       "",
			Phone:         "",
			Email:         "",
			AdminName:     "",
			AdminPhone:    "",
			AdminEmail:    "",
			AdminPassword: "",
		})
		if err != nil {
			return utils.HandleErrorResponse(ctx, err)
		}

		return utils.HandleSuccessResponse(ctx, register)
	}
}

func Login(service auth.Service) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var requestBody LoginRequest
		err := ctx.BodyParser(&requestBody)
		if err != nil {
			ctx.Status(http.StatusBadRequest)
			return ctx.JSON(utils.HandleErrorResponse(ctx, err))
		}

		return utils.HandleSuccessResponse(ctx, ctx.Get("Authorization"))
	}
}
