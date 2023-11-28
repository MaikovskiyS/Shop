package usecase

import (
	"context"
	"myproject/internal/domain"
	"myproject/internal/services/order_service/model"
)

type User interface {
	GetAll(ctx context.Context) ([]domain.User, error)
	Save(ctx context.Context, u domain.User) (uint64, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetById(ctx context.Context, id uint64) (*domain.User, error)
}
type Product interface {
	Save(ctx context.Context, p domain.Product) (uint64, error)
	GetById(ctx context.Context, id uint64) (*domain.Product, error)
	GetAll(ctx context.Context) ([]*domain.Product, error)
}
type Storage interface {
	Save(ctx context.Context, p *domain.Order) (uint64, error)
	GetById(ctx context.Context, id uint64) (*model.StoreOrder, error)
	GetAll(ctx context.Context) ([]*model.StoreOrder, error)
}
type Cache interface {
	Get(ctx context.Context, key uint64) (*domain.Order, error)
	Set(ctx context.Context, key uint64, o *domain.Order) error
}
