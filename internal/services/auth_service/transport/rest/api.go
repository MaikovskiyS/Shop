package rest

import (
	"context"
	"encoding/json"
	"myproject/internal/apperrors"
	"myproject/internal/services/auth_service/usecase"

	"net/http"
	"time"
)

const location = "Auth_Service-Api-"

var (
	ErrBadRequest = apperrors.New(apperrors.ErrBadRequest, location)
	ErrInternal   = apperrors.New(apperrors.ErrInternal, location)
)

type api struct {
	timeout time.Duration
	auth    usecase.Auth
}

func New(u usecase.Auth) *api {
	return &api{
		timeout: time.Second * 3,
		auth:    u,
	}
}
func (a *api) SignUp(w http.ResponseWriter, r *http.Request) error {
	var input *signUpRequest
	w.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		ErrBadRequest.AddLocation("SignUp-Decode")
		ErrBadRequest.SetErr(err)
		return ErrBadRequest
	}
	r.Body.Close()
	user, err := input.toModel()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), a.timeout)
	defer cancel()
	err = a.auth.SignUp(ctx, user)
	if err != nil {
		return err
	}
	respBytes, err := json.Marshal("success")
	if err != nil {
		ErrInternal.AddLocation("SignUp-json.Marshal")
		ErrInternal.SetErr(err)
		return ErrInternal
	}
	w.Write(respBytes)
	return nil
}
func (a *api) SignIn(w http.ResponseWriter, r *http.Request) error {

	var input *signInRequest
	w.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		ErrBadRequest.AddLocation("SignIn-Decode")
		ErrBadRequest.SetErr(err)
		return ErrBadRequest
	}
	r.Body.Close()
	user, err := input.toModel()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), a.timeout)
	defer cancel()
	token, err := a.auth.SignIn(ctx, user)
	if err != nil {
		return err
	}
	resp := &signinResponse{Msg: "login success", Token: token}
	respBytes, err := json.Marshal(&resp)
	if err != nil {
		ErrInternal.AddLocation("SignIn-json.Marshal")
		ErrInternal.SetErr(err)
		return ErrInternal
	}
	w.Write(respBytes)
	return nil
}
