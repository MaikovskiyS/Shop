package order_service

import (
	"context"
	"errors"
	"log"
	"myproject/internal/apperrors"
	"myproject/internal/config"
	"myproject/internal/server/router"
	"myproject/internal/services/order_service/adapter/cache"
	"myproject/internal/services/order_service/adapter/repository"
	"myproject/internal/services/order_service/adapter/storage"
	"myproject/internal/services/order_service/transport/msgbroker/consumer"
	"myproject/internal/services/order_service/transport/rest"
	"myproject/internal/services/order_service/usecase"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func New(ctx context.Context, r *router.Router, pcl *pgxpool.Pool, pr usecase.Product, us usecase.User, rCl *redis.Client, mongoDb *mongo.Database, cfg *config.Config) usecase.Order {
	kafkaBrokers := []string{cfg.Kafka.HostPort()}

	mongoDb.CreateCollection(context.Background(), "orders", &options.CreateCollectionOptions{})

	mongoRepo := repository.New(mongoDb.Collection("orders"))

	msgReader := kafka.NewReader(kafka.ReaderConfig{Brokers: kafkaBrokers, Topic: cfg.Kafka.Topic, GroupID: cfg.Kafka.GroupId})

	cache := cache.New(rCl)

	store := storage.New(pcl)

	svc := usecase.New(store, pr, us, cache, mongoRepo)

	consumer := consumer.New(msgReader, svc)
	go func(context.Context) {
		err := consumer.SaveOrders(ctx)
		if err != nil {
			var er *apperrors.AppErr
			if errors.As(err, &er) {
				log.Println(er.Log())
			}
			log.Println("consumer.SaveOrdersErr:", err)
			<-ctx.Done()
			return
		}
	}(ctx)

	api := rest.New(svc)

	api.RegisterRoutes(r)
	return svc
}
