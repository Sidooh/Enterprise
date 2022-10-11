package routes

import (
	"github.com/gofiber/fiber/v2"
	"sidooh-enterprise-gateway/api/handlers"
	"sidooh-enterprise-gateway/pkg/enterprise"
)

func EnterpriseRouter(app fiber.Router, service enterprise.Service) {
	app.Get("/enterprises/:id", handlers.GetEnterprise(service))
}
