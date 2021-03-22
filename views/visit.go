package views

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/tnyie/tny-api/models"
	"github.com/tnyie/tny-api/util"
)

func GetVisits(w http.ResponseWriter, r *http.Request) {
	link := &models.Link{ID: chi.URLParam(r, "id")}

	if link.ID == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	link.Get()

	if _, authorized := util.CheckLogin(r, link.OwnerID); !authorized {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	start_date := chi.URLParam(r, "start_date")

	visits := models.GetVisits(link.ID, start_date)
	encoded, err := json.Marshal(visits)
	if err != nil {
		log.Println("error marshalling visits\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	respondJSON(w, encoded, http.StatusOK)
}
