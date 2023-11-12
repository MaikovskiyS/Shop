package consumer

import (
	"context"
	"encoding/json"
	"log"
	"myproject/internal/apperrors"
	"myproject/internal/services/order_service/model"
	"myproject/internal/services/order_service/usecase"
	"time"

	"github.com/segmentio/kafka-go"
)

const location = "Order_Service-Consumer-"

var (
	ErrBadRequest = apperrors.New(apperrors.ErrBadRequest, location)
	ErrInternal   = apperrors.New(apperrors.ErrInternal, location)
)

type consumer struct {
	timeout time.Duration
	reader  *kafka.Reader
	order   usecase.Order
}

func New(reader *kafka.Reader, o usecase.Order) *consumer {
	return &consumer{
		timeout: time.Second * 30,
		reader:  reader,
		order:   o,
	}
}
func (c *consumer) SaveOrders(ctx context.Context) error {
	for {
		ctxt, cancel := context.WithTimeout(ctx, c.timeout)
		defer cancel()
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:

			msg, err := c.reader.ReadMessage(ctxt)
			if err != nil {
				return err
			}
			err = c.reader.CommitMessages(ctx, msg)
			if err != nil {
				return err
			}
			var do Order
			err = json.Unmarshal(msg.Value, &do)
			if err != nil {
				return err
			}
			log.Println("DOMAINORDER: ", do)
			mo := model.Order{UserId: do.UserID, ProductsIds: make([]uint64, len(do.Products))}
			for i, product := range do.Products {
				mo.ProductsIds[i] = product.ID
			}
			id, err := c.order.Save(ctxt, mo)
			if err != nil {
				return err
			}

			log.Printf("order created. ID: %v", id)
		}
	}
}
