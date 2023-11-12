package rest

import (
	"errors"
	"myproject/internal/domain"
	"strconv"
)

type signUpRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Age      string `json:"age"`
}

func (s *signUpRequest) toModel() (domain.User, error) {
	if s.Email == "" || s.Password == "" || s.Name == "" {
		ErrBadRequest.AddLocation("SignUp-Tomodel-Validate")
		ErrBadRequest.SetErr(errors.New("email or password, or name cant be empty"))
		return domain.User{}, ErrBadRequest
	}
	age, err := strconv.Atoi(s.Age)
	if err != nil {
		ErrBadRequest.AddLocation("SignUp-Tomodel-Validate")
		ErrBadRequest.SetErr(errors.New("wrong age"))
		return domain.User{}, ErrBadRequest
	}
	user := domain.User{
		Id:       0,
		Name:     s.Name,
		Email:    s.Email,
		Password: s.Password,
		Age:      uint8(age),
	}
	return user, nil
}

type signinResponse struct {
	Msg   string `json:"msg"`
	Token string `json:"token"`
}
type signInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *signInRequest) toModel() (domain.User, error) {
	if s.Email == "" || s.Password == "" {
		ErrBadRequest.AddLocation("SignIn-Tomodel-Validate")
		ErrBadRequest.SetErr(errors.New("wrong age"))
		return domain.User{}, ErrBadRequest
	}
	user := domain.User{
		Email:    s.Email,
		Password: s.Password,
	}
	return user, nil
}
