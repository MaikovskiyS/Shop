package domain

import "time"

type Order struct {
	ID           uint64     `json:"id"`
	UserID       uint64     `json:"user_id"`
	CustomerName string     `json:"customer_name"`
	TotalPrice   float64    `json:"total_price"`
	CreatedAt    time.Time  `json:"created_at"`
	Status       string     `json:"status"`
	Products     []*Product `json:"products"`
}
