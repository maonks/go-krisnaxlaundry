package utils

import (
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetPagination(c *fiber.Ctx) (int, int, int) {

	pageStr := c.Query("page", "1")

	page, err := strconv.Atoi(pageStr)

	if err != nil || page < 1 {
		page = 1
	}

	limitStr := os.Getenv("APP_LOADDATA")

	limit, err := strconv.Atoi(limitStr)

	if err != nil || limit <= 0 {
		limit = 10
	}

	start := (page - 1) * limit

	return page, limit, start
}
