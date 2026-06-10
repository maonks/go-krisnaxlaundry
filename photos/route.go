package photos

import (
	"github.com/gofiber/fiber/v2"
	"github.com/maonkscode/go-kresnaxlaundry/middlewares"
	"gorm.io/gorm"
)

func PhotoRoute(
	app *fiber.App,
	db *gorm.DB,
) {

	api := app.Group(
		"/api/photos",
		middlewares.Auth(),
	)

	api.Post(
		"/order/:id",
		Upload(db),
	)

	api.Get(
		"/order/:id",
		List(db),
	)

	api.Delete(
		"/:id",
		Delete(db),
	)
}
