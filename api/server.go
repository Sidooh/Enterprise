package api

import (
	"enterprise.sidooh/api/middleware"
	"enterprise.sidooh/api/middleware/jwt"
	"enterprise.sidooh/api/routes"
	"enterprise.sidooh/pkg/clients"
	"enterprise.sidooh/pkg/logger"
	"enterprise.sidooh/pkg/services/account"
	"enterprise.sidooh/pkg/services/auth"
	"enterprise.sidooh/pkg/services/enterprise"
	"enterprise.sidooh/pkg/services/team"
	"enterprise.sidooh/pkg/services/user"
	"enterprise.sidooh/utils"
	"errors"
	"fmt"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/helmet/v2"
	"github.com/spf13/viper"
	"time"
)

func setMiddleware(app *fiber.App) {
	//app.Use(func(c *fiber.Ctx) error {
	//	// TODO: URGENT: Check out these headers and review them
	//	// Set some security headers:
	//	c.Set("X-XSS-Protection", "1; mode=block")
	//	c.Set("X-Content-Type-Options", "nosniff")
	//	c.Set("X-Download-Options", "noopen")
	//	c.Set("Strict-Transport-Security", "max-age=5184000")
	//	c.Set("X-Frame-Options", "SAMEORIGIN")
	//	c.Set("X-DNS-Prefetch-Control", "off")
	//
	//	// Go to next middleware:
	//	return c.Next()
	//})

	app.Use(helmet.New())
	app.Use(cors.New())
	app.Use(limiter.New())
	app.Use(recover.New())
	app.Use(fiberLogger.New(fiberLogger.Config{Output: utils.GetLogFile("stats.log")}))

	app.Use(favicon.New(favicon.Config{Next: func(c *fiber.Ctx) bool {
		return true
	}}))

	middleware.Validator = validator.New()

}

func setHealthCheckRoutes(app *fiber.App) {
	app.Get("/200", func(ctx *fiber.Ctx) error {
		return ctx.JSON("200")
	})
}

func setHandlers(app *fiber.App) {
	api := app.Group("/api")
	v1 := api.Group("/v1")

	// Initialize rest clients
	clients.InitAccountClient()
	clients.InitPaymentClient()
	clients.InitNotifyClient()

	userRep := user.NewRepo()

	authSrv := auth.NewService(userRep)

	userSrv := user.NewService(userRep)

	enterpriseRep := enterprise.NewRepo()
	enterpriseSrv := enterprise.NewService(enterpriseRep)

	accountRep := account.NewRepo()
	accountSrv := account.NewService(accountRep)

	teamRep := team.NewRepo()
	teamSrv := team.NewService(teamRep)

	routes.AuthRouter(v1, authSrv)

	app.Use(jwt.New(jwt.Config{
		Secret: viper.GetString("JWT_KEY"),
		Expiry: time.Duration(15) * time.Minute,
	}))

	routes.EnterpriseRouter(v1, enterpriseSrv)
	routes.UserRouter(v1, userSrv)
	routes.AccountRouter(v1, accountSrv)
	routes.TeamRouter(v1, teamSrv)
}

func Server() *fiber.App {
	// Create a new fiber instance with custom config
	app := fiber.New(fiber.Config{
		Prefork: viper.GetBool("PREFORK"),

		// Override default error handler
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			logger.ServerLog.Error(err)
			// Status code defaults to 500
			code := fiber.StatusInternalServerError

			// Retrieve the custom status code if it's a *fiber.Error
			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}

			// Send custom error page
			err = ctx.Status(code).SendFile(fmt.Sprintf("./%d.html", code))
			if err != nil {
				// In case the SendFile fails
				return ctx.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
			}

			// Return from handler
			return nil
		},
	})

	// ...

	setMiddleware(app)
	setHealthCheckRoutes(app)
	setHandlers(app)

	return app
}
