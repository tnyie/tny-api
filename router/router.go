package router

import (
	"net/http"

	"github.com/go-chi/chi"
	mid "github.com/go-chi/chi/middleware"

	"github.com/gal/tny/middleware"
	"github.com/gal/tny/views"
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
	r.Use(middleware.FirebaseAuth)
	r.Mount("/link", linkRouter())
	return r
}
