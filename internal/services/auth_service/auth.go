package auth_service

import (
	"myproject/internal/server/router"
	"myproject/internal/services/auth_service/transport/rest"
	"myproject/internal/services/auth_service/usecase"

	"github.com/jackc/pgx/v5"
)

func New(r *router.Router, pCl *pgx.Conn, auth usecase.TokenService, u usecase.UserService) usecase.Auth {
	svc := usecase.New(auth, u)
	api := rest.New(svc)

	api.RegisterRoutes(r)
	return svc
}
