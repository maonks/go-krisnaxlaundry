package tracking

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func TrackingRoute(
	app *fiber.App,
	db *gorm.DB,
) {

	api := app.Group("/api/public")

	api.Get(
		"/tracking/:order_no",
		Detail(db),
	)
}
