package routes

import (
	"enterprise.sidooh/api/handlers"
	"enterprise.sidooh/pkg/services/float"
	"github.com/gofiber/fiber/v2"
)

func FloatRouter(app fiber.Router, service float.Service) {
	app.Get("/float-account", handlers.GetFloatAccount(service))
	app.Get("/float-account/transactions", handlers.GetFloatAccountTransactions(service))
}
