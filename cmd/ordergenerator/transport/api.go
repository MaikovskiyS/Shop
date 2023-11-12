package transport

import (
	"context"
	"myproject/cmd/ordergenerator/usecase"
	"net/http"
)

type api struct {
	order   usecase.Generator
	exitCtx context.Context
}

func New(c context.Context, o usecase.Generator) *api {
	return &api{
		exitCtx: c,
		order:   o,
	}
}
func (a *api) Start(w http.ResponseWriter, r *http.Request) {
	go func() {
		err := a.order.Generate(a.exitCtx)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
	}()
	w.Write([]byte("start generator"))
}
