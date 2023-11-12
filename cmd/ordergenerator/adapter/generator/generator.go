package generator

import (
	"myproject/internal/domain"
	"time"
)

type generator struct {
}

func New() *generator {
	return &generator{}
}
func (g *generator) GenerateProduct() (*domain.Product, error) {
	p := &domain.Product{
		ID:       1,
		Sku:      "sku",
		Category: "category",
		Name:     "name",
		Price:    2.0,
		Image:    "image",
	}
	return p, nil
}
func (g *generator) GenerateOrder() (domain.Order, error) {

	products := make([]*domain.Product, 3)
	for i := 0; i < 3; i++ {
		pr, err := g.GenerateProduct()
		if err != nil {
			return domain.Order{}, err
		}
		products[i] = pr
	}
	o := domain.Order{
		ID:           1,
		UserID:       1,
		CustomerName: "cat name",
		TotalPrice:   3.0,
		CreatedAt:    time.Now().Add(time.Second * 2),
		Status:       "test",
		Products:     products,
	}
	return o, nil
}
