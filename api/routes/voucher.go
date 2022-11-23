package routes

import (
	"enterprise.sidooh/api/handlers"
	"enterprise.sidooh/pkg/services/voucher"
	"github.com/gofiber/fiber/v2"
)

func VoucherRouter(app fiber.Router, service voucher.Service) {
	app.Get("/voucher-types", handlers.GetVoucherTypes(service))
	app.Post("/voucher-types", handlers.CreateVoucherType(service))
	app.Get("/voucher-types/:id", handlers.GetVoucherType(service))
}
