package rest

import (
	"myproject/internal/server/router"
)

func (a *api) RegisterRoutes(r *router.Router) {
	r.HandleFunc("/products", r.ErrorHandle(r.Logging(r.Auth(a.GetById))))
	r.HandleFunc("/product", r.ErrorHandle(r.Logging(r.Auth(a.Save))))
	r.HandleFunc("/products/all", r.ErrorHandle(r.Logging(r.Auth(a.GetAll))))
}