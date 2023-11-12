package user_service

import (
	"myproject/internal/server/router"
	psql "myproject/internal/services/user_service/adapter/dbs/postgres"
	"myproject/internal/services/user_service/transport/rest"
	"myproject/internal/services/user_service/usecase"

	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

func New(r *router.Router, pCl *pgx.Conn, rCl *redis.Client) usecase.User {
	store := psql.New(pCl)
	svc := usecase.New(store)
	api := rest.New(svc)
	api.RegisterRoutes(r)
	return svc
}
