package cache

import (
	"context"
	"encoding/json"
	"myproject/internal/apperrors"
	"myproject/internal/domain"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ErrInternal = apperrors.New(apperrors.ErrInternal, "Order-Service-")
	ErrNotFound = apperrors.New(apperrors.ErrNotFound, "Order-Service-")
)

type cache struct {
	ttl time.Duration
	r   *redis.Client
}

func New(rCl *redis.Client) *cache {
	return &cache{
		ttl: time.Hour * 2,
		r:   rCl,
	}
}
func (c *cache) Set(ctx context.Context, key uint64, o *domain.Order) error {
	sId := strconv.Itoa(int(key))
	orderBytes, err := json.Marshal(&o)
	if err != nil {
		ErrInternal.AddLocation("Set-MarshalOrder")
		ErrInternal.SetErr(err)
		return ErrInternal
	}

	st := c.r.Set(ctx, sId, orderBytes, c.ttl)
	return st.Err()
}
func (c *cache) Get(ctx context.Context, key uint64) (*domain.Order, error) {
	sId := strconv.Itoa(int(key))

	res := c.r.Get(ctx, sId)
	respBytes, err := res.Bytes()
	if err != nil {
		return &domain.Order{}, err
	}
	var order *domain.Order
	err = json.Unmarshal(respBytes, &order)
	if err != nil {
		ErrInternal.AddLocation("Get-Unmarshal")
		ErrInternal.SetErr(err)
		return &domain.Order{}, ErrInternal
	}

	return order, nil
}
