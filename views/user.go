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
	"github.com/tnyie/tny-api/middleware"
	"github.com/tnyie/tny-api/models"
)

const emailVerificationHeader = "Verification needed for Tny.ie"
const emailVerificationParagraph = `Welcome to Tny.ie!
Click the button below verify your email address.
<div style="display:flex;align-items:center;justify-content:center">
<a style="padding: 1em 2em; background-color: #009688; text-decoration: none; color: white; border-radius: 6%%;" href=" %s ">Verify Email</a>
</div>
If you did not make an account on https://tny.ie , please ignore this email.
`

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
	if err != nil {
		log.Println("Error creating user\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)

	expirationTime := time.Now().Add(time.Hour).Unix()
	claims := &models.EmailVerification{
		Key: user.Email,
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

	response, err := mail.SendMail(userAuth, "Email Verification for TnyIE", tokenString)
	if err != nil {
		log.Println("Error sending email to ", userAuth.Email)
		log.Println(err)
		return
	}
	log.Println("Sent email verification to ", userAuth.Username, " with response: ", response)
}
