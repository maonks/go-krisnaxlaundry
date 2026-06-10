package photos

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/maonkscode/go-kresnaxlaundry/utils"
	"gorm.io/gorm"
)

func Upload(
	db *gorm.DB,
) fiber.Handler {

	return func(c *fiber.Ctx) error {

		orderID := utils.StringToInt64(
			c.Params("id"),
		)

		file, err := c.FormFile(
			"photo",
		)

		if err != nil {

			return utils.Error(
				c,
				"photo required",
			)
		}

		dir := fmt.Sprintf(
			"./storage/orders/%d",
			orderID,
		)

		os.MkdirAll(
			dir,
			os.ModePerm,
		)

		fileName := fmt.Sprintf(
			"%d_%s",
			time.Now().Unix(),
			file.Filename,
		)

		path := dir + "/" + fileName

		if err := c.SaveFile(
			file,
			path,
		); err != nil {

			return utils.Error(
				c,
				err.Error(),
			)
		}

		photo := OrderPhoto{
			OrderID: orderID,

			FileName: fileName,

			FilePath: path,
		}

		if err := db.Create(
			&photo,
		).Error; err != nil {

			return utils.Error(
				c,
				err.Error(),
			)
		}

		return utils.Success(
			c,
			photo,
		)
	}
}

func List(
	db *gorm.DB,
) fiber.Handler {

	return func(c *fiber.Ctx) error {

		orderID := utils.StringToInt64(
			c.Params("id"),
		)

		data, err := GetOrderPhotos(
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

func Delete(
	db *gorm.DB,
) fiber.Handler {

	return func(c *fiber.Ctx) error {

		photoID := utils.StringToInt64(
			c.Params("id"),
		)

		photo, err := GetPhotoByID(
			db,
			photoID,
		)

		if err != nil {

			return utils.Error(
				c,
				"photo not found",
			)
		}

		// hapus file fisik

		if photo.FilePath != "" {

			_ = os.Remove(
				photo.FilePath,
			)
		}

		// hapus record database

		if err := DeletePhoto(
			db,
			photoID,
		); err != nil {

			return utils.Error(
				c,
				err.Error(),
			)
		}

		return utils.Success(
			c,
			"photo deleted",
		)
	}
}
