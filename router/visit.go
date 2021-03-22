package router

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/tnyie/tny-api/views"
)

func visitRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/{id}/{start_date}", views.GetVisits)
	return r
}
