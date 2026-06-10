package customers

import "time"

type Customer struct {
	ID        int64     `gorm:"column:id;primaryKey" json:"id"`
	Name      string    `gorm:"column:name" json:"name"`
	Phone     string    `gorm:"column:phone" json:"phone"`
	Email     *string   `gorm:"column:email" json:"email"`
	Address   *string   `gorm:"column:address" json:"address"`
	Notes     *string   `gorm:"column:notes" json:"notes"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (Customer) TableName() string {
	return "customers"
}
