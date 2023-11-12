package rest

import (
	"myproject/internal/server/router"
)

func (a *api) RegisterRoutes(r *router.Router) {
	r.HandleFunc("/users", r.ErrorHandle(r.Logging(r.Auth(a.GetAll))))
	r.HandleFunc("/user", r.ErrorHandle(r.Logging(r.Auth(a.Save))))
}
