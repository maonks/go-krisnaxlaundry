package orders

import (
	"github.com/gofiber/fiber/v2"
	"github.com/maonkscode/go-kresnaxlaundry/utils"
	"gorm.io/gorm"
)

func Create(db *gorm.DB) fiber.Handler {

	return func(c *fiber.Ctx) error {

		var req CreateOrderRequest

		if err := c.BodyParser(&req); err != nil {

			return utils.Error(
				c,
				"invalid request",
			)
		}

		userID := c.Locals("user_id").(int64)

		if err := CreateOrder(
			db,
			req,
			userID,
		); err != nil {

			return utils.Error(
				c,
				err.Error(),
			)
		}

		return utils.Success(
			c,
			"order created",
		)
	}
}

func Index(db *gorm.DB) fiber.Handler {

	return func(c *fiber.Ctx) error {

		search := c.Query("search")
		status := c.Query("status")

		dateFrom := c.Query("date_from")
		dateTo := c.Query("date_to")

		page, limit, start := utils.GetPagination(c)

		rows, total, err := GetOrders(
			db,
			search,
			status,
			dateFrom,
			dateTo,
			start,
			limit,
		)

		if err != nil {

			return utils.Error(
				c,
				err.Error(),
			)
		}

		return c.JSON(fiber.Map{
			"success": true,
			"count":   len(rows),

			"count_total": total,

			"page": page,

			"data": rows,
		})
	}
}

func Detail(db *gorm.DB) fiber.Handler {

	return func(c *fiber.Ctx) error {

		id := utils.StringToInt64(
			c.Params("id"),
		)

		data, err := GetOrderDetail(
			db,
			id,
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

func UpdateStatus(
	db *gorm.DB,
) fiber.Handler {

	return func(c *fiber.Ctx) error {

		id := utils.StringToInt64(
			c.Params("id"),
		)

		userID := c.Locals(
			"user_id",
		).(int64)

		var req UpdateStatusRequest

		if err := c.BodyParser(
			&req,
		); err != nil {

			return utils.Error(
				c,
				"invalid request",
			)
		}

		err := UpdateOrderStatus(
			db,
			id,
			req.StatusID,
			req.Remarks,
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
			"status updated",
		)
	}
}

type CancelOrderRequest struct {
	Remarks string `json:"remarks"`
}

func Cancel(
	db *gorm.DB,
) fiber.Handler {

	return func(c *fiber.Ctx) error {

		orderID := utils.StringToInt64(
			c.Params("id"),
		)

		userID := c.Locals(
			"user_id",
		).(int64)

		var req CancelOrderRequest

		if err := c.BodyParser(
			&req,
		); err != nil {

			return utils.Error(
				c,
				"invalid request",
			)
		}

		err := CancelOrder(
			db,
			orderID,
			req.Remarks,
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
			"order cancelled",
		)
	}
}
