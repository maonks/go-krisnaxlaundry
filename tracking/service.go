package tracking

import "gorm.io/gorm"

type TrackingTimeline struct {
	StatusName string `json:"status_name"`

	Remarks string `json:"remarks"`

	ChangedAt string `json:"changed_at"`
}

type TrackingResponse struct {
	OrderNo string `json:"order_no"`

	CustomerName string `json:"customer_name"`

	CustomerPhone string `json:"customer_phone"`

	Status string `json:"status"`

	OrderDate string `json:"order_date"`

	EstimatedFinish string `json:"estimated_finish"`

	Timeline []TrackingTimeline `json:"timeline"`
}

func GetTracking(
	db *gorm.DB,
	orderNo string,
) (*TrackingResponse, error) {

	var result TrackingResponse

	err := db.Raw(`
		SELECT

			o.order_no,

			c.name customer_name,

			c.phone customer_phone,

			s.name status,

			TO_CHAR(
				o.order_date,
				'DD Mon YYYY HH24:MI'
			) order_date,

			COALESCE(
				TO_CHAR(
					o.estimated_finish_at,
					'DD Mon YYYY HH24:MI'
				),
				'-'
			) estimated_finish

		FROM laundry_orders o

		INNER JOIN customers c
			ON c.id = o.customer_id

		INNER JOIN laundry_statuses s
			ON s.id = o.current_status_id

		WHERE o.order_no = ?
	`,
		orderNo,
	).Scan(&result).Error

	if err != nil {
		return nil, err
	}

	var timeline []TrackingTimeline

	err = db.Raw(`
		SELECT

			s.name status_name,

			COALESCE(
				h.remarks,
				''
			) remarks,

			TO_CHAR(
				h.changed_at,
				'DD Mon YYYY HH24:MI'
			) changed_at

		FROM laundry_status_histories h

		INNER JOIN laundry_statuses s
			ON s.id = h.status_id

		WHERE h.order_id = (
			SELECT id
			FROM laundry_orders
			WHERE order_no = ?
		)

		ORDER BY h.changed_at ASC
	`,
		orderNo,
	).Scan(&timeline).Error

	if err != nil {
		return nil, err
	}

	result.Timeline = timeline

	return &result, nil
}
