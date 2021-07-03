package views

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/tnyie/tny-api/models"
	"github.com/tnyie/tny-api/util"
)

// GetUserGDPR returns a gdpr dump to the user
func GetUserGDPR(w http.ResponseWriter, r *http.Request) {
	log.Println("User fetching GDPR data")
	userAuth, authenticated := util.CheckLogin(r, "")

	if !authenticated || userAuth.UID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		log.Println("Failed to authenticate user")
		return
	}

	user := &models.User{UID: userAuth.UID}
	if err := user.Get(); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		log.Println("Failed to fetch user data")
		return
	}

	links, err := models.GetLinksByUser(user.UID)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		log.Println("Failed to get user links\n", err)
		return
	}

	gdprData := &models.GDPRData{
		UserData:     *user,
		UserAuthData: *userAuth,
		Links:        *links,
	}

	encoded, err := json.Marshal(gdprData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Failed to encode gdpr data")
		return
	}

	respondJSON(w, encoded, http.StatusOK)
}
