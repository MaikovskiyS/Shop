package apperrors

import (
	"fmt"
	"net/http"
)

type errType string

var (
	ErrBadRequest   errType = "badRequest"
	ErrNotFound     errType = "notFound"
	ErrInternal     errType = "internalErr"
	ErrUnauthorized errType = "Unauthorized"
)

type AppErr struct {
	t    errType
	err  error
	loc  string
	pkg  string
	code int
}

func New(t errType, pkg string) *AppErr {
	e := &AppErr{t: t, pkg: pkg}
	switch t {
	case ErrBadRequest:
		e.code = http.StatusBadRequest
	case ErrInternal:
		e.code = http.StatusInternalServerError
	case ErrNotFound:
		e.code = http.StatusNotFound
	case ErrUnauthorized:
		e.code = http.StatusUnauthorized
	}
	return e
}
func (e *AppErr) Error() string {
	return e.err.Error()
}
func (e *AppErr) AddLocation(l string) {
	e.loc = l
}
func (e *AppErr) Type() errType {
	return e.t
}
func (e *AppErr) Code() int {
	return e.code
}
func (e *AppErr) Log() string {
	return fmt.Sprintf("%s  Error: %s", e.pkg+e.loc, e.err.Error())
}
func (e *AppErr) SetErr(er error) {
	e.err = er
}

func (e *AppErr) Unwrap() error { return e.err }
