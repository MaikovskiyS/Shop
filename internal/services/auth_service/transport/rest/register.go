package rest

import (
	"myproject/internal/server/router"
)

func (a *api) RegisterRoutes(r *router.Router) {
	r.HandleFunc("/sign_up", r.ErrorHandle(r.Logging((a.SignUp))))
	r.HandleFunc("/sign_in", r.ErrorHandle(r.Logging((a.SignIn))))

}
