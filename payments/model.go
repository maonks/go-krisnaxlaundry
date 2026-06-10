package payments

import "time"

type Payment struct {
	ID int64 `gorm:"column:id"`

	OrderID int64 `gorm:"column:order_id"`

	PaymentMethodID int64 `gorm:"column:payment_method_id"`

	Amount float64 `gorm:"column:amount"`

	PaidAt time.Time `gorm:"column:paid_at"`

	PaidBy int64 `gorm:"column:paid_by"`

	ReferenceNo string `gorm:"column:reference_no"`

	Notes string `gorm:"column:notes"`

	CreatedAt time.Time `gorm:"column:created_at"`
}

func (Payment) TableName() string {
	return "payments"
}
