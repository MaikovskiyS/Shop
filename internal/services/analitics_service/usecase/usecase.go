package usecase

import (
	"context"
	"myproject/internal/domain"
)

type usecase struct {
	repo Storage
}

func New(s Storage) *usecase {
	return &usecase{
		repo: s,
	}
}
func (u *usecase) SaveOrder(ctx context.Context, or *domain.Order) (uint64, error) {
	return u.repo.Save(ctx, or)

}
func (u *usecase) GetAll(ctx context.Context) ([]*domain.Order, error) {
	return u.repo.GetAll(ctx)
}
