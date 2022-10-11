package routes

import (
	"enterprise-gateway.sidooh/api/handlers"
	"enterprise-gateway.sidooh/pkg/enterprise"
	"github.com/gofiber/fiber/v2"
)

func EnterpriseRouter(app fiber.Router, service enterprise.Service) {
	app.Get("/enterprises/:id", handlers.GetEnterprise(service))
}
