package payments

import (
	"github.com/maonkscode/go-kresnaxlaundry/middlewares"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func PaymentRoute(
	app *fiber.App,
	db *gorm.DB,
) {

	api := app.Group(
		"/api/payments",
		middlewares.Auth(),
	)

	api.Post(
		"/",
		Create(db),
	)

	api.Get(
		"/history/:id",
		History(db),
	)
}
