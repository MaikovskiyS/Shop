package rest

import (
	"context"
	"encoding/json"
	"errors"
	"myproject/internal/apperrors"
	"myproject/internal/services/user_service/usecase"
	"net/http"
	"time"
)

//go:generate mockgen -source=handler.go -destination=mocks/mock.go
const location = "User_Service-Api"

var (
	ErrBadRequest = apperrors.New(apperrors.ErrBadRequest, location)
	ErrInternal   = apperrors.New(apperrors.ErrInternal, location)
)

type api struct {
	svc usecase.User
}

var timeout = time.Second * 5

func New(svc usecase.User) *api {
	return &api{svc: svc}
}

// Validate input params from request; mapping to domain entity; save User
func (a *api) Save(w http.ResponseWriter, r *http.Request) error {
	var request *SaveRequest
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		ErrBadRequest.AddLocation("Save-json.Decode")
		ErrBadRequest.SetErr(err)
		return ErrBadRequest
	}
	params, err := request.toModel()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	id, err := a.svc.Save(ctx, params)
	if err != nil {
		return err
	}
	response := SaveResponse{Id: id}
	resBytes, err := json.Marshal(&response)
	if err != nil {
		ErrInternal.AddLocation("Save-json.Marshal")
		ErrInternal.SetErr(err)
		return ErrInternal
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(resBytes)
	return nil
}

// Get all Users
func (a *api) GetAll(w http.ResponseWriter, r *http.Request) error {
	w.Header().Add("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		ErrBadRequest.AddLocation("GetAll-CheckMethod")
		ErrBadRequest.SetErr(errors.New("wrong method"))
		return ErrBadRequest
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	users, err := a.svc.GetAll(ctx)
	if err != nil {
		return err
	}
	resp := &GetAllResponse{Result: users}
	respBytes, err := json.Marshal(&resp)
	if err != nil {
		ErrInternal.AddLocation("GetAll-json.Marshal")
		ErrInternal.SetErr(err)
		return ErrInternal
	}
	w.WriteHeader(http.StatusOK)
	w.Write(respBytes)
	return nil
}
