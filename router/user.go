package router

import (
	"net/http"

	"github.com/tnyie/tny-api/views"
	"github.com/go-chi/chi"
)

func userRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/{id}", views.GetUser)
	r.Post("/", views.PostUser)
	return r
}
