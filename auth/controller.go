package auth

import (
	"github.com/maonkscode/go-kresnaxlaundry/users"
	"github.com/maonkscode/go-kresnaxlaundry/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GoogleLogin(db *gorm.DB) fiber.Handler {

	return func(c *fiber.Ctx) error {

		type Request struct {
			IDToken string `json:"id_token"`
		}

		var req Request

		if err := c.BodyParser(&req); err != nil {
			return utils.Error(
				c,
				"invalid request",
			)
		}

		user, err := LoginWithGoogle(
			db,
			req.IDToken,
		)

		if err != nil {
			return utils.Error(
				c,
				"unauthorized",
			)
		}

		token, err := utils.GenerateJWT(
			user.ID,
			user.RoleID,
		)
		if err != nil {
			return utils.Error(
				c,
				err.Error(),
			)
		}

		return utils.Success(
			c,
			fiber.Map{
				"token": token,
				"user":  user,
			},
		)
	}
}

func Me(db *gorm.DB) fiber.Handler {

	return func(c *fiber.Ctx) error {

		userID := c.Locals("user_id")

		var user users.User

		if err := db.
			First(
				&user,
				userID,
			).Error; err != nil {

			return utils.Error(
				c,
				"user not found",
			)
		}

		return utils.Success(
			c,
			user,
		)
	}
}
