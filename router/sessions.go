package router

import (
	"net/http"

	"github.com/gal/tny/views"
	"github.com/go-chi/chi"
)

func tokenRouter() http.Handler {
	r := chi.NewRouter()
	r.Post("/", views.CreateToken)
	r.Get("/", views.InspectToken)
	return r
}
