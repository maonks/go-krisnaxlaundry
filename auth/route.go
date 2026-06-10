package auth

import (
	"github.com/maonkscode/go-kresnaxlaundry/middlewares"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func AuthRoute(
	app *fiber.App,
	db *gorm.DB,
) {

	api := app.Group("/api")

	api.Post(
		"/auth/google",
		GoogleLogin(db),
	)

	api.Get(
		"/auth/me",
		middlewares.Auth(),
		Me(db),
	)
}
