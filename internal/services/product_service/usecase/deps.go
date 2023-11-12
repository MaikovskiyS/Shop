package usecase

import (
	"context"
	"myproject/internal/domain"
)

type Storage interface {
	Save(ctx context.Context, p domain.Product) (uint64, error)
	GetByID(ctx context.Context, id uint64) (*domain.Product, error)
	GetAll(ctx context.Context) ([]*domain.Product, error)
}
