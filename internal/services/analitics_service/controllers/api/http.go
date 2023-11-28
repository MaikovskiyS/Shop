package api

import (
	"context"
	"encoding/json"
	"log"
	"myproject/internal/domain"
	"net/http"
)

type Service interface {
	GetAll(ctx context.Context) ([]*domain.Order, error)
}
type api struct {
	usecase Service
}

func New(s Service) *api {
	return &api{
		usecase: s,
	}
}

func (a *api) GetAll(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	orders, err := a.usecase.GetAll(context.Background())
	if err != nil {
		log.Println("getAllOrder err", err)
		return err
	}
	resp := GetAllResponse{Result: orders}
	respBytes, err := json.Marshal(&resp)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)
	w.Write(respBytes)
	return nil
}
