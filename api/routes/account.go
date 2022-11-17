package routes

import (
	"enterprise.sidooh/api/handlers"
	"enterprise.sidooh/pkg/services/account"
	"github.com/gofiber/fiber/v2"
)

func AccountRouter(app fiber.Router, service account.Service) {
	app.Get("/accounts", handlers.GetAccounts(service))
	app.Post("/accounts", handlers.CreateAccount(service))
	app.Get("/accounts/:id", handlers.GetAccount(service))
}
