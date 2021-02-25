package router

import (
	"net/http"

	"github.com/gal/tny/views"
	"github.com/go-chi/chi"
)

func linkRouter() http.Handler {
	r := chi.NewRouter()

	r.Get("/{id}", views.GetLink)
	r.Get("/{id}/{attr}", views.GetLinkAttribute)
	r.Put("/{id}/{attr}", views.PutLinkAttribute)
	r.Post("/", views.CreateLink)

	r.Get("/search/{query}", views.SearchLink)
	return r
}
