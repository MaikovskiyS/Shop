package rest

import (
	"myproject/internal/server/router"
)

func (a *api) RegisterRoutes(r *router.Router) {
	r.HandleFunc("/orders/{id}", r.ErrorHandle(r.Logging(r.Auth(a.GetById))))
	r.HandleFunc("/order", r.ErrorHandle(r.Logging(r.Auth(a.Save))))
	r.HandleFunc("/orders/all", r.ErrorHandle(r.Spammer(r.Logging(r.Auth(a.GetAll)))))
}
