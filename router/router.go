package router

import (
	"net/http"

	"github.com/go-chi/chi"
	mid "github.com/go-chi/chi/middleware"
	"github.com/spf13/viper"

	"github.com/tnyie/tny-api/middleware"
	"github.com/tnyie/tny-api/oidc"
	"github.com/tnyie/tny-api/views"
)

// Route sets up routes
func Route(r *chi.Mux) {
	r.Use(mid.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://"+viper.GetString("tny.ui.url"), http.StatusTemporaryRedirect)
	})
	r.Get("/{slug}", views.RedirectSlug)
	r.Mount("/api", apiHandler())
}

func apiHandler() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.CheckToken)
	r.Mount("/gdpr", gdprRouter())
	r.Mount("/tokens", tokenRouter())
	r.Mount("/keys", keyRouter())
	r.Mount("/verify", verificationRouter())
	r.Mount("/links", linkRouter())
	r.Mount("/users", userRouter())
	r.Mount("/visits", visitRouter())

	r.Get("/auth/callback", oidc.HandleCallback)
	return r
}
