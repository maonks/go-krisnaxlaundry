package users

import (
	"github.com/maonkscode/go-kresnaxlaundry/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Index(
	db *gorm.DB,
) fiber.Handler {

	return func(c *fiber.Ctx) error {

		data, err := GetUsers(db)

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
