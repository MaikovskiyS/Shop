package rest

import (
	"myproject/internal/domain"
	"strconv"
)

type SaveResponse struct {
	Id uint64 `json:"id"`
}
type SaveRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   string `json:"age"`
}

func (r *SaveRequest) toModel() (domain.User, error) {
	// if r == nil {
	// 	return domain.User{}, apperrors.New(apperrors.ErrBadRequest, "user-api-validateRequest-save", errors.New("bad params"))
	// }
	// if r.Name == "" {
	// 	return domain.User{}, apperrors.New(apperrors.ErrBadRequest, "user-api-validateRequest-save", errors.New("bad params"))
	// }
	age, err := strconv.Atoi(r.Age)
	if err != nil {
		ErrBadRequest.AddLocation("Validator-strconv.Atoi")
		ErrBadRequest.SetErr(err)
		return domain.User{}, ErrBadRequest
	}
	// if age <= 0 {
	// 	return domain.User{}, apperrors.New(apperrors.ErrBadRequest, "user-api-validateRequest-save", errors.New("bad params"))
	// }
	data := domain.User{
		Name:  r.Name,
		Email: r.Email,
		Age:   uint8(age),
	}
	return data, nil
}

type GetAllResponse struct {
	Result []domain.User `json:"result"`
}
