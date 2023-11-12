package usecase

import (
	"context"
	"errors"
	"myproject/internal/apperrors"
	"myproject/internal/domain"
	"myproject/internal/services/order_service/model"
	"time"
)

var ErrNotFound = apperrors.New(apperrors.ErrNotFound, "Order_Service-Usecase-")

type Order interface {
	Save(ctx context.Context, o model.Order) (uint64, error)
	GetById(ctx context.Context, id uint64) (*domain.Order, error)
}
type usecase struct {
	user    User
	cache   Cache
	product Product
	store   Storage
}

func New(s Storage, pr Product, u User, c Cache) *usecase {
	return &usecase{
		cache:   c,
		product: pr,
		store:   s,
		user:    u,
	}
}

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
			return 0, err
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
	err = u.cache.Set(ctx, orderId, order)
	if err != nil {
		return 0, nil
	}
	return orderId, nil
}
