package producer

import (
	"encoding/json"
	"log"
	"myproject/internal/domain"

	"github.com/segmentio/kafka-go"
)

type msgbroker struct {
	conn *kafka.Conn
}

func New(c *kafka.Conn) *msgbroker {
	return &msgbroker{
		conn: c,
	}
}

func (m *msgbroker) Send(o *domain.Order) error {
	oBytes, err := json.Marshal(&o)
	if err != nil {
		return err
	}
	msg := kafka.Message{
		Value: oBytes,
	}
	_, err = m.conn.WriteMessages(msg)
	if err != nil {
		log.Println("writemsg err", err)
		return err
	}
	log.Println("generated order send")
	return nil
}
