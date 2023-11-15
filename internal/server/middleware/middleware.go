package middleware

import (
	"errors"
	"fmt"
	"log"
	"myproject/internal/apperrors"
	"net/http"
	"strings"
	"time"
)

var (
	ErrBadRequest   = apperrors.New(apperrors.ErrBadRequest, "Server-Middleware-")
	ErrUnauthorized = apperrors.New(apperrors.ErrUnauthorized, "Server-Middleware-")
)

const (
	authorizationHeaderKey = "authorization"
	authorizationType      = "bearer"
)

type tokenService interface {
	VerifyToken(accessToken string) (bool, error)
}
type middleware struct {
	t tokenService
}
type AppHandler func(w http.ResponseWriter, r *http.Request) error

func New(ts tokenService) *middleware {
	return &middleware{t: ts}
}

// Auth middleware checking autorization header and verify bearer token
func (m *middleware) Auth(h AppHandler) AppHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		authorizationHeader := r.Header.Get(authorizationHeaderKey)
		isEmpty := len(authorizationHeader) == 0
		if isEmpty {
			ErrBadRequest.AddLocation("Auth-CheckHeader")
			ErrBadRequest.SetErr(errors.New("empty autorization header"))
			return ErrBadRequest
		}

		fields := strings.Fields(authorizationHeader)
		isValid := len(fields) == 2
		if !isValid {
			ErrBadRequest.AddLocation("Auth-CheckLenHeader")
			ErrBadRequest.SetErr(errors.New("wrong header len"))
			return ErrBadRequest
		}

		currentAuthorizationType := strings.ToLower(fields[0])
		if currentAuthorizationType != authorizationType {
			ErrBadRequest.AddLocation("Auth-CheckAuthorizationType")
			ErrBadRequest.SetErr(errors.New("wrong authorizationType"))
			return ErrBadRequest
		}

		accessToken := fields[1]
		_, err := m.t.VerifyToken(accessToken)
		if err != nil {
			ErrUnauthorized.AddLocation("Auth-VerifyToken")
			ErrUnauthorized.SetErr(err)
			return ErrUnauthorized
		}

		err = h(w, r)
		if err != nil {
			return err
		}
		return nil
	}
}

// Log request letency and answer from server
func (m *middleware) Logging(h AppHandler) AppHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		t := time.Now()
		err := h(w, r)
		log.Println(time.Since(t))
		if err != nil {
			var er *apperrors.AppErr
			if errors.As(err, &er) {
				log.Println(er.Log())
				return err
			}
			log.Printf("unknown error: %s", err)
			return err
		}
		log.Println("success")
		return nil
	}
}

// Handle Error and sending answer to client
func (m *middleware) ErrorHandle(h AppHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h(w, r)
		if err != nil {
			var er *apperrors.AppErr
			if errors.As(err, &er) {
				apperrors.ErrResponse(w, er)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("unknown error:  %s", err.Error())))
			return

		}
	}
}
