package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"myproject/internal/apperrors"
	"myproject/internal/domain"
	"myproject/internal/services/order_service/model"
	"time"

	"github.com/segmentio/kafka-go"
)

var (
	ErrNotFound = apperrors.New(apperrors.ErrNotFound, "Order_Service-Usecase-")
	ErrInternal = apperrors.New(apperrors.ErrInternal, "Order_Service-Usecase-")
)

type Order interface {
	Save(ctx context.Context, o model.Order) (uint64, error)
	GetById(ctx context.Context, id uint64) (*domain.Order, error)
	GetAll(ctx context.Context) ([]*domain.Order, error)
}
type usecase struct {
	producer *kafka.Writer
	user     User
	cache    Cache
	product  Product
	store    Storage
}

func New(s Storage, pr Product, u User, c Cache, ow *kafka.Writer) *usecase {
	return &usecase{
		cache:    c,
		product:  pr,
		store:    s,
		user:     u,
		producer: ow,
	}
}

// Get order by id from cache if not exist, collect order data from order repository and product_service; set order to cache. Returning order
func (u *usecase) GetById(ctx context.Context, id uint64) (*domain.Order, error) {
	order, err := u.cache.Get(ctx, id)
	if err != nil {
		var er *apperrors.AppErr
		if errors.As(err, &er) {
			if er.Type() == apperrors.ErrNotFound {
				mOrder, err := u.store.GetById(ctx, id)
				if err != nil {
					return &domain.Order{}, err
				}
				products := make([]*domain.Product, len(mOrder.Products))
				for i, pId := range mOrder.Products {
					p, err := u.product.GetById(ctx, pId)
					if err != nil {
						return &domain.Order{}, err
					}
					products[i] = p
				}
				o := &domain.Order{
					ID:           mOrder.ID,
					Products:     products,
					CustomerName: mOrder.CustomerName,
					Status:       mOrder.Status,
					CreatedAt:    mOrder.CreatedAt,
					TotalPrice:   mOrder.TotalPrice,
				}
				err = u.cache.Set(ctx, o.ID, o)
				if err != nil {
					return &domain.Order{}, err
				}
				return o, nil
			}
			return &domain.Order{}, err
		}

	}

	return order, nil
}

// Save checking user
func (u *usecase) Save(ctx context.Context, or model.Order) (uint64, error) {
	user, err := u.user.GetById(ctx, or.UserId)
	if err != nil {
		var er *apperrors.AppErr
		if errors.As(err, &er) {
			if er.Type() == apperrors.ErrNotFound {
				ErrNotFound.AddLocation("Save-CheckUser")
				ErrNotFound.SetErr(errors.New("user not found"))
				return 0, ErrNotFound
			}
			return 0, er
		}
		return 0, err
	}
	products := make([]*domain.Product, len(or.ProductsIds))
	totalPrice := 0.0
	for i, pId := range or.ProductsIds {
		product, err := u.product.GetById(ctx, pId)
		if err != nil {
			return 0, err
		}
		totalPrice += product.Price
		products[i] = product
	}

	order := &domain.Order{
		UserID:       user.Id,
		CustomerName: user.Name,
		TotalPrice:   totalPrice,
		Products:     products,
		CreatedAt:    time.Now(),
		//add storage with statuses
		Status: "confirm",
	}
	orderId, err := u.store.Save(ctx, order)
	if err != nil {
		return 0, err
	}

	order.ID = orderId
	oBytes, err := json.Marshal(&order)
	if err != nil {
		return 0, err
	}
	msg := kafka.Message{
		Value: oBytes,
	}
	err = u.producer.WriteMessages(ctx, msg)
	if err != nil {
		log.Println("order producer err", err)
		return 0, err
	}
	err = u.cache.Set(ctx, orderId, order)
	if err != nil {
		return 0, nil
	}
	return orderId, nil
}
func (u *usecase) GetAll(ctx context.Context) ([]*domain.Order, error) {
	mOrders, err := u.store.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	if mOrders == nil {
		ErrInternal.AddLocation("GetAll-CheckOrders")
		ErrInternal.SetErr(errors.New("nil orders from storage"))
		return nil, ErrInternal
	}

	orders := make([]*domain.Order, len(mOrders))

	for _, mOrder := range mOrders {
		order := &domain.Order{
			ID:           mOrder.ID,
			UserID:       mOrder.UserID,
			CustomerName: mOrder.CustomerName,
			TotalPrice:   mOrder.TotalPrice,
			CreatedAt:    mOrder.CreatedAt,
			Status:       mOrder.Status,
			Products:     make([]*domain.Product, len(mOrder.Products)),
		}

		tick := time.Now()
		for i, pId := range mOrder.Products {

			product, err := u.product.GetById(ctx, pId)
			if err != nil {
				return nil, err
			}

			order.Products[i] = product
		}
		log.Println(time.Since(tick))

		orders = append(orders, order)
	}
	return orders, nil
}
