package usecase

import (
	"context"
	"myproject/internal/domain"
)

type Storage interface {
	Save(ctx context.Context, or *domain.Order) (uint64, error)
	GetAll(ctx context.Context) ([]*domain.Order, error)
}
