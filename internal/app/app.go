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
	"myproject/internal/services/auth_service"
	"myproject/internal/services/order_service"
	"myproject/internal/services/product_service"
	"myproject/internal/services/user_service"
	"myproject/pkg/autorization"
	"net/http"
	"os"
	"os/signal"

	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
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
	go func() {
		//TODO: gorutine leaking
		log.Println("register metrics")
		metrics.Register(router)
	}()
	//run order generator
	genCtx, genExit := context.WithCancel(context.Background())

	ogen := ordergenerator.New(genCtx, kbroker, router)
	go func(context.Context) {
		err := ogen.Generate(genCtx)
		if err != nil {
			log.Println("genStop Err: ", err)
			return
		}
	}(genCtx)

	//init services
	appCtx, appExit := context.WithCancel(context.Background())

	us := user_service.New(router, pdb, rdb)
	auth_service.New(router, pdb, tokenSvc, us)
	pr := product_service.New(router, pdb)
	order_service.New(appCtx, router, pdb, pr, us, rdb, cfg)

	//init server
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
		log.Println("connections closed")

		if err := srv.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}
	}()
	log.Println("Starting server")
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}
	return errors.New("app closed")
}
