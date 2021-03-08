package router

import (
	"net/http"

	"github.com/tnyie/tny-api/views"
	"github.com/go-chi/chi"
)

func tokenRouter() http.Handler {
	r := chi.NewRouter()
	r.Post("/", views.CreateToken)
	r.Get("/", views.InspectToken)
	return r
}
