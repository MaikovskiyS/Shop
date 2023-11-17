package storage

import (
	"context"
	"errors"
	"fmt"
	"log"
	"myproject/internal/apperrors"
	"myproject/internal/domain"
	"myproject/internal/services/order_service/model"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const location = "Order_Service-Storage-"

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
	p := &model.StoreOrder{}

	sql := "SELECT id, customer_name,total_price,status,created_at FROM orders WHERE id=$1"
	err := s.conn.QueryRow(ctx, sql, id).Scan(&p.ID, &p.CustomerName, &p.TotalPrice, &p.Status, &p.CreatedAt)
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

	defer rows.Close()

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

// Get all orders
func (s *storage) GetAll(ctx context.Context) ([]*model.StoreOrder, error) {
	orders := make([]*model.StoreOrder, 0)

	// tx, err := s.conn.Begin(ctx)
	// if err != nil {
	// 	ErrInternal.AddLocation("GetAll-s.conn.Begin")
	// 	ErrInternal.SetErr(err)
	// 	return nil, ErrInternal
	// }
	// defer tx.Rollback(ctx)
	sql := "SELECT * FROM orders"
	rows, err := s.conn.Query(ctx, sql)
	if err != nil {
		ErrInternal.AddLocation("GetAll-tx.Query")
		ErrInternal.SetErr(err)
		return nil, ErrInternal
	}
	defer rows.Close()
	ticker := time.Now()
	for rows.Next() {
		o := &model.StoreOrder{Products: []uint64{1, 2, 3}}
		err = rows.Scan(&o.ID, &o.UserID, &o.CustomerName, &o.TotalPrice, &o.Status, &o.CreatedAt)
		if err != nil {
			ErrInternal.AddLocation(fmt.Sprintf("GetAll-rows.Scan-Value: %v", rows.CommandTag().RowsAffected()))
			ErrInternal.Log()
			continue
		}
		log.Println(time.Since(ticker))
		//ptx, err := tx.Begin(ctx)
		//defer ptx.Rollback(ctx)
		if err != nil {
			ErrInternal.AddLocation("GetAll-tx.BeginProduct")
			ErrInternal.SetErr(err)
			ErrInternal.Log()
		}
		// pSql := "SELECT product_id FROM order_products WHERE order_id=$1"
		// rows, err := s.conn.Query(ctx, pSql, o.ID)
		// if err != nil {
		// 	ErrInternal.AddLocation("GetAll-tx.QueryProduct")
		// 	ErrInternal.SetErr(err)
		// 	return nil, ErrInternal
		// }

		// defer rows.Close()

		// for rows.Next() {
		// 	var pId uint64
		// 	err = rows.Scan(&pId)
		// 	if err != nil {
		// 		ErrInternal.AddLocation("GetAll-rows.ScanProduct")
		// 		ErrInternal.Log()
		// 		continue
		// 	}
		// 	o.Products = append(o.Products, pId)
		// }
		//	ptx.Commit(ctx)
		orders = append(orders, o)
	}
	//err = tx.Commit(ctx)
	// if err != nil {
	// 	ErrInternal.AddLocation("GetAll-tx.Commit")
	// 	ErrInternal.SetErr(err)
	// 	return nil, ErrInternal
	// }
	return orders, nil
}
