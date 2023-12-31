package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"myproject/cmd/ordergenerator"
	"myproject/internal/config"
	"myproject/internal/metrics"
	"myproject/internal/server"
	"myproject/internal/server/router"
	"myproject/internal/services/analitics_service"
	"myproject/internal/services/auth_service"
	"myproject/internal/services/order_service"
	"myproject/internal/services/product_service"
	"myproject/internal/services/user_service"
	"myproject/pkg/autorization"
	"net/http"
	"os"
	"os/signal"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// init dependences; starting server
func Run(cfg *config.Config) error {

	//init redis connection
	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.Redis.HostPort(),
	})

	//init postgres connection
	pdb, err := pgx.Connect(context.Background(), cfg.Psql.ConnString())
	if err != nil {
		return fmt.Errorf("pgx conn Err: %w", err)
	}
	cc, err := pgxpool.ParseConfig(cfg.Psql.ConnString())
	if err != nil {
		return err
	}
	pool, err := pgxpool.NewWithConfig(context.Background(), cc)
	if err != nil {
		return err
	}

	// init mongo client
	clientOptions := options.Client()
	clientOptions.SetMaxPoolSize(80)
	clientOptions.ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return err
	}
	mongoDb := client.Database("shop")

	//init kafka connection
	kbroker, err := kafka.DialLeader(context.Background(), cfg.Kafka.Protocol, cfg.Kafka.HostPort(), cfg.Kafka.Topic, 0)
	if err != nil {
		return fmt.Errorf("kafka conn err %w", err)
	}

	//init token service
	tokenSvc, err := autorization.New()
	if err != nil {
		return err
	}

	//init router
	router := router.New(tokenSvc)

	//register metrics
	go metrics.Run()

	//run order generator
	genCtx, genExit := context.WithCancel(context.Background())

	ogen := ordergenerator.New(genCtx, kbroker, router)

	//init services
	appCtx, appExit := context.WithCancel(context.Background())

	us := user_service.New(router, pdb, rdb)
	auth_service.New(router, pdb, tokenSvc, us)
	pr := product_service.New(router, pool)
	order_service.New(appCtx, router, pool, pr, us, rdb, cfg)
	analitics_service.New(appCtx, router, kbroker, mongoDb, cfg)

	//init http server
	srv := server.New(cfg)

	srv.SetHandler(router)

	//shutdown server and connections
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint
		genExit()
		appExit()
		if err = pdb.Close(context.Background()); err != nil {
			log.Printf("Postgres conn Close Err: %s", err)
		}
		if err = rdb.Close(); err != nil {
			log.Printf("Redis conn Close Err: %s", err)
		}
		if err = kbroker.Close(); err != nil {
			log.Printf("Kafka conn Close Err: %s", err)
		}
		if err = mongoDb.Client().Disconnect(context.Background()); err != nil {
			log.Printf("Mongo conn Close Err: %s", err)
		}
		log.Println("connections closed")

		if err := srv.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}
	}()
	go func(context.Context) {
		err := ogen.Generate(genCtx)
		if err != nil {
			log.Println("genStop Err: ", err)
			return
		}
	}(genCtx)
	log.Println("Starting server")
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}
	return errors.New("app closed")
}
