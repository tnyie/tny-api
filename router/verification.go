package router

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/tnyie/tny-api/views"
)

func verificationRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/{token}", views.VerifyEmail)
	return r
}
