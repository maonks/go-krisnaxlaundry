package payments

import (
	"github.com/gofiber/fiber/v2"
	"github.com/maonkscode/go-kresnaxlaundry/utils"
	"gorm.io/gorm"
)

func Create(
	db *gorm.DB,
) fiber.Handler {

	return func(c *fiber.Ctx) error {

		userID := c.Locals(
			"user_id",
		).(int64)

		var req CreatePaymentRequest

		if err := c.BodyParser(
			&req,
		); err != nil {

			return utils.Error(
				c,
				"invalid request",
			)
		}

		err := CreatePayment(
			db,
			req,
			userID,
		)

		if err != nil {

			return utils.Error(
				c,
				err.Error(),
			)
		}

		return utils.Success(
			c,
			"payment success",
		)
	}
}

func History(
	db *gorm.DB,
) fiber.Handler {

	return func(c *fiber.Ctx) error {

		orderID := utils.StringToInt64(
			c.Params("id"),
		)

		data, err := GetPaymentHistory(
			db,
			orderID,
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
