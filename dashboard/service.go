package dashboard

import "gorm.io/gorm"

type Summary struct {
	TotalOrders      int64   `json:"total_orders"`
	TodayOrders      int64   `json:"today_orders"`
	ProcessingOrders int64   `json:"processing_orders"`
	CompletedOrders  int64   `json:"completed_orders"`
	Outstanding      float64 `json:"outstanding"`
}

func GetSummary(
	db *gorm.DB,
) (Summary, error) {

	var result Summary

	db.Raw(`
		SELECT
			(SELECT COUNT(*) FROM laundry_orders) total_orders,

			(
			  SELECT COUNT(*)
			  FROM laundry_orders
			  WHERE DATE(order_date)=CURRENT_DATE
			) today_orders,

			(
			  SELECT COUNT(*)
			  FROM laundry_orders
			  WHERE current_status_id NOT IN (
			      SELECT id
			      FROM laundry_statuses
			      WHERE code='COMPLETED'
			  )
			) processing_orders,

			(
			  SELECT COUNT(*)
			  FROM laundry_orders
			  WHERE current_status_id IN (
			      SELECT id
			      FROM laundry_statuses
			      WHERE code='COMPLETED'
			  )
			) completed_orders,

			(
			  SELECT COALESCE(
			        SUM(net_amount-paid_amount),
			        0
			  )
			  FROM laundry_orders
			) outstanding
	`).Scan(&result)

	return result, nil
}

type RecentOrder struct {
	OrderNo string `json:"order_no"`

	CustomerName string `json:"customer_name"`

	Status string `json:"status"`

	Amount float64 `json:"amount"`

	OrderDate string `json:"order_date"`
}

func GetRecentOrders(
	db *gorm.DB,
) ([]RecentOrder, error) {

	var rows []RecentOrder

	err := db.Raw(`
		SELECT

			o.order_no,

			c.name customer_name,

			s.name status,

			o.net_amount amount,

			TO_CHAR(
				o.order_date,
				'DD Mon YYYY HH24:MI'
			) order_date

		FROM laundry_orders o

		INNER JOIN customers c
			ON c.id=o.customer_id

		INNER JOIN laundry_statuses s
			ON s.id=o.current_status_id

		ORDER BY o.id DESC

		LIMIT 10
	`).Scan(&rows).Error

	return rows, err
}
