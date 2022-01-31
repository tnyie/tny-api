package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/rs/cors"

	"github.com/tnyie/tny-api/config"
	"github.com/tnyie/tny-api/models"
	"github.com/tnyie/tny-api/router"
)

func main() {
	r := chi.NewRouter()

	config.InitConfig()
	router.Route(r)

	// wait for database
	time.Sleep(time.Second * 10)
	models.InitModels()

	handler := cors.AllowAll().Handler(r)

	log.Println("bruhe")
	http.ListenAndServe(":8080", handler)
}
