package domain

import "time"

type Product struct {
	ID        uint64
	Sku       string
	Category  string
	Name      string
	Price     float64
	Image     string
	CreatedAt time.Time
}
