package views

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"github.com/tnyie/tny-api/models"
	"github.com/tnyie/tny-api/util"
)

// InspectToken returns a status 200 if logged in, 403 if not
func InspectToken(w http.ResponseWriter, r *http.Request) {
	if user, valid, admin := util.CheckLogin(r, ""); valid {
		w.WriteHeader(http.StatusAccepted)
		jsonResp := make(map[string]interface{})
		jsonResp["user_id"] = user.UID
		if admin {
			jsonResp["admin"] = admin
		}
		encoded, err := json.Marshal(jsonResp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("Couldn't unmarshall response")
			return
		}
		respondJSON(w, encoded, http.StatusOK)
	}

	w.WriteHeader(http.StatusUnauthorized)
}

// CreateToken creates an API token
func CreateToken(w http.ResponseWriter, r *http.Request) {

	userAuth := getLogin(r)

	expirationTime := time.Now().Add(time.Hour * 5)

	claims := &models.JWTClaims{
		UserID: userAuth.UID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(viper.GetString("tny.auth.key")))
	if err != nil {
		log.Println("Couldn't sign jwt\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	jsonResp := make(map[string]interface{})
	jsonResp["token"] = tokenString
	encoded, err := json.Marshal(jsonResp)
	if err != nil {
		log.Println("Couldn't unmarshall token response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	respondJSON(w, encoded, http.StatusCreated)
}

func CreateAPIKey(w http.ResponseWriter, r *http.Request) {
	user := getLogin(r)
	if user != nil {
		key, err := models.GenerateAPIKey(user.UID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("Couldn't create api key")
			return
		}

		jsonResp := make(map[string]interface{})
		jsonResp["key"] = key.ID
		encoded, err := json.Marshal(jsonResp)
		if err != nil {
			log.Println("Couldn't unmarshall token response")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		respondJSON(w, encoded, http.StatusCreated)
		return
	}
	w.WriteHeader(http.StatusUnauthorized)
}

func getLogin(r *http.Request) *models.UserAuth {
	jsonMap := make(map[string]string)

	bd, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading request body\n", err)
		return nil
	}
	err = json.Unmarshal(bd, &jsonMap)
	if err != nil {
		log.Println("Couldn't unmarshall json body")
		return nil
	}
	userAuth := &models.UserAuth{
		Email:    jsonMap["email"],
		Username: jsonMap["username"],
	}

	err = userAuth.VerifyPassword(jsonMap["password"])
	if err != nil {
		log.Println("Incorrect password")
		return nil
	}

	return userAuth
}
