package analitics_service

import (
	"context"
	"log"
	"myproject/internal/config"
	"myproject/internal/server/router"
	"myproject/internal/services/analitics_service/controllers/api"
	"myproject/internal/services/analitics_service/controllers/msgbroker"
	"myproject/internal/services/analitics_service/repository"
	"myproject/internal/services/analitics_service/usecase"

	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/mongo"
)

func New(ctx context.Context, r *router.Router, msgBr *kafka.Conn, mDb *mongo.Database, cfg *config.Config) {

	kafkaBrokers := []string{cfg.Kafka.HostPort()}
	msgReader := kafka.NewReader(kafka.ReaderConfig{Brokers: kafkaBrokers, Topic: "analitics", GroupID: "anId"})

	mColl := mDb.Collection("orders")

	repo := repository.New(mColl)
	svc := usecase.New(repo)
	consumer := msgbroker.New(msgReader, svc)
	api := api.New(svc)

	api.RegisterRoutes(r)

	go func() {
		err := consumer.SaveOrders(context.Background())
		if err != nil {
			log.Println("analitics consumer err:", err)
			return
		}
	}()

}
