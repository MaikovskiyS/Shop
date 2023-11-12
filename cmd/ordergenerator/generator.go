package ordergenerator

import (
	"context"
	"myproject/cmd/ordergenerator/adapter/generator"
	"myproject/cmd/ordergenerator/adapter/msgbroker/producer"
	"myproject/cmd/ordergenerator/usecase"
	"myproject/internal/server/router"

	"github.com/segmentio/kafka-go"
)

// TODO: add transport
func New(ctx context.Context, c *kafka.Conn, r *router.Router) usecase.Generator {
	gen := generator.New()
	msg := producer.New(c)
	svc := usecase.New(gen, msg)
	// api := transport.New(ctx, svc)
	// api.RegisterRoutes(r)
	return svc
}
