package transport

import "myproject/internal/server/router"

func (a *api) RegisterRoutes(r *router.Router) {
	r.HandleFunc("/gen/start", a.Start)
}
