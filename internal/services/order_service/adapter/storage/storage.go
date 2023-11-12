package storage

import (
	"context"
	"errors"
	"myproject/internal/apperrors"
	"myproject/internal/domain"
	"myproject/internal/services/order_service/model"

	"github.com/jackc/pgx/v5"
)

const location = "Order_Service-Storage"

var (
	ErrBadRequest = apperrors.New(apperrors.ErrBadRequest, location)
	ErrInternal   = apperrors.New(apperrors.ErrInternal, location)
	ErrNotFound   = apperrors.New(apperrors.ErrNotFound, location)
)

type storage struct {
	conn *pgx.Conn
}

func New(c *pgx.Conn) *storage {
	return &storage{
		conn: c,
	}
}

// Save order into order_table and products ids relations into order_products. Returning id.
func (s *storage) Save(ctx context.Context, p *domain.Order) (uint64, error) {
	tx, err := s.conn.Begin(ctx)
	defer func() {
		tx.Rollback(ctx)
	}()
	if err != nil {
		ErrInternal.AddLocation("Save-s.conn.Begin")
		ErrInternal.SetErr(err)
		return 0, ErrInternal
	}
	sql := "INSERT INTO orders(user_id,customer_name,total_price,created_at,status) VALUES($1,$2,$3,$4,$5) RETURNING ID"
	row := tx.QueryRow(ctx, sql, p.UserID, p.CustomerName, p.TotalPrice, p.CreatedAt, p.Status)
	var id uint64
	err = row.Scan(&id)
	if err != nil {
		if err == pgx.ErrNoRows {
			ErrNotFound.AddLocation("Save-ErrNoRows")
			ErrNotFound.SetErr(errors.New("order not found"))
			return 0, ErrNotFound
		}
		ErrInternal.AddLocation("Save-row.Scan")
		ErrInternal.SetErr(err)
		return 0, ErrInternal
	}
	for _, product := range p.Products {
		sql := "INSERT INTO order_products(order_id,product_id) VALUES($1,$2)"
		_, err := tx.Exec(ctx, sql, id, product.ID)
		if err != nil {
			ErrInternal.AddLocation("Save-tx.Exec")
			ErrInternal.SetErr(err)
			return 0, ErrInternal
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}
	return id, nil
}

// Getting order by id from 'orders' and order products ids from 'order_products'. Returning Order
func (s *storage) GetById(ctx context.Context, id uint64) (*model.StoreOrder, error) {
	sql := "SELECT id, customer_name,total_price,status,created_at FROM orders WHERE id=$1"
	row := s.conn.QueryRow(ctx, sql, id)
	p := &model.StoreOrder{}
	err := row.Scan(&p.ID, &p.CustomerName, &p.TotalPrice, &p.Status, &p.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			ErrNotFound.AddLocation("GetByID-ErrNoRows")
			ErrNotFound.SetErr(errors.New("order not found"))
			return nil, ErrNotFound
		}
		ErrInternal.AddLocation("GetByID-row.Scan")
		ErrInternal.SetErr(err)
		return &model.StoreOrder{}, ErrInternal
	}
	pIds := make([]uint64, 0)
	rows, err := s.conn.Query(ctx, "SELECT product_id FROM order_products WHERE order_id=$1", id)
	if err != nil {
		ErrInternal.AddLocation("GetByID-s.conn.Query")
		ErrInternal.SetErr(err)
		return &model.StoreOrder{}, ErrInternal
	}
	for rows.Next() {
		var id uint64
		err = rows.Scan(&id)
		if err != nil {
			ErrInternal.AddLocation("GetByID-row.Scan")
			ErrInternal.SetErr(err)
			return &model.StoreOrder{}, ErrInternal
		}
		pIds = append(pIds, id)
	}
	p.Products = pIds
	return p, nil
}
