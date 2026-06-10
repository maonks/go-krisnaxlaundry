package tracking

import (
	"github.com/maonkscode/go-kresnaxlaundry/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Detail(
	db *gorm.DB,
) fiber.Handler {

	return func(c *fiber.Ctx) error {

		orderNo := c.Params("order_no")

		data, err := GetTracking(
			db,
			orderNo,
		)

		if err != nil {

			return utils.Error(
				c,
				"order not found",
			)
		}

		return utils.Success(
			c,
			data,
		)
	}
}
