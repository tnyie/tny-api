package router

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/tnyie/tny-api/views"
)

func tokenRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/", views.InspectToken)
	r.Post("/", views.CreateToken)
	return r
}

func keyRouter() http.Handler {
	r := chi.NewRouter()
	r.Post("/", views.CreateAPIKey)
	return r
}
