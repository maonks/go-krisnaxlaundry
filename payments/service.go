package payments

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type CreatePaymentRequest struct {
	OrderID int64 `json:"order_id"`

	PaymentMethodID int64 `json:"payment_method_id"`

	ReferenceNo string `json:"reference_no"`

	Notes string `json:"notes"`
}
type OrderPaymentInfo struct {
	ID int64

	NetAmount float64

	PaidAmount float64
}

func GetOrderPaymentInfo(
	tx *gorm.DB,
	orderID int64,
) (*OrderPaymentInfo, error) {

	var row OrderPaymentInfo

	err := tx.Raw(`
		SELECT
			id,
			net_amount,
			paid_amount
		FROM laundry_orders
		WHERE id=?
	`,
		orderID,
	).
		Scan(&row).
		Error

	if err != nil {
		return nil, err
	}

	return &row, nil
}

func CreatePayment(
	db *gorm.DB,
	req CreatePaymentRequest,
	userID int64,
) error {

	return db.Transaction(func(tx *gorm.DB) error {

		orderInfo, err := GetOrderPaymentInfo(
			tx,
			req.OrderID,
		)

		if err != nil {
			return err
		}

		if orderInfo.PaidAmount > 0 {

			return errors.New(
				"order already paid",
			)
		}

		payment := Payment{
			OrderID: req.OrderID,

			PaymentMethodID: req.PaymentMethodID,

			Amount: orderInfo.NetAmount,

			PaidAt: time.Now(),

			PaidBy: userID,

			ReferenceNo: req.ReferenceNo,

			Notes: req.Notes,
		}

		if err := tx.Create(&payment).Error; err != nil {
			return err
		}

		if err := tx.Exec(`
			UPDATE laundry_orders
			SET paid_amount = ?
			WHERE id = ?
		`,
			orderInfo.NetAmount,
			req.OrderID,
		).Error; err != nil {

			return err
		}

		return nil
	})
}

type PaymentHistoryRow struct {
	ID int64 `json:"id"`

	Amount float64 `json:"amount"`

	Method string `json:"method"`

	PaidAt string `json:"paid_at"`

	ReferenceNo string `json:"reference_no"`
}

func GetPaymentHistory(
	db *gorm.DB,
	orderID int64,
) ([]PaymentHistoryRow, error) {

	var rows []PaymentHistoryRow

	err := db.Raw(`
		SELECT

			p.id,

			p.amount,

			pm.name method,

			TO_CHAR(
				p.paid_at,
				'YYYY-MM-DD HH24:MI'
			) paid_at,

			p.reference_no

		FROM payments p

		INNER JOIN payment_methods pm
			ON pm.id=p.payment_method_id

		WHERE p.order_id=?

		ORDER BY p.id DESC
	`,
		orderID,
	).
		Scan(&rows).
		Error

	return rows, err
}
