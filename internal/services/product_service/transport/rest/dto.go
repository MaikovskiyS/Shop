package rest

import (
	"errors"
	"myproject/internal/domain"
	"strconv"
	"time"
)

type GetAllRrsponse struct {
	Result []*domain.Product `json:"products"`
}
type GetResponse struct {
	Result *domain.Product `json:"product"`
}
type SaveRequest struct {
	Sku       string    `json:"sku"`
	Category  string    `json:"category"`
	Name      string    `json:"name"`
	Price     string    `json:"price"`
	Image     string    `json:"image"`
	CreatedAt time.Time `json:"created_at"`
}

func (r *SaveRequest) toModel() (domain.Product, error) {
	if r.Name == "" || r.Category == "" || r.Price == "" {
		ErrBadRequest.AddLocation("SaveRequest-Tomodel")
		ErrBadRequest.SetErr(errors.New("name, category, price required"))
		return domain.Product{}, ErrInternal
	}
	price, err := strconv.ParseFloat(r.Price, 32)
	if err != nil {
		ErrBadRequest.AddLocation("SaveRequest-strconv.ParseFloat")
		ErrBadRequest.SetErr(errors.New("wrong price"))
		return domain.Product{}, ErrInternal
	}
	p := domain.Product{
		Sku:      r.Sku,
		Category: r.Category,
		Name:     r.Name,
		Price:    price,
		Image:    r.Image,
	}

	return p, nil
}
