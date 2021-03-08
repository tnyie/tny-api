package main

import (
	"net/http"

	"github.com/go-chi/chi"

	"github.com/tnyie/tny-api/config"
	"github.com/tnyie/tny-api/models"
	"github.com/tnyie/tny-api/router"

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
