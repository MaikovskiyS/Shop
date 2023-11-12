package model

import "time"

type Order struct {
	UserId      uint64
	ProductsIds []uint64
}

type StoreOrder struct {
	ID           uint64
	UserID       uint64
	CustomerName string
	TotalPrice   float64
	CreatedAt    time.Time
	Status       string
	Products     []uint64
}
