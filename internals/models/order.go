package models

import "time"

// this referencies a product and quantity
type OrderItem struct {
	ProductID uint  `json:"product_id"`
	Qty       int   `json:"qty"`
	UnitPrice int64 `json:"unit_price_cents"`
}

// represents an simple order

type Order struct {
	ID         uint        `json:"id"`
	Customer   string      `json:"customer" binding:"required"`
	Items      []OrderItem `json:"items" binding:"required"`
	TotalCents int64       `json:"total_cents"`
	CreatedAt  time.Time   `json:"created_at"`
}
