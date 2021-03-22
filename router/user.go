package router

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/tnyie/tny-api/views"
)

func userRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/{id}", views.GetUser)
	// r.Put("/{id}/password", views.ResetPassword)
	r.Post("/", views.PostUser)
	return r
}
