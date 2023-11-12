package usecase

import (
	"context"
	"myproject/internal/domain"
)

type User interface {
	GetAll(ctx context.Context) ([]domain.User, error)
	Save(ctx context.Context, u domain.User) (uint64, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetById(ctx context.Context, id uint64) (*domain.User, error)
}
type usecase struct {
	user Storage
}

func New(s Storage) *usecase {
	return &usecase{
		user: s,
	}
}

func (u *usecase) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return u.user.GetByEmail(ctx, email)
}
func (u *usecase) GetById(ctx context.Context, id uint64) (*domain.User, error) {
	return u.user.GetById(ctx, id)
}
func (u *usecase) Save(ctx context.Context, us domain.User) (uint64, error) {
	return u.user.Save(ctx, us)
}
func (u *usecase) GetAll(ctx context.Context) ([]domain.User, error) {
	return u.user.GetAll(ctx)
}
