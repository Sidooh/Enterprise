package handlers

import (
	"enterprise.sidooh/api/middleware"
	"enterprise.sidooh/api/middleware/jwt"
	"enterprise.sidooh/api/presenter"
	"enterprise.sidooh/pkg"
	"enterprise.sidooh/pkg/services/auth"
	"enterprise.sidooh/utils"
	"fmt"
	"github.com/gofiber/fiber/v2"
	jwt2 "github.com/golang-jwt/jwt/v4"
	"github.com/hesahesa/pwdbro"
	"github.com/spf13/viper"
	"golang.org/x/exp/slices"
	"net/http"
	"strings"
	"time"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=64"`
}

type OTPGenerationRequest struct {
	Id      int    `json:"id" validate:"required,numeric"`
	Channel string `json:"channel"`
}

type OTPValidationRequest struct {
	Id  int `json:"id" validate:"required,numeric"`
	Otp int `json:"otp" validate:"required,numeric,min=100000,max=999999"`
}

type RegisterRequest struct {
	Name      string `json:"name" validate:"required"`
	Phone     string `json:"phone" validate:"required,numeric,min=9"`
	Email     string `json:"email" validate:"required,email"`
	AdminName string `json:"admin_name" validate:"required"`
	Password  string `json:"password" validate:"required,min=10,max=64"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

func GetAuthUser(service auth.Service) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		claims := ctx.Locals("jwtClaims").(jwt2.MapClaims)
		fmt.Println(ctx.Locals("user"))

		user, err := service.User(int(claims["id"].(float64)))
		if err != nil {
			return utils.HandleErrorResponse(ctx, err)
		}

		return utils.HandleSuccessResponse(ctx, user)
	}
}

func Register(service auth.Service) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var request RegisterRequest
		if err := middleware.BindAndValidateRequest(ctx, &request); err != nil {
			return ctx.Status(http.StatusUnprocessableEntity).JSON(err)
		}

		// TODO: Check password rules
		// use owasp recommendations

		// OWASP placeholder
		pwd := pwdbro.NewDefaultPwdBro()
		status, err := pwd.RunParallelChecks(request.Password)
		if err != nil {
			return utils.HandleErrorResponse(ctx, err)
		}

		for _, resp := range status {
			// in practice, you will want to evaluate based on the
			// resp.Safe field and not just printing it\
			if !resp.Safe {
				return ctx.Status(http.StatusUnprocessableEntity).JSON(utils.ValidationErrorResponse(utils.ValidationError{
					Field:   "password",
					Message: "password does not meet requirements",
					Param:   "password",
					Value:   request.Password,
				}))
			}
		}

		register, err := service.Register(presenter.Registration{
			Name:      request.Name,
			AdminName: request.AdminName,
			Phone:     request.Phone,
			Email:     request.Email,
			Password:  request.Password,
		})
		if err != nil {
			return ctx.Status(http.StatusUnprocessableEntity).JSON(utils.ValidationErrorResponse(err))
			//return utils.HandleErrorResponse(ctx, err)
		}

		return utils.HandleSuccessResponse(ctx, register)
	}
}

func Login(service auth.Service) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var request LoginRequest
		if err := middleware.BindAndValidateRequest(ctx, &request); err != nil {
			return ctx.Status(http.StatusUnprocessableEntity).JSON(err)
		}

		authData, err := service.Login(presenter.Login{
			Email:    request.Email,
			Password: strings.TrimSpace(request.Password),
		})
		if err != nil {
			return utils.HandleErrorResponse(ctx, err)
		}

		validity := time.Duration(viper.GetInt("ACCESS_TOKEN_VALIDITY")) * time.Minute
		token, err := jwt.Encode(&jwt2.MapClaims{
			"name":      authData.User.Name,
			"email":     authData.User.Email,
			"id":        authData.User.Id,
			"valid_mfa": false,
		}, validity)

		authData.Token = token

		return utils.HandleSuccessResponse(ctx, authData)
	}
}

func GenerateOtp(service auth.Service) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var request OTPGenerationRequest
		if err := middleware.BindAndValidateRequest(ctx, &request); err != nil {
			return ctx.Status(http.StatusUnprocessableEntity).JSON(err)
		}

		if !slices.Contains([]string{utils.SMS, utils.MAIL}, request.Channel) {
			return utils.HandleErrorResponse(ctx, pkg.ErrInvalidChannel)
		}

		err := service.GenerateOTP(request.Id, request.Channel)
		if err != nil {
			return utils.HandleErrorResponse(ctx, err)
		}

		return utils.HandleSuccessResponse(ctx, true)
	}
}

func VerifyOtp(service auth.Service) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var request OTPValidationRequest
		if err := middleware.BindAndValidateRequest(ctx, &request); err != nil {
			return ctx.Status(http.StatusUnprocessableEntity).JSON(err)
		}

		authData, err := service.ValidateOTP(request.Id, request.Otp)
		if err != nil {
			return utils.HandleErrorResponse(ctx, err)
		}

		validity := time.Duration(viper.GetInt("ACCESS_TOKEN_VALIDITY")) * time.Minute
		token, err := jwt.Encode(&jwt2.MapClaims{
			"name":      authData.User.Name,
			"email":     authData.User.Email,
			"id":        authData.User.Id,
			"valid_mfa": true,
		}, validity)

		authData.Token = token

		return utils.HandleSuccessResponse(ctx, authData)
	}
}
