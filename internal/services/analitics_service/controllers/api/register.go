package api

import "myproject/internal/server/router"

func (a *api) RegisterRoutes(r *router.Router) {
	r.HandleFunc("/orders/all", r.ErrorHandle(r.Metrics(r.Logging(r.Auth(a.GetAll)))))
}
