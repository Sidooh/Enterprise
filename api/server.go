package api

import (
	"enterprise.sidooh/api/middleware"
	"enterprise.sidooh/api/routes"
	"enterprise.sidooh/pkg/client"
	"enterprise.sidooh/pkg/services/auth"
	"enterprise.sidooh/pkg/services/enterprise"
	"enterprise.sidooh/utils"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/spf13/viper"
	"net/http"
)

func setMiddleware(app *fiber.App) {
	app.Use(func(c *fiber.Ctx) error {
		// TODO: URGENT: Check out these headers and review them
		// Set some security headers:
		c.Set("X-XSS-Protection", "1; mode=block")
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("X-Download-Options", "noopen")
		c.Set("Strict-Transport-Security", "max-age=5184000")
		c.Set("X-Frame-Options", "SAMEORIGIN")
		c.Set("X-DNS-Prefetch-Control", "off")

		// Go to next middleware:
		return c.Next()
	})

	app.Use(cors.New())
	app.Use(limiter.New())
	app.Use(recover.New())
	app.Use(fiberLogger.New(fiberLogger.Config{Output: utils.GetLogFile("server.log")}))

	app.Use(favicon.New(favicon.Config{Next: func(c *fiber.Ctx) bool {
		return true
	}}))

	middleware.Validator = validator.New()
}

func setHealthCheckRoutes(app *fiber.App) {
	app.Get("/200", func(ctx *fiber.Ctx) error {
		return ctx.JSON("200")
	})

	app.Get("/500", func(ctx *fiber.Ctx) error {
		ctx.Status(http.StatusInternalServerError)
		return ctx.JSON("500")
	})
}

func setHandlers(app *fiber.App) {
	api := app.Group("/api")
	v1 := api.Group("/v1")

	authRep := auth.NewRepo()
	accountApi := client.InitAccountClient()
	authSrv := auth.NewService(authRep, accountApi)

	enterpriseRep := enterprise.NewRepo()
	enterpriseSrv := enterprise.NewService(enterpriseRep)

	routes.AuthRouter(v1, authSrv)
	routes.EnterpriseRouter(v1, enterpriseSrv)
}

func Server() *fiber.App {
	app := fiber.New(fiber.Config{
		Prefork: viper.GetBool("PREFORK"),
	})

	setMiddleware(app)
	setHealthCheckRoutes(app)

	setHandlers(app)

	return app
}
