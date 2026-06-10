package utils

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

func GenerateOrderNumber(
	db *gorm.DB,
) (string, error) {

	today := time.Now().Format("20060102")

	var total int64

	err := db.Raw(`
		SELECT COUNT(*)
		FROM laundry_orders
		WHERE DATE(created_at)=CURRENT_DATE
	`).Scan(&total).Error

	if err != nil {
		return "", err
	}

	return fmt.Sprintf(
		"KXL-%s-%05d",
		today,
		total+1,
	), nil
}
