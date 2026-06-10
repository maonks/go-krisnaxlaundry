package orders

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type OrderSequence struct {
	OrderDate  time.Time `gorm:"column:order_date"`
	LastNumber int64     `gorm:"column:last_number"`
}

func (OrderSequence) TableName() string {
	return "order_sequences"
}

func GenerateOrderNumber(
	tx *gorm.DB,
) (string, error) {

	today := time.Now()

	dateOnly := today.Format("2006-01-02")

	var seq OrderSequence

	err := tx.Raw(`
		SELECT
			order_date,
			last_number
		FROM order_sequences
		WHERE order_date = ?
		FOR UPDATE
	`, dateOnly).
		Scan(&seq).
		Error

	if err != nil {
		return "", err
	}

	if seq.OrderDate.IsZero() {

		err = tx.Exec(`
			INSERT INTO order_sequences
			(
				order_date,
				last_number
			)
			VALUES
			(
				?,
				1
			)
		`, dateOnly).
			Error

		if err != nil {
			return "", err
		}

		return fmt.Sprintf(
			"KXL-%s-%05d",
			today.Format("20060102"),
			1,
		), nil
	}

	nextNumber := seq.LastNumber + 1

	err = tx.Exec(`
		UPDATE order_sequences
		SET last_number = ?
		WHERE order_date = ?
	`,
		nextNumber,
		dateOnly,
	).Error

	if err != nil {
		return "", err
	}

	return fmt.Sprintf(
		"KXL-%s-%05d",
		today.Format("20060102"),
		nextNumber,
	), nil
}

type CreateOrderRequest struct {
	CustomerID     int64   `json:"customer_id"`
	DiscountAmount float64 `json:"discount_amount"`
	Notes          string  `json:"notes"`

	Items []CreateOrderItem `json:"items"`
}

type CreateOrderItem struct {
	ServiceID int64   `json:"service_id"`
	Qty       float64 `json:"qty"`
}

type LaundryStatus struct {
	ID   int64
	Code string
	Name string
}

func (LaundryStatus) TableName() string {
	return "laundry_statuses"
}

func GetReceivedStatusID(
	tx *gorm.DB,
) (int64, error) {

	var status LaundryStatus

	err := tx.
		Where("code = ?", "RECEIVED").
		First(&status).
		Error

	if err != nil {
		return 0, err
	}

	return status.ID, nil
}

type ServicePrice struct {
	ID    int64
	Price float64
}

func GetServicePrice(
	tx *gorm.DB,
	serviceID int64,
) (*ServicePrice, error) {

	var service ServicePrice

	err := tx.Raw(`
		SELECT
			id,
			price
		FROM services
		WHERE id = ?
	`, serviceID).
		Scan(&service).
		Error

	if err != nil {
		return nil, err
	}

	return &service, nil
}

func CreateOrder(
	db *gorm.DB,
	req CreateOrderRequest,
	userID int64,
) error {

	return db.Transaction(func(tx *gorm.DB) error {

		orderNo, err := GenerateOrderNumber(tx)

		if err != nil {
			return err
		}

		statusID, err := GetReceivedStatusID(tx)

		if err != nil {
			return err
		}

		var grossAmount float64
		var totalQty float64

		order := Order{
			OrderNo:         orderNo,
			CustomerID:      req.CustomerID,
			CurrentStatusID: statusID,
			OrderDate:       time.Now(),
			DiscountAmount:  req.DiscountAmount,
			CreatedBy:       userID,
			Notes:           req.Notes,
		}

		if err := tx.Create(&order).Error; err != nil {
			return err
		}

		for _, item := range req.Items {

			service, err := GetServicePrice(
				tx,
				item.ServiceID,
			)

			if err != nil {
				return err
			}

			subtotal := item.Qty * service.Price

			detail := OrderDetail{
				OrderID:   order.ID,
				ServiceID: item.ServiceID,
				Qty:       item.Qty,
				UnitPrice: service.Price,
				Subtotal:  subtotal,
			}

			if err := tx.Create(&detail).Error; err != nil {
				return err
			}

			grossAmount += subtotal
			totalQty += item.Qty
		}

		netAmount := grossAmount - req.DiscountAmount

		if err := tx.Model(&Order{}).
			Where("id = ?", order.ID).
			Updates(map[string]interface{}{
				"gross_amount": grossAmount,
				"net_amount":   netAmount,
				"total_qty":    totalQty,
			}).Error; err != nil {

			return err
		}

		if err := tx.Exec(`
			INSERT INTO laundry_status_histories
			(
				order_id,
				status_id,
				changed_by
			)
			VALUES
			(
				?,
				?,
				?
			)
		`,
			order.ID,
			statusID,
			userID,
		).Error; err != nil {

			return err
		}

		return nil
	})
}

type OrderListRow struct {
	ID int64 `json:"id"`

	OrderNo string `json:"order_no"`

	CustomerName string `json:"customer_name"`

	CustomerPhone string `json:"customer_phone"`

	StatusID int64 `json:"status_id"`

	StatusName string `json:"status_name"`

	NetAmount float64 `json:"net_amount"`

	OrderDate string `json:"order_date"`
}

func GetOrders(
	db *gorm.DB,
	search string,
	status string,
	dateFrom string,
	dateTo string,
	start int,
	limit int,
) ([]OrderListRow, int64, error) {

	var rows []OrderListRow

	query := `
	SELECT
		o.id,
		o.order_no,

		c.name AS customer_name,
		c.phone AS customer_phone,

		s.id AS status_id,
		s.name AS status_name,

		o.net_amount,

		TO_CHAR(
			o.order_date,
			'YYYY-MM-DD HH24:MI'
		) AS order_date

	FROM laundry_orders o

	INNER JOIN customers c
		ON c.id = o.customer_id

	INNER JOIN laundry_statuses s
		ON s.id = o.current_status_id

	WHERE 1=1
	`

	args := []interface{}{}

	if search != "" {

		query += `
		AND (
			c.name ILIKE ?
			OR c.phone ILIKE ?
			OR o.order_no ILIKE ?
		)
		`

		searchLike := "%" + search + "%"

		args = append(
			args,
			searchLike,
			searchLike,
			searchLike,
		)
	}

	if status != "" {

		query += `
		AND s.id = ?
		`

		args = append(
			args,
			status,
		)
	}

	if dateFrom != "" {

		query += `
		AND DATE(o.order_date) >= ?
		`

		args = append(
			args,
			dateFrom,
		)
	}

	if dateTo != "" {

		query += `
		AND DATE(o.order_date) <= ?
		`

		args = append(
			args,
			dateTo,
		)
	}

	countQuery := `
	SELECT COUNT(*)
	FROM (
	` + query + `
	) x
	`

	var total int64

	if err := db.
		Raw(
			countQuery,
			args...,
		).
		Scan(&total).
		Error; err != nil {

		return nil, 0, err
	}

	query += `
	ORDER BY o.id DESC
	LIMIT ?
	OFFSET ?
	`

	args = append(
		args,
		limit,
		start,
	)

	if err := db.
		Raw(
			query,
			args...,
		).
		Scan(&rows).
		Error; err != nil {

		return nil, 0, err
	}

	return rows, total, nil
}

type OrderDetailResponse struct {
	ID int64 `json:"id"`

	OrderNo string `json:"order_no"`

	CustomerName string `json:"customer_name"`

	CustomerPhone string `json:"customer_phone"`

	StatusName string `json:"status_name"`

	GrossAmount float64 `json:"gross_amount"`

	DiscountAmount float64 `json:"discount_amount"`

	NetAmount float64 `json:"net_amount"`

	OrderDate string `json:"order_date"`

	Notes string `json:"notes"`

	Items []OrderItemResponse `json:"items"`
}

type OrderItemResponse struct {
	ServiceName string `json:"service_name"`

	Qty float64 `json:"qty"`

	UnitPrice float64 `json:"unit_price"`

	Subtotal float64 `json:"subtotal"`
}

func GetOrderDetail(
	db *gorm.DB,
	id int64,
) (*OrderDetailResponse, error) {

	var result OrderDetailResponse

	err := db.Raw(`
	SELECT
		o.id,
		o.order_no,

		c.name customer_name,
		c.phone customer_phone,

		s.name status_name,

		o.gross_amount,
		o.discount_amount,
		o.net_amount,

		o.notes,

		TO_CHAR(
			o.order_date,
			'YYYY-MM-DD HH24:MI'
		) order_date

	FROM laundry_orders o

	INNER JOIN customers c
		ON c.id=o.customer_id

	INNER JOIN laundry_statuses s
		ON s.id=o.current_status_id

	WHERE o.id=?
	`, id).
		Scan(&result).
		Error

	if err != nil {
		return nil, err
	}

	var items []OrderItemResponse

	db.Raw(`
	SELECT

		s.name service_name,

		d.qty,

		d.unit_price,

		d.subtotal

	FROM laundry_order_details d

	INNER JOIN services s
		ON s.id=d.service_id

	WHERE d.order_id=?
	`, id).
		Scan(&items)

	result.Items = items

	return &result, nil
}

type UpdateStatusRequest struct {
	StatusID int64 `json:"status_id"`

	Remarks string `json:"remarks"`
}

func UpdateOrderStatus(
	db *gorm.DB,
	orderID int64,
	statusID int64,
	remarks string,
	userID int64,
) error {

	return db.Transaction(func(tx *gorm.DB) error {

		err := tx.Exec(`
			UPDATE laundry_orders
			SET current_status_id = ?
			WHERE id = ?
		`,
			statusID,
			orderID,
		).Error

		if err != nil {
			return err
		}

		err = tx.Exec(`
			INSERT INTO laundry_status_histories
			(
				order_id,
				status_id,
				remarks,
				changed_by
			)
			VALUES
			(
				?,
				?,
				?,
				?
			)
		`,
			orderID,
			statusID,
			remarks,
			userID,
		).Error

		return err
	})
}

func GetCancelledStatusID(
	tx *gorm.DB,
) (int64, error) {

	var status LaundryStatus

	err := tx.
		Where("code = ?", "CANCELLED").
		First(&status).
		Error

	if err != nil {
		return 0, err
	}

	return status.ID, nil
}

func CancelOrder(
	db *gorm.DB,
	orderID int64,
	remarks string,
	userID int64,
) error {

	return db.Transaction(func(tx *gorm.DB) error {

		cancelStatusID, err := GetCancelledStatusID(tx)

		if err != nil {
			return err
		}

		err = tx.Exec(`
			UPDATE laundry_orders
			SET
				current_status_id = ?,
				updated_at = NOW()
			WHERE id = ?
		`,
			cancelStatusID,
			orderID,
		).Error

		if err != nil {
			return err
		}

		err = tx.Exec(`
			INSERT INTO laundry_status_histories
			(
				order_id,
				status_id,
				remarks,
				changed_by
			)
			VALUES
			(
				?,
				?,
				?,
				?
			)
		`,
			orderID,
			cancelStatusID,
			remarks,
			userID,
		).Error

		if err != nil {
			return err
		}

		return nil
	})
}
