package users

import "time"

type User struct {
	ID          int64      `gorm:"column:id;primaryKey"`
	GoogleID    *string    `gorm:"column:google_id"`
	Email       string     `gorm:"column:email"`
	FullName    string     `gorm:"column:full_name"`
	RoleID      int64      `gorm:"column:role_id"`
	IsActive    bool       `gorm:"column:is_active"`
	LastLoginAt *time.Time `gorm:"column:last_login_at"`

	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (User) TableName() string {
	return "users"
}
