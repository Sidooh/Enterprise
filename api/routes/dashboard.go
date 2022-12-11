package routes

import (
	"enterprise.sidooh/api/handlers"
	"enterprise.sidooh/pkg/services/dashboard"
	"github.com/gofiber/fiber/v2"
)

func DashboardRouter(app fiber.Router, service dashboard.Service) {
	app.Get("/dashboard/statistics", handlers.GetStatistics(service))
	app.Get("/dashboard/recent-voucher-transactions", handlers.GetRecentFloatTransactionsForEnterprise(service))
	app.Get("/dashboard/recent-float-transactions", handlers.GetRecentFloatTransactionsForEnterprise(service))
}
