package repository

import (
	"context"
	"log"
	"myproject/internal/apperrors"
	"myproject/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrInternal = apperrors.New(apperrors.ErrInternal, "Order_Service-Repository-")
)

type store struct {
	coll *mongo.Collection
}

func New(coll *mongo.Collection) *store {
	return &store{
		coll: coll,
	}
}
func (s *store) Save(ctx context.Context, or *domain.Order) (uint64, error) {
	_, err := s.coll.InsertOne(ctx, or)
	if err != nil {
		ErrInternal.AddLocation("Save-IncertOne")
		ErrInternal.SetErr(err)
		return 0, ErrInternal
	}
	log.Println("order saved")
	return 0, nil
}
func (s *store) GetAll(ctx context.Context) ([]*domain.Order, error) {
	filter := bson.D{{}}
	cur, err := s.coll.Find(ctx, &filter)
	if err != nil {
		ErrInternal.AddLocation("GetAll-Find")
		ErrInternal.SetErr(err)
		return nil, ErrInternal
	}
	var orders []*domain.Order
	for cur.Next(ctx) {
		var t *domain.Order
		err := cur.Decode(&t)
		if err != nil {
			ErrInternal.AddLocation("GetAll-cur.Decode")
			ErrInternal.SetErr(err)
			return nil, ErrInternal
		}

		orders = append(orders, t)
	}

	if err := cur.Err(); err != nil {
		ErrInternal.AddLocation("GetAll-cur.Err")
		ErrInternal.SetErr(err)
		return nil, ErrInternal
	}
	return orders, nil
}
