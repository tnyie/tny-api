package router

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/tnyie/tny-api/views"
)

func linkRouter() http.Handler {
	r := chi.NewRouter()

	r.Get("/{slug}", views.GetLink)
	r.Get("/user/{id}", views.GetLinksByUser)
	r.Get("/authenticated/{slug}", views.GetAuthenticatedLink)
	// r.Get("/{id}/{attr}", views.GetLinkAttribute)
	r.Put("/{id}/{attr}", views.PutLinkAttribute)
	r.Post("/", views.CreateLink)
	r.Delete("/{id}", views.DeleteLink)

	r.Get("/search/{query}", views.SearchLink)
	return r
}
