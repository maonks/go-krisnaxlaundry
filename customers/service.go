package customers

import (
	"errors"

	"gorm.io/gorm"
)

func GetCustomers(db *gorm.DB, search string) ([]Customer, error) {

	var data []Customer

	query := db.Model(&Customer{})

	if search != "" {

		query = query.Where(
			"name ILIKE ? OR phone ILIKE ?",
			"%"+search+"%",
			"%"+search+"%",
		)
	}

	err := query.
		Order("id DESC").
		Find(&data).
		Error

	return data, err
}

func FindByPhone(db *gorm.DB, phone string) (*Customer, error) {

	var customer Customer

	err := db.
		Where("phone = ?", phone).
		First(&customer).
		Error

	if err != nil {
		return nil, err
	}

	return &customer, nil
}

func CreateCustomer(db *gorm.DB, customer *Customer) error {

	var total int64

	db.Model(&Customer{}).
		Where("phone = ?", customer.Phone).
		Count(&total)

	if total > 0 {
		return errors.New("phone already registered")
	}

	return db.Create(customer).Error
}

func UpdateCustomer(db *gorm.DB, id int64, data map[string]interface{}) error {

	return db.
		Model(&Customer{}).
		Where("id = ?", id).
		Updates(data).
		Error
}

func DeleteCustomer(db *gorm.DB, id int64) error {

	return db.
		Delete(
			&Customer{},
			id,
		).
		Error
}
