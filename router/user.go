package router

import (
	"net/http"

	"github.com/gal/tny/views"
	"github.com/go-chi/chi"
)

func userRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/{id}", views.GetUser)
	r.Post("/", views.PostUser)
	return r
}
