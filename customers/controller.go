package customers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/maonkscode/go-kresnaxlaundry/utils"
	"gorm.io/gorm"
)

func Index(db *gorm.DB) fiber.Handler {

	return func(c *fiber.Ctx) error {

		search := c.Query("search")

		data, err := GetCustomers(
			db,
			search,
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

func SearchPhone(db *gorm.DB) fiber.Handler {

	return func(c *fiber.Ctx) error {

		phone := c.Query("phone")

		if phone == "" {

			return utils.Error(
				c,
				"phone required",
			)
		}

		customer, err := FindByPhone(
			db,
			phone,
		)

		if err != nil {

			return utils.Error(
				c,
				"customer not found",
			)
		}

		return utils.Success(
			c,
			customer,
		)
	}
}

func Detail(db *gorm.DB) fiber.Handler {

	return func(c *fiber.Ctx) error {

		id := c.Params("id")

		var customer Customer

		if err := db.
			Where("id = ?", id).
			First(&customer).
			Error; err != nil {

			return utils.Error(
				c,
				"customer not found",
			)
		}

		return utils.Success(
			c,
			customer,
		)
	}
}

func Create(db *gorm.DB) fiber.Handler {

	return func(c *fiber.Ctx) error {

		var customer Customer

		if err := c.BodyParser(&customer); err != nil {

			return utils.Error(
				c,
				"invalid request",
			)
		}

		if err := CreateCustomer(
			db,
			&customer,
		); err != nil {

			return utils.Error(
				c,
				err.Error(),
			)
		}

		return utils.Success(
			c,
			customer,
		)
	}
}

func Update(db *gorm.DB) fiber.Handler {

	return func(c *fiber.Ctx) error {

		id := c.Params("id")

		var customer Customer

		if err := c.BodyParser(
			&customer,
		); err != nil {

			return utils.Error(
				c,
				"invalid request",
			)
		}

		data := map[string]interface{}{
			"name":    customer.Name,
			"phone":   customer.Phone,
			"email":   customer.Email,
			"address": customer.Address,
			"notes":   customer.Notes,
		}

		if err := UpdateCustomer(
			db,
			utils.StringToInt64(id),
			data,
		); err != nil {

			return utils.Error(
				c,
				err.Error(),
			)
		}

		return utils.Success(
			c,
			"customer updated",
		)
	}
}

func Delete(db *gorm.DB) fiber.Handler {

	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		err := DeleteCustomer(
			db,
			utils.StringToInt64(id),
		)

		if err != nil {

			return utils.Error(
				c,
				err.Error(),
			)
		}

		return utils.Success(
			c,
			"customer deleted",
		)
	}
}
