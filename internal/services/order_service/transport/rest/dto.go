package rest

import (
	"errors"
	"myproject/internal/domain"
	"myproject/internal/services/order_service/model"
	"strconv"
)

type GetByIdResponse struct {
	Result *domain.Order `json:"order"`
}
type SaveRequest struct {
	UserId      string   `json:"user_id"`
	ProductsIds []string `json:"products_ids"`
}

func (r *SaveRequest) toModel() (model.Order, error) {
	if r.UserId == "" {
		ErrBadRequest.AddLocation("SaveRequestValidator")
		ErrBadRequest.SetErr(errors.New("user_id cant be empty"))
		return model.Order{}, ErrBadRequest
	}
	userId, err := strconv.Atoi(r.UserId)
	if err != nil {
		ErrBadRequest.AddLocation("SaveRequest-strconv.Atoi")
		ErrBadRequest.SetErr(err)
		return model.Order{}, ErrBadRequest
	}
	productIds := make([]uint64, 0)
	for _, strId := range r.ProductsIds {
		pId, err := strconv.Atoi(strId)
		if err != nil {
			ErrBadRequest.AddLocation("SaveRequest-strconv.Atoi")
			ErrBadRequest.SetErr(err)
			return model.Order{}, ErrBadRequest
		}
		productIds = append(productIds, uint64(pId))
	}
	or := model.Order{
		UserId:      uint64(userId),
		ProductsIds: productIds,
	}
	return or, nil
}

type GetAllResponse struct {
	Rows   uint64          `json:"rows"`
	Result []*domain.Order `json:"orders"`
}
