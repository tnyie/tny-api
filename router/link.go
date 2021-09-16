package router

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/tnyie/tny-api/views"
)

func linkRouter() http.Handler {
	r := chi.NewRouter()

	r.Get("/link/{id}", views.GetLink)
	r.Get("/{slug}", views.GetLinkFromSlug)

	r.Get("/user/{id}", views.GetLinksByUser)
	r.Put("/authenticated/{slug}", views.GetAuthenticatedLink)
	r.Put("/{id}/{attr}", views.PutLinkAttribute)
	r.Put("/{id}", views.UpdateLinkLease)
	r.Post("/", views.CreateLink)
	r.Delete("/{id}", views.DeleteLink)
	return r
}
