package order_service

import (
	"context"
	"errors"
	"myproject/internal/apperrors"
	"myproject/internal/config"
	"myproject/internal/server/router"
	"myproject/internal/services/order_service/adapter/cache"
	"myproject/internal/services/order_service/adapter/storage"
	"myproject/internal/services/order_service/transport/msgbroker/consumer"
	"myproject/internal/services/order_service/transport/rest"
	"myproject/internal/services/order_service/usecase"

	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
)

func New(ctx context.Context, r *router.Router, pcl *pgx.Conn, pr usecase.Product, us usecase.User, rCl *redis.Client, cfg *config.Config) usecase.Order {
	br := []string{cfg.Kafka.HostPort()}
	msgReader := kafka.NewReader(kafka.ReaderConfig{Brokers: br, Topic: cfg.Kafka.Topic, GroupID: cfg.Kafka.GroupId})
	cache := cache.New(rCl)
	store := storage.New(pcl)
	svc := usecase.New(store, pr, us, cache)
	consumer := consumer.New(msgReader, svc)
	go func(context.Context) {
		err := consumer.SaveOrders(ctx)
		if err != nil {
			var er *apperrors.AppErr
			if errors.As(err, &er) {
				er.Log()
			}
		}
	}(ctx)

	api := rest.New(svc)

	api.RegisterRoutes(r)
	return svc
}
