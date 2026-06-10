package middlewares

import (
	"fmt"
	"strings"

	"github.com/maonkscode/go-kresnaxlaundry/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func Auth() fiber.Handler {

	return func(c *fiber.Ctx) error {

		authHeader := c.Get("Authorization")

		fmt.Println("AUTH:", authHeader)

		if authHeader == "" {

			return c.Status(401).JSON(
				fiber.Map{
					"success": false,
					"msg":     "unauthorized",
				},
			)
		}

		tokenString := strings.ReplaceAll(
			authHeader,
			"Bearer ",
			"",
		)

		token, err := utils.ParseJWT(
			tokenString,
		)

		if err != nil || !token.Valid {

			return c.Status(401).JSON(
				fiber.Map{
					"success": false,
					"msg":     "invalid token",
				},
			)
		}

		claims := token.Claims.(jwt.MapClaims)

		c.Locals(
			"user_id",
			int64(claims["user_id"].(float64)),
		)

		c.Locals(
			"role_id",
			int64(claims["role_id"].(float64)),
		)

		return c.Next()
	}
}
