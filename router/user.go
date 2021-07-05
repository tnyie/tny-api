package router

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/tnyie/tny-api/views"
)

func userRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/{id}", views.GetUser)
	r.Post("/", views.PostUser)
	r.Put("/password", views.PasswordResetRequest)
	r.Put("/password/{token}", views.PasswordReset)
	return r
}
