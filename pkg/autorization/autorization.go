package autorization

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
)

/**
 * PasetoToken implements port.TokenService interface
 * and provides an access to the paseto library
 */
const tokenTTL = 12 * time.Hour

type jwtAuth struct {
	expairesAt *jwt.Time
	secretKey  []byte
}
type tokenClaims struct {
	jwt.StandardClaims
	UserEmail string `json:"user_email"`
}

// NewToken creates a new paseto instance
func New() (*jwtAuth, error) {

	return &jwtAuth{
		expairesAt: jwt.NewTime(float64(time.Now().Add(tokenTTL).Unix())),
		secretKey:  []byte("mysecretkey"),
	}, nil
}

// CreateToken creates a new paseto token
func (pt *jwtAuth) CreateToken(email string) (string, error) {

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{jwt.StandardClaims{
		ExpiresAt: pt.expairesAt,
	}, email})
	token, err := jwtToken.SignedString(pt.secretKey)
	if err != nil {
		return "", fmt.Errorf("jwt-createToken Err: %w", err)
	}
	return token, nil

}

// VerifyToken verifies the paseto token
func (pt *jwtAuth) VerifyToken(accessToken string) (bool, error) {

	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return pt.secretKey, nil
	})
	if err != nil {
		return false, err
	}
	_, ok := token.Claims.(*tokenClaims)
	if !ok {
		return false, errors.New("token claims are not of type *tokenClaims")
	}

	return true, nil

}
