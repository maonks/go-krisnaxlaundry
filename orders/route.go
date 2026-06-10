package orders

import (
	"github.com/maonkscode/go-kresnaxlaundry/middlewares"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func OrderRoute(
	app *fiber.App,
	db *gorm.DB,
) {

	api := app.Group("/api/orders", middlewares.Auth())

	api.Post("/", Create(db))
	api.Get("/", Index(db))
	api.Get("/", Index(db))

	api.Get("/:id", Detail(db))

	api.Post("/", Create(db))

	api.Put("/:id/status", UpdateStatus(db))

	api.Put("/:id/cancel", Cancel(db))
}
