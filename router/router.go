package router

import (
	"net/http"

	"github.com/go-chi/chi"
	mid "github.com/go-chi/chi/middleware"

	"github.com/tnyie/tny-api/middleware"
	"github.com/tnyie/tny-api/views"
)

// Route sets up routes
func Route(r *chi.Mux) {
	r.Use(mid.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://tny.ie/ui", http.StatusTemporaryRedirect)
	})
	r.Get("/{slug}", views.GetLink)
	r.Mount("/api", apiHandler())
}

func apiHandler() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.CheckToken)
	r.Mount("/tokens", tokenRouter())
	r.Mount("/verify", verificationRouter())
	r.Mount("/links", linkRouter())
	r.Mount("/users", userRouter())
	return r
}
