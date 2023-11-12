package rest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"myproject/internal/apperrors"
	"myproject/internal/services/product_service/usecase"
	"net/http"
	"strconv"
	"time"
)

const location = "Product_Service-Api-"

var (
	ErrBadRequest = apperrors.New(apperrors.ErrBadRequest, location)
	ErrInternal   = apperrors.New(apperrors.ErrInternal, location)
)

type api struct {
	timeout time.Duration
	product usecase.Product
}

func New(u usecase.Product) *api {
	return &api{
		timeout: time.Second * 3,
		product: u,
	}
}
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
	product, err := input.toModel()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), a.timeout)
	defer cancel()
	id, err := a.product.Save(ctx, product)
	if err != nil {
		return err
	}
	str := fmt.Sprintf("product created. ID: %v", id)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(str))
	return nil
}
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

	product, err := a.product.GetById(ctx, uint64(id))
	if err != nil {
		return err
	}
	resp := GetResponse{Result: product}
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
	ctx, cancel := context.WithTimeout(context.Background(), a.timeout)
	defer cancel()
	products, err := a.product.GetAll(ctx)
	if err != nil {
		return err
	}
	resp := GetAllRrsponse{Result: products}
	respBytes, err := json.Marshal(&resp)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)
	w.Write(respBytes)
	return nil
}
