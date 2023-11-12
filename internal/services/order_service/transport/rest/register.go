package rest

import (
	"myproject/internal/server/router"
)

func (a *api) RegisterRoutes(r *router.Router) {
	r.HandleFunc("/orders", r.ErrorHandle(r.Logging(r.Auth(a.GetById))))
	r.HandleFunc("/order", r.ErrorHandle(r.Logging(r.Auth(a.Save))))
}
