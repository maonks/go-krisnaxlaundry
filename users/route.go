package users

import (
	"github.com/maonkscode/go-kresnaxlaundry/middlewares"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func UserRoute(
	app *fiber.App,
	db *gorm.DB,
) {

	api := app.Group(
		"/api/users",
		middlewares.Auth(),
	)

	api.Get(
		"/",
		Index(db),
	)
}
