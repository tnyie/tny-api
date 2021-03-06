package main

import (
	"net/http"

	"github.com/go-chi/chi"

	"github.com/gal/tny/config"
	"github.com/gal/tny/models"
	"github.com/gal/tny/router"

	"github.com/rs/cors"
)

func main() {
	r := chi.NewRouter()

	config.InitConfig()
	router.Route(r)
	models.InitModels()

	handler := cors.AllowAll().Handler(r)

	http.ListenAndServe(":8080", handler)
}
