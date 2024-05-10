package routes

import (
	"enterprise.sidooh/api/handlers"
	"enterprise.sidooh/pkg/services/user"
	"github.com/gofiber/fiber/v2"
)

func UserRouter(app fiber.Router, service user.Service) {
	app.Get("/users/:id", handlers.GetUser(service))
}
