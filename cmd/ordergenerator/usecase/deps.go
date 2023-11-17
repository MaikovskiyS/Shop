package usecase

import (
	"myproject/internal/domain"
)

type Sender interface {
	Send(o *domain.Order) error
}
type OrderGenerator interface {
	GenerateOrder() (*domain.Order, error)
}
