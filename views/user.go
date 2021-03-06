package views

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gal/tny/middleware"
	"github.com/gal/tny/models"
	"github.com/go-chi/chi"
)

func GetUser(w http.ResponseWriter, r *http.Request) {
	uid := ""
	if claims, ok := r.Context().Value(middleware.AuthCtx{}).(jwt.MapClaims); ok {
		uid = claims["UserID"].(string)
	}

	if uid == "" || uid != chi.URLParam(r, "id") {
		log.Println("User unauthorized to access resource")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	user := &models.User{
		UID: uid,
	}
	err := user.Get()
	if err != nil {
		log.Println("Error fetching user\n", err)
	}
	encoded, err := json.Marshal(user)
	if err != nil {
		log.Println("Error marshalling user object")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	respondJSON(w, encoded, http.StatusOK)
}

func PostUser(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]string)

	bd, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Couldn't read json Body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(bd, &data); err != nil {
		log.Println("Couldn't parse json body\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userAuth := &models.UserAuth{
		Username: data["username"],
		Email:    data["email"],
	}

	err = userAuth.Create(data["password"])
	if err != nil {
		log.Println("Error creating user\n", err)
		if userAuth.UID == "" {
			w.WriteHeader(http.StatusBadRequest)
		}
		w.WriteHeader(http.StatusInternalServerError)
	}

	user := &models.User{
		UID:      userAuth.UID,
		Username: userAuth.Username,
		Email:    userAuth.Email,
	}

	err = user.Create()

}
