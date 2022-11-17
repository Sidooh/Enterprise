package routes

import (
	"enterprise.sidooh/api/handlers"
	"enterprise.sidooh/pkg/services/team"
	"github.com/gofiber/fiber/v2"
)

func TeamRouter(app fiber.Router, service team.Service) {
	app.Get("/teams", handlers.GetTeams(service))
	app.Post("/teams", handlers.CreateTeam(service))
	app.Get("/teams/:id", handlers.GetTeam(service))
}
