package main

import (
	"net/http"

	"github.com/go-chi/chi"

	"github.com/gal/tny/config"
	"github.com/gal/tny/models"
	"github.com/gal/tny/router"
)

func main() {
	r := chi.NewRouter()

	config.InitConfig()
	router.Route(r)
	models.InitModels()
	

	http.ListenAndServe(":8080", r)
}
