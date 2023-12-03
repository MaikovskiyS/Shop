package rest

import (
	"myproject/internal/server/router"
)

// Register product routes
func (a *api) RegisterRoutes(r *router.Router) {
	r.HandleFunc("/products", r.ErrorHandle(r.Metrics(r.Logging(r.Auth(a.GetById)))))
	r.HandleFunc("/product", r.ErrorHandle(r.Metrics(r.Logging(r.Auth(a.Save)))))
	r.HandleFunc("/products/all", r.ErrorHandle(r.Metrics(r.Logging(r.Auth(a.GetAll)))))
}
