package photos

import "gorm.io/gorm"

func GetOrderPhotos(
	db *gorm.DB,
	orderID int64,
) ([]PhotoResponse, error) {

	var rows []PhotoResponse

	err := db.Raw(`
		SELECT
			id,
			file_name,
			'/storage/orders/' ||
			order_id ||
			'/' ||
			file_name
			AS file_url
		FROM order_photos
		WHERE order_id = ?
		ORDER BY id DESC
	`,
		orderID,
	).
		Scan(&rows).
		Error

	return rows, err
}

type PhotoResponse struct {
	ID       int64  `json:"id"`
	FileName string `json:"file_name"`
	FileURL  string `json:"file_url"`
}

func GetPhotoByID(
	db *gorm.DB,
	id int64,
) (*OrderPhoto, error) {

	var photo OrderPhoto

	err := db.
		Where("id = ?", id).
		First(&photo).
		Error

	if err != nil {
		return nil, err
	}

	return &photo, nil
}

func DeletePhoto(
	db *gorm.DB,
	id int64,
) error {

	return db.
		Delete(
			&OrderPhoto{},
			id,
		).
		Error
}
