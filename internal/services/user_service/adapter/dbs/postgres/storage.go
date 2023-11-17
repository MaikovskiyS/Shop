package psql

import (
	"context"
	"errors"
	"myproject/internal/apperrors"
	"myproject/internal/domain"

	"github.com/jackc/pgx/v5"
)

const location = "User_Service-Storage-"

var (
	ErrBadRequest = apperrors.New(apperrors.ErrBadRequest, location)
	ErrInternal   = apperrors.New(apperrors.ErrInternal, location)
	ErrNotFound   = apperrors.New(apperrors.ErrNotFound, location)
)

type store struct {
	conn *pgx.Conn
}

func New(c *pgx.Conn) *store {
	return &store{
		conn: c,
	}
}

// Save Users in database
func (s *store) Save(ctx context.Context, u domain.User) (uint64, error) {
	tx, err := s.conn.Begin(ctx)
	defer func() {
		tx.Rollback(ctx)
	}()
	if err != nil {
		ErrInternal.AddLocation("Save-s.conn.Begin")
		ErrInternal.SetErr(err)
		return 0, ErrInternal
	}
	d, err := tx.Prepare(ctx, "insert user", "INSERT INTO users(name,age) values ($1,$2) RETURNING ID")
	if err != nil {
		ErrInternal.AddLocation("Save-tx.Prepare")
		ErrInternal.SetErr(err)
		return 0, ErrInternal
	}
	var id uint64
	err = tx.QueryRow(ctx, d.SQL, u.Name, u.Age).Scan(&id)
	if err != nil {
		ErrInternal.AddLocation("Save-tx.QueryRow")
		ErrInternal.SetErr(err)
		return 0, ErrInternal
	}
	tag, err := tx.Exec(ctx, "INSERT INTO auth_users(user_id,email,password) VALUES ($1,$2,$3)", id, u.Email, u.Password)
	if err != nil {
		ErrInternal.AddLocation("Save-tx.Exec")
		ErrInternal.SetErr(err)
		return 0, ErrInternal
	}
	if tag.RowsAffected() == 0 {
		ErrInternal.AddLocation("Save-tag.RowsAffected")
		ErrInternal.SetErr(err)
		return 0, ErrInternal
	}
	if err = tx.Commit(ctx); err != nil {
		return 0, err
	}
	return id, nil
}

// Get all Users
func (s *store) GetAll(ctx context.Context) ([]domain.User, error) {
	sql := "SELECT user_id, name, email,age FROM users JOIN auth_users ON users.id=auth_users.id"
	rows, err := s.conn.Query(ctx, sql)
	if err != nil {
		if err == pgx.ErrNoRows {
			ErrNotFound.AddLocation("GetAll-ErrNoRows")
			ErrNotFound.SetErr(errors.New("user not found"))
			return nil, ErrNotFound
		}
		ErrInternal.AddLocation("Save-s.conn.Query")
		ErrInternal.SetErr(err)
		return nil, ErrInternal
	}
	users := make([]domain.User, 0)
	for rows.Next() {
		var user domain.User
		err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.Age)
		if err != nil {
			continue
		}
		users = append(users, user)

	}
	return users, nil
}

// Get User by email
func (s *store) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	user := &domain.User{}

	sql := "SELECT user_id, name, email,password,age FROM auth_users JOIN users ON users.id=auth_users.user_id WHERE email=$1"
	d, err := s.conn.Prepare(ctx, "userByEmail", sql)
	if err != nil {
		ErrInternal.AddLocation("GetByEmail-s.conn.Prepare")
		ErrInternal.SetErr(err)
		return &domain.User{}, ErrInternal
	}

	var age int
	err = s.conn.QueryRow(ctx, d.SQL, email).Scan(&user.Id, &user.Name, &user.Email, &user.Password, &age)
	if err != nil {
		if err == pgx.ErrNoRows {
			ErrNotFound.AddLocation("GetByEmail-ErrNoRows")
			ErrNotFound.SetErr(errors.New("user not found"))
			return &domain.User{}, ErrNotFound
		}
		ErrInternal.AddLocation("GetByEmail-s.conn.QueryRow")
		ErrInternal.SetErr(err)
		return &domain.User{}, ErrInternal
	}
	user.Age = uint8(age)
	return user, nil
}

// Get User by id
func (s *store) GetById(ctx context.Context, id uint64) (*domain.User, error) {
	user := &domain.User{}

	sql := "SELECT user_id, name, email,password,age FROM auth_users JOIN users ON users.id=auth_users.user_id WHERE users.id=$1"
	d, err := s.conn.Prepare(ctx, "userById", sql)
	if err != nil {
		ErrInternal.AddLocation("GetByID-s.conn.Prepare")
		ErrInternal.SetErr(err)
		return &domain.User{}, ErrInternal
	}

	var age int
	err = s.conn.QueryRow(ctx, d.SQL, id).Scan(&user.Id, &user.Name, &user.Email, &user.Password, &age)
	if err != nil {
		if err == pgx.ErrNoRows {
			ErrNotFound.AddLocation("GetByID-ErrNoRows")
			ErrNotFound.SetErr(errors.New("user not found"))
			return &domain.User{}, ErrNotFound
		}
		ErrInternal.AddLocation("GetByID-s.conn.QueryRow")
		ErrInternal.SetErr(err)
		return &domain.User{}, ErrInternal
	}
	user.Age = uint8(age)
	return user, nil
}
