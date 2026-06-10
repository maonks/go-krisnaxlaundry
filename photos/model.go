package photos

import "time"

type OrderPhoto struct {
	ID int64 `gorm:"column:id"`

	OrderID int64 `gorm:"column:order_id"`

	FileName string `gorm:"column:file_name"`

	FilePath string `gorm:"column:file_path"`

	CreatedAt time.Time `gorm:"column:created_at"`
}

func (OrderPhoto) TableName() string {
	return "order_photos"
}
