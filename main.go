package main

import (
	"log"
	"os"

	"github.com/maonkscode/go-kresnaxlaundry/auth"
	"github.com/maonkscode/go-kresnaxlaundry/configs"
	"github.com/maonkscode/go-kresnaxlaundry/customers"
	"github.com/maonkscode/go-kresnaxlaundry/dashboard"
	"github.com/maonkscode/go-kresnaxlaundry/orders"
	"github.com/maonkscode/go-kresnaxlaundry/payments"
	"github.com/maonkscode/go-kresnaxlaundry/photos"
	"github.com/maonkscode/go-kresnaxlaundry/tracking"
	"github.com/maonkscode/go-kresnaxlaundry/users"

	"github.com/gofiber/fiber/v2"
)

func main() {

	configs.LoadEnv()

	db, err := configs.KonekDB()

	if err != nil {
		log.Fatal(err)
	}

	app := fiber.New()

	app.Static(
		"/storage",
		"./storage",
	)

	auth.AuthRoute(app, db)

	users.UserRoute(app, db)

	customers.CustomerRoute(app, db)

	orders.OrderRoute(app, db)

	payments.PaymentRoute(app, db)

	photos.PhotoRoute(app, db)

	tracking.TrackingRoute(app, db)

	dashboard.DashboardRoute(app, db)

	log.Fatal(
		app.Listen(
			":" + os.Getenv("APP_PORT"),
		),
	)
}
