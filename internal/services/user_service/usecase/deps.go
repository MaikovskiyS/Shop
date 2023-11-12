package usecase

import (
	"context"
	"myproject/internal/domain"
)

type Storage interface {
	GetById(ctx context.Context, id uint64) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Save(ctx context.Context, u domain.User) (uint64, error)
	GetAll(ctx context.Context) ([]domain.User, error)
}
