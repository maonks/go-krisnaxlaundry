package middlewares

import (
	"github.com/gofiber/fiber/v2"
)

func Role(
	roleIDs ...int64,
) fiber.Handler {

	return func(c *fiber.Ctx) error {

		currentRole := c.Locals(
			"role_id",
		).(int64)

		for _, role := range roleIDs {

			if currentRole == role {
				return c.Next()
			}
		}

		return c.Status(
			fiber.StatusForbidden,
		).JSON(
			fiber.Map{
				"success": false,
				"msg":     "access denied",
			},
		)
	}
}
