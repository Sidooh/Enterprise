package api

import (
	"enterprise.sidooh/api/routes"
	"enterprise.sidooh/pkg/enterprise"
	"enterprise.sidooh/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"net/http"
)

func Server() *fiber.App {
	app := fiber.New(fiber.Config{
		Prefork: true,
	})
	app.Use(cors.New())
	app.Use(fiberLogger.New(fiberLogger.Config{Output: utils.GetLogFile("server.log")}))

	app.Get("/200", func(ctx *fiber.Ctx) error {
		return ctx.JSON("200")
	})

	app.Get("/500", func(ctx *fiber.Ctx) error {
		ctx.Status(http.StatusInternalServerError)
		return ctx.JSON("500")
	})

	api := app.Group("/api")
	v1 := api.Group("/v1")

	routes.EnterpriseRouter(v1, enterprise.NewService(enterprise.NewRepo()))

	return app
}
