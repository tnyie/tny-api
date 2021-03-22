package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/spf13/viper"

	"github.com/tnyie/tny-api/config"
	"github.com/tnyie/tny-api/models"
	"github.com/tnyie/tny-api/router"

	"github.com/rs/cors"
)

func seed() {
	createUsers()
}

func main() {
	r := chi.NewRouter()

	config.InitConfig()
	router.Route(r)
	time.Sleep(time.Second)
	models.InitModels()

	handler := cors.AllowAll().Handler(r)

	if viper.GetString("debug") == "true" {
		go seed()
	}

	http.ListenAndServe(":8080", handler)
}
