package usecase

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"log"
	"myproject/internal/apperrors"
	"myproject/internal/domain"
)

const location = "Auth_Service-Usecase-"

var (
	ErrBadRequest = apperrors.New(apperrors.ErrBadRequest, location)
	ErrNotFound   = apperrors.New(apperrors.ErrNotFound, location)
)

type Auth interface {
	SignUp(ctx context.Context, user domain.User) error
	SignIn(ctx context.Context, inputUser domain.User) (string, error)
}

type usecase struct {
	user UserService
	auth TokenService
}

func New(t TokenService, u UserService) *usecase {
	return &usecase{
		user: u,
		auth: t,
	}
}

// need return value?
func (u *usecase) SignUp(ctx context.Context, user domain.User) error {

	_, err := u.user.GetByEmail(ctx, user.Email)
	if err != nil {
		var er *apperrors.AppErr
		if errors.As(err, &er) {
			if er.Type() == apperrors.ErrNotFound {
				log.Println("in my err")
				user.Password = generatePasswordHash(user.Password)
				id, err := u.user.Save(ctx, user)
				log.Println(id)
				if err != nil {
					log.Println("from save", err)
					return err
				}
				return nil
			}
		}
		return err
	}
	ErrBadRequest.AddLocation("SignUp-CheckUser")
	ErrBadRequest.SetErr(errors.New("user already exist"))
	return ErrBadRequest
}
func (u *usecase) SignIn(ctx context.Context, inputUser domain.User) (string, error) {
	user, err := u.user.GetByEmail(ctx, inputUser.Email)
	if err != nil {

		return "", err
	}
	if user == nil {
		ErrNotFound.AddLocation("SignInCheckUser")
		ErrNotFound.SetErr(errors.New("user not found"))
		return "", ErrNotFound
	}
	if generatePasswordHash(inputUser.Password) != user.Password {
		ErrBadRequest.AddLocation("SignIn-CheckPassword")
		ErrBadRequest.SetErr(errors.New("wrong password"))
		return "", ErrBadRequest
	}

	token, err := u.auth.CreateToken(user.Email)
	if err != nil {
		return "", err
	}

	return token, nil
}

const salt = "hjqrhjqw124617ajfhajs"

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
