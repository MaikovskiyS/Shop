package usecase

import (
	"context"
	"myproject/internal/domain"
)

type TokenService interface {
	CreateToken(email string) (string, error)
}
type UserService interface {
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Save(ctx context.Context, u domain.User) (uint64, error)
}
