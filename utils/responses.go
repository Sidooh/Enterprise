package utils

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
)

type JsonResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

func (j JsonResponse) Error() string {
	v, _ := jsoniter.MarshalToString(j)
	return v
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Param   string `json:"param"`
	Value   string `json:"value,omitempty"`
}

func SuccessResponse(data interface{}) JsonResponse {
	return JsonResponse{
		Status: true,
		Data:   data,
	}
}

func ErrorResponse(message string, errors interface{}) JsonResponse {
	return JsonResponse{
		Status:  false,
		Message: message,
		Errors:  errors,
	}
}

func ServerErrorResponse() JsonResponse {
	return ErrorResponse("something went wrong, please try again", nil)
}

func UnauthorizedErrorResponse() JsonResponse {
	return ErrorResponse("invalid credentials", nil)
}

func NotFoundErrorResponse() JsonResponse {
	return ErrorResponse("not found", nil)
}

func ValidationErrorResponse(errors interface{}) JsonResponse {
	return ErrorResponse("the request is invalid", errors)
}

func HandleErrorResponse(ctx *fiber.Ctx, err error) error {
	log.Error(err)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ctx.Status(http.StatusNotFound).JSON(NotFoundErrorResponse())
	}

	if err.Error() == "invalid credentials" {
		return ctx.Status(http.StatusUnauthorized).JSON(UnauthorizedErrorResponse())
	}

	return ctx.Status(http.StatusInternalServerError).JSON(ServerErrorResponse())
}

func HandleSuccessResponse(ctx *fiber.Ctx, data interface{}) error {
	return ctx.Status(http.StatusOK).JSON(SuccessResponse(data))
}
