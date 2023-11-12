package product_service

import (
	"myproject/internal/server/router"
	"myproject/internal/services/product_service/adapter/storage"
	"myproject/internal/services/product_service/transport/rest"
	"myproject/internal/services/product_service/usecase"

	"github.com/jackc/pgx/v5"
)

func New(r *router.Router, pCl *pgx.Conn) usecase.Product {
	store := storage.New(pCl)
	svc := usecase.New(store)
	api := rest.New(svc)

	api.RegisterRoutes(r)
	return svc
}
