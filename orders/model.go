package orders

import "time"

type Order struct {
	ID                int64      `gorm:"column:id" json:"id"`
	OrderNo           string     `gorm:"column:order_no" json:"order_no"`
	CustomerID        int64      `gorm:"column:customer_id" json:"customer_id"`
	CurrentStatusID   int64      `gorm:"column:current_status_id" json:"current_status_id"`
	OrderDate         time.Time  `gorm:"column:order_date" json:"order_date"`
	EstimatedFinishAt *time.Time `gorm:"column:estimated_finish_at" json:"estimated_finish_at"`

	TotalQty       float64 `gorm:"column:total_qty" json:"total_qty"`
	GrossAmount    float64 `gorm:"column:gross_amount" json:"gross_amount"`
	DiscountAmount float64 `gorm:"column:discount_amount" json:"discount_amount"`
	NetAmount      float64 `gorm:"column:net_amount" json:"net_amount"`
	PaidAmount     float64 `gorm:"column:paid_amount" json:"paid_amount"`

	Notes string `gorm:"column:notes" json:"notes"`

	CreatedBy int64 `gorm:"column:created_by" json:"created_by"`

	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (Order) TableName() string {
	return "laundry_orders"
}

type OrderDetail struct {
	ID int64 `gorm:"column:id"`

	OrderID int64 `gorm:"column:order_id"`

	ServiceID int64 `gorm:"column:service_id"`

	Qty float64 `gorm:"column:qty"`

	UnitPrice float64 `gorm:"column:unit_price"`

	Subtotal float64 `gorm:"column:subtotal"`

	Notes string `gorm:"column:notes"`
}

func (OrderDetail) TableName() string {
	return "laundry_order_details"
}
