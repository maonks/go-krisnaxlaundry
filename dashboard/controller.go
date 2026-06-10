package dashboard

import (
	"github.com/maonkscode/go-kresnaxlaundry/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SummaryDashboard(db *gorm.DB) fiber.Handler {

	return func(c *fiber.Ctx) error {

		data, err := GetSummary(
			db,
		)

		if err != nil {

			return utils.Error(
				c,
				err.Error(),
			)
		}

		return utils.Success(
			c,
			data,
		)
	}
}

func RecentOrders(
	db *gorm.DB,
) fiber.Handler {

	return func(c *fiber.Ctx) error {

		data, err := GetRecentOrders(
			db,
		)

		if err != nil {

			return utils.Error(
				c,
				err.Error(),
			)
		}

		return utils.Success(
			c,
			data,
		)
	}
}
