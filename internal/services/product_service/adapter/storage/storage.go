package storage

import (
	"context"
	"errors"
	"log"
	"myproject/internal/apperrors"
	"myproject/internal/domain"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const location = "Product_Service-Storage-"

var (
	ErrBadRequest = apperrors.New(apperrors.ErrBadRequest, location)
	ErrInternal   = apperrors.New(apperrors.ErrInternal, location)
	ErrNotFound   = apperrors.New(apperrors.ErrNotFound, location)
)

type storage struct {
	conn *pgxpool.Pool
}

func New(c *pgxpool.Pool) *storage {
	return &storage{
		conn: c,
	}
}

// Save product in database
func (s *storage) Save(ctx context.Context, p domain.Product) (uint64, error) {
	tx, err := s.conn.Begin(ctx)
	if err != nil {
		ErrInternal.AddLocation("Save-s.conn.Begin")
		ErrInternal.SetErr(err)
		return 0, ErrInternal
	}
	var catId int
	err = tx.QueryRow(ctx, "SELECT id FROM categories WHERE name=$1", p.Category).Scan(&catId)
	if err != nil {
		if err == pgx.ErrNoRows {
			tx.Rollback(ctx)
			ErrNotFound.AddLocation("Save-SelectErrNoRows")
			ErrNotFound.SetErr(errors.New("category not found"))
			return 0, ErrNotFound
		}
		tx.Rollback(ctx)
		ErrInternal.AddLocation("Save-tx.QueryRow")
		ErrInternal.SetErr(err)
		return 0, ErrInternal
	}
	p.CreatedAt = time.Now()
	var id uint64
	row := tx.QueryRow(ctx, "INSERT INTO products(category_id,sku,name,price, image, created_at) VALUES($1,$2,$3,$4,$5,$6) RETURNING ID", catId, p.Sku, p.Name, p.Price, p.Image, p.CreatedAt)
	err = row.Scan(&id)
	if err != nil {
		tx.Rollback(ctx)
		ErrInternal.AddLocation("Save-row.Scan")
		ErrInternal.SetErr(err)
		return 0, ErrInternal
	}
	tx.Commit(ctx)
	return id, nil
}

// Get product from database by id
func (s *storage) GetByID(ctx context.Context, id uint64) (*domain.Product, error) {
	sql := "SELECT categories.name,sku,products.name,price,image,created_at FROM products JOIN categories ON products.category_id=categories.id WHERE products.id=$1"
	row := s.conn.QueryRow(ctx, sql, &id)
	p := &domain.Product{}

	err := row.Scan(&p.Category, &p.Sku, &p.Name, &p.Price, &p.Image, &p.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			ErrNotFound.AddLocation("GetByID-ErrNoRows")
			ErrNotFound.SetErr(errors.New("product not found"))
			return nil, ErrNotFound
		}
		ErrInternal.AddLocation("GetByID-row.Scan")
		ErrInternal.SetErr(err)
		return &domain.Product{}, ErrInternal
	}
	p.ID = id
	return p, nil
}

// Get all products
func (s *storage) GetAll(ctx context.Context) ([]*domain.Product, error) {
	products := make([]*domain.Product, 0)

	sql := "SELECT categories.name,sku,products.name,price,image,created_at FROM products JOIN categories ON products.category_id=categories.id"
	rows, err := s.conn.Query(ctx, sql)
	if err != nil {
		ErrInternal.AddLocation("GetByID-s.conn.Query")
		ErrInternal.SetErr(err)
		return nil, ErrInternal
	}
	defer rows.Close()

	for rows.Next() {
		var p domain.Product
		err := rows.Scan(&p.Category, &p.Sku, &p.Name, &p.Price, &p.Image, &p.CreatedAt)
		if err != nil {
			log.Println(err)
			continue
		}
		products = append(products, &p)
	}
	return products, nil
}
