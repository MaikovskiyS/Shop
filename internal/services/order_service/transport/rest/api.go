package rest

import (
	"context"
	"encoding/json"
	"errors"
	"myproject/internal/apperrors"
	"myproject/internal/services/order_service/usecase"
	"net/http"
	"strconv"
	"time"
)

const location = "Order_Service-Api-"

var (
	ErrBadRequest = apperrors.New(apperrors.ErrBadRequest, location)
	ErrInternal   = apperrors.New(apperrors.ErrInternal, location)
)

type api struct {
	timeout time.Duration
	order   usecase.Order
}

func New(o usecase.Order) *api {
	return &api{
		timeout: time.Second * 3,
		order:   o,
	}
}

// Validate input params from request; mapping to domain entity; save Order
func (a *api) Save(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		ErrBadRequest.AddLocation("Save-CheckMethod")
		ErrBadRequest.SetErr(errors.New("wrong method"))
		return ErrBadRequest
	}
	input := &SaveRequest{}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		ErrBadRequest.AddLocation("Save-json.Decode")
		ErrBadRequest.SetErr(err)
		return ErrBadRequest
	}
	r.Body.Close()
	order, err := input.toModel()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), a.timeout)
	defer cancel()
	_, err = a.order.Save(ctx, order)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("order created"))
	return nil
}

// Validate input params from request; mapping to domain entity; get Order by id
func (a *api) GetById(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		ErrBadRequest.AddLocation("GetByID-CheckMethod")
		ErrBadRequest.SetErr(errors.New("wrong method"))
		return ErrBadRequest
	}
	r.ParseForm()
	data := r.Form.Get("id")
	if data == "" {
		ErrBadRequest.AddLocation("GetByID-ValidateId")
		ErrBadRequest.SetErr(errors.New("id required"))
		return ErrBadRequest
	}
	id, err := strconv.Atoi(data)
	if err != nil {
		ErrBadRequest.AddLocation("GetByID-strconv.Atoi")
		ErrBadRequest.SetErr(err)
		return ErrBadRequest
	}
	ctx, cancel := context.WithTimeout(context.Background(), a.timeout)
	defer cancel()

	order, err := a.order.GetById(ctx, uint64(id))
	if err != nil {
		return err
	}
	resp := GetByIdResponse{Result: order}
	respBytes, err := json.Marshal(&resp)
	if err != nil {
		ErrInternal.AddLocation("GetByID-json.Marshal")
		ErrInternal.SetErr(err)
		return ErrInternal
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respBytes)
	return nil
}
func (a *api) GetAll(w http.ResponseWriter, r *http.Request) error {

	if r.Method != http.MethodGet {
		ErrBadRequest.AddLocation("GetAll-CheckMethod")
		ErrBadRequest.SetErr(errors.New("wrong method"))
		return ErrBadRequest
	}
	orders, err := a.order.GetAll(r.Context())
	if err != nil {
		return err
	}

	if orders == nil {
		ErrInternal.AddLocation("GetAll-CheckNilOrders")
		ErrInternal.SetErr(errors.New("cant get orders"))
		return ErrInternal
	}
	resp := &GetAllResponse{Rows: uint64(len(orders)), Result: orders}

	oBytes, err := json.Marshal(resp)
	if err != nil {
		ErrInternal.AddLocation("GetAll-json.Marshal")
		ErrInternal.SetErr(err)
		return ErrInternal
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(oBytes)
	return nil
}
