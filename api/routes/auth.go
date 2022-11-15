package routes

import (
	"enterprise.sidooh/api/handlers"
	"enterprise.sidooh/api/middleware/jwt"
	"enterprise.sidooh/pkg/services/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"time"
)

func AuthRouter(app fiber.Router, service auth.Service) {
	app.Get("/auth/me", jwt.New(jwt.Config{
		Secret: viper.GetString("JWT_KEY"),
		Expiry: time.Duration(15) * time.Minute,
	}), handlers.GetAuthUser(service))
	app.Post("/auth/register", handlers.Register(service))
	app.Post("/auth/login", handlers.Login(service))
}
