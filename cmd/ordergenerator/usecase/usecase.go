package usecase

import (
	"context"
	"errors"
	"time"
)

type Generator interface {
	Generate(ctx context.Context) error
}
type usecase struct {
	timeout time.Duration
	sender  Sender
	gen     OrderGenerator
}

func New(g OrderGenerator, s Sender) Generator {
	return &usecase{
		timeout: time.Second * 5,
		gen:     g,
		sender:  s,
	}
}
func (u *usecase) Generate(ctx context.Context) error {
	ticker := 0
MainLoop:
	for {
		select {
		case <-ctx.Done():
			break MainLoop
		default:

			order, err := u.gen.GenerateOrder()
			if err != nil {
				return err
			}
			err = u.sender.Send(order)
			if err != nil {
				return err
			}

		}
		time.Sleep(3 * time.Second)
		if ticker == 2 {
			return errors.New("ticker end")
		}
		ticker++
	}
	return ctx.Err()
}
