package dashboard

import (
	"github.com/maonkscode/go-kresnaxlaundry/middlewares"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func DashboardRoute(
	app *fiber.App,
	db *gorm.DB,
) {

	api := app.Group(
		"/api/dashboard",
		middlewares.Auth(),
	)

	api.Get(
		"/summary",
		SummaryDashboard(db),
	)

	api.Get(
		"/recent-orders",
		RecentOrders(db),
	)
}
