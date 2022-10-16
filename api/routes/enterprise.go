package routes

import (
	"enterprise.sidooh/api/handlers"
	"enterprise.sidooh/pkg/enterprise"
	"github.com/gofiber/fiber/v2"
)

func EnterpriseRouter(app fiber.Router, service enterprise.Service) {
	app.Get("/enterprises", handlers.GetEnterprises(service))
	app.Get("/enterprises/:id", handlers.GetEnterprise(service))
}
