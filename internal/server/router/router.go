package router

import (
	"myproject/internal/server/middleware"
	"net/http"
)

type Middle interface {
	Auth(h middleware.AppHandler) middleware.AppHandler
	Logging(h middleware.AppHandler) middleware.AppHandler
	ErrorHandle(h middleware.AppHandler) http.HandlerFunc
	Spammer(h middleware.AppHandler) middleware.AppHandler
}
type tokenService interface {
	VerifyToken(accessToken string) (bool, error)
}
type Router struct {
	Middle
	*http.ServeMux
}

func New(t tokenService) *Router {
	mux := http.NewServeMux()
	m := middleware.New(t)

	return &Router{m, mux}
}
