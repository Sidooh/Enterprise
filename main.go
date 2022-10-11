package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"log"
	"sidooh-enterprise-gateway/api/routes"
	"sidooh-enterprise-gateway/pkg/enterprise"
)

func main() {
	app := fiber.New()
	app.Use(cors.New())
	app.Use(logger.New())

	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.JSON("Yaaay!!!")
	})

	api := app.Group("/api")
	v1 := api.Group("/v1")

	//enterpriseService := enterprise.NewService()
	routes.EnterpriseRouter(v1, enterprise.NewService())

	log.Fatal(app.Listen(":8006"))
}
