package msgbroker

import (
	"context"
	"encoding/json"
	"log"
	"myproject/internal/apperrors"
	"myproject/internal/domain"
	"time"

	"github.com/segmentio/kafka-go"
)

const location = "Order_Service-Consumer-"

var (
	ErrBadRequest = apperrors.New(apperrors.ErrBadRequest, location)
	ErrInternal   = apperrors.New(apperrors.ErrInternal, location)
)

type Service interface {
	SaveOrder(ctx context.Context, or *domain.Order) (uint64, error)
	GetAll(ctx context.Context) ([]*domain.Order, error)
}
type consumer struct {
	timeout time.Duration
	reader  *kafka.Reader
	order   Service
}

func New(reader *kafka.Reader, o Service) *consumer {
	return &consumer{
		timeout: time.Second * 5,
		reader:  reader,
		order:   o,
	}
}

// Save orders
func (c *consumer) SaveOrders(ctx context.Context) error {
	for {
		ctxt, cancel := context.WithTimeout(ctx, c.timeout)
		defer cancel()
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:

			msg, err := c.reader.ReadMessage(ctx)
			if err != nil {
				log.Println("read msg err")
				return err
			}
			err = c.reader.CommitMessages(ctx, msg)
			if err != nil {
				return err
			}
			var order *domain.Order
			err = json.Unmarshal(msg.Value, &order)
			if err != nil {
				log.Println("unmarshal err")
				return err
			}

			_, err = c.order.SaveOrder(ctxt, order)
			if err != nil {
				return err
			}
		}
	}
}
