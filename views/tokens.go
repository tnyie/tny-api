package views

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"github.com/tnyie/tny-api/middleware"
	"github.com/tnyie/tny-api/models"
)

// CreateToken creates an API token
func CreateToken(w http.ResponseWriter, r *http.Request) {
	jsonMap := make(map[string]string)

	bd, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading request body\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(bd, &jsonMap)
	if err != nil {
		log.Println("Couldn't unmarshall json body")
	}
	userAuth := &models.UserAuth{
		Email:    jsonMap["email"],
		Username: jsonMap["username"],
	}

	err = userAuth.VerifyPassword(jsonMap["password"])
	if err != nil {
		log.Println("Incorrect password")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

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
	log.Println(encoded)
	respondJSON(w, encoded, http.StatusCreated)
}

// InspectToken returns a status 200 if logged in, 403 if not
func InspectToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if claims, ok := ctx.Value(middleware.AuthCtx{}).(jwt.MapClaims); ok {
		log.Println(ctx.Value(middleware.AuthCtx{}))
		w.WriteHeader(http.StatusAccepted)
		jsonResp := make(map[string]interface{})
		jsonResp["user_id"] = claims["UserID"]
		encoded, err := json.Marshal(jsonResp)
		if err != nil {
			log.Println("Couldn't unmarshall response")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		respondJSON(w, encoded, http.StatusCreated)
	}
	log.Println(ctx.Value(middleware.AuthCtx{}))
	w.WriteHeader(http.StatusUnauthorized)
}
