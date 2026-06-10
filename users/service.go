package users

import "gorm.io/gorm"

func GetUsers(db *gorm.DB) ([]User, error) {

	var users []User

	err := db.
		Order("id DESC").
		Find(&users).
		Error

	return users, err
}
