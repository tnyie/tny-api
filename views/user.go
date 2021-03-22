package views

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
	"github.com/spf13/viper"

	"github.com/tnyie/tny-api/mail"
	"github.com/tnyie/tny-api/models"
	"github.com/tnyie/tny-api/util"
)

// GetUser gets user data
func GetUser(w http.ResponseWriter, r *http.Request) {

	uid := chi.URLParam(r, "id")

	_, authorized := util.CheckLogin(r, uid)
	if !authorized {
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

// PostUser creates a user
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
	if err != nil {
		log.Println("Error creating user\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)

	expirationTime := time.Now().Add(time.Hour).Unix()
	claims := &models.EmailVerification{
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime,
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(
		viper.GetString("tny.auth.key"),
	))
	if err != nil {
		log.Println("Couldn't create token for email verification")
		return
	}

	err = mail.SendMail(userAuth, tokenString)
	if err != nil {
		log.Println("Error sending email to ", userAuth.Email)
		log.Println(err)
		return
	}
	log.Println("Sent email verification to ", userAuth.Username)
}

// // ResetPassword sends email with a link
// func ResetPassword(w http.ResponseWriter, r *http.Request) {
// 	userAuth, authorized := util.CheckLogin(r, chi.URLParam(r, "id"))
// 	if authorized && userAuth != nil {

// 	}
// }
