package usecase

import (
	"context"
	"myproject/internal/domain"
)

type Product interface {
	Save(ctx context.Context, p domain.Product) (uint64, error)
	GetById(ctx context.Context, id uint64) (*domain.Product, error)
	GetAll(ctx context.Context) ([]*domain.Product, error)
}
type usecase struct {
	store Storage
}

func New(s Storage) *usecase {
	return &usecase{
		store: s,
	}
}

func (u *usecase) GetById(ctx context.Context, id uint64) (*domain.Product, error) {

	return u.store.GetByID(ctx, id)
}
func (u *usecase) Save(ctx context.Context, p domain.Product) (uint64, error) {
	return u.store.Save(ctx, p)
}
func (u *usecase) GetAll(ctx context.Context) ([]*domain.Product, error) {
	return u.store.GetAll(ctx)
}
