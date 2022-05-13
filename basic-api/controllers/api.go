package controllers

import (
	"basic-api/engine"
	"context"
	"fmt"
	"net/http"
)

type API struct {
	Logger func(http.Handler) http.Handler
	User   *engine.Route
}

func NewAPI() *API {
	return &API{
		Logger: engine.Logger,
	}
}

func (a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx = context.WithValue(ctx, engine.ContextOriginalPath, r.URL.Path)

	var next *engine.Route
	var head string
	head, r.URL.Path = engine.ShiftPath(r.URL.Path)
	if head == "user" {
		next = NewUser()
	} else {
		next = newError(fmt.Errorf("path not found"), http.StatusNotFound)
	}

	if next.Logger {
		next.Handler = a.Logger(next.Handler)
	}

	next.Handler.ServeHTTP(w, r.WithContext(ctx))
}

func newError(err error, statusCode int) *engine.Route {
	return &engine.Route{
		Logger: true,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			engine.Respond(w, r, statusCode, err)
		}),
	}
}
