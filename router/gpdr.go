package router

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/tnyie/tny-api/views"
)

func gdprRouter() http.Handler {
	r := chi.NewRouter()

	r.Get("/", views.GetUserGDPR)
	return r
}
