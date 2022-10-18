package routes

import (
	"enterprise.sidooh/api/handlers"
	"enterprise.sidooh/pkg/services/auth"
	"github.com/gofiber/fiber/v2"
)

func AuthRouter(app fiber.Router, service auth.Service) {
	app.Get("/auth/me", handlers.GetAuthAccount(service))
	app.Post("/auth/register", handlers.Register(service))
	app.Post("/auth/login", handlers.Login(service))
}
