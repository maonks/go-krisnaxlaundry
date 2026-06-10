package customers

import (
	"github.com/maonkscode/go-kresnaxlaundry/middlewares"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CustomerRoute(
	app *fiber.App,
	db *gorm.DB,
) {

	api := app.Group(
		"/api/customers",
		middlewares.Auth(),
	)

	api.Get("/", Index(db))

	api.Get("/search", SearchPhone(db))

	api.Post("/", Create(db))

	api.Put("/:id", Update(db))

	api.Delete("/:id", Delete(db))
}
