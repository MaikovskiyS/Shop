package storage

import (
	"context"
	"myproject/internal/apperrors"
	"myproject/internal/services/auth_service/model"

	"github.com/jackc/pgx/v5"
)

const location = "Auth_Service-Storage-"

var (
	ErrBadRequest = apperrors.New(apperrors.ErrBadRequest, location)
	ErrNotFound   = apperrors.New(apperrors.ErrNotFound, location)
	ErrInternal   = apperrors.New(apperrors.ErrInternal, location)
)

type storage struct {
	conn *pgx.Conn
}

func New(c *pgx.Conn) *storage {
	return &storage{
		conn: c,
	}
}
func (s *storage) Save(ctx context.Context, u model.User) (uint8, error) {
	sql := "INSERT INTO auth_users(email,password) values ($1,$2) RETURNING ID"
	row := s.conn.QueryRow(ctx, sql, u.Email, u.Password)
	var id uint8
	err := row.Scan(&id)
	if err != nil {
		if err == pgx.ErrNoRows {
			ErrNotFound.AddLocation("Save-ErrNoRows")
			ErrNotFound.SetErr(err)
			return 0, ErrNotFound
		}
		ErrInternal.AddLocation("Save-row.Scan")
		ErrInternal.SetErr(err)
		return 0, ErrInternal
	}
	return id, nil
}
func (s *storage) Get(ctx context.Context, email string) (*model.User, error) {
	sql := "SELECT * FROM auth_users WHERE email=$1"
	row := s.conn.QueryRow(ctx, sql, email)
	var user model.User
	var id uint64
	err := row.Scan(&id, &user.Email, &user.Password)
	if err != nil {
		if err == pgx.ErrNoRows {
			ErrNotFound.AddLocation("Get-ErrNoRows")
			ErrNotFound.SetErr(err)
			return &model.User{}, ErrNotFound
		}
		ErrInternal.AddLocation("Get-row.Scan")
		ErrInternal.SetErr(err)
		return &model.User{}, ErrInternal
	}
	return &user, nil
}
