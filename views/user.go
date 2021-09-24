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

	_, authorized, admin := util.CheckLogin(r, uid)
	if !authorized && !admin {
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

	if !viper.GetBool("tny.self.signup") {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

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

	sendEmailVerification(userAuth)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error sending email to ", userAuth.Email)
		log.Println(err)
		return
	}
	log.Println("Sent email verification to ", userAuth.Username)
	w.WriteHeader(http.StatusCreated)
}

// Send email Verification
func VerifyEmail(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")

	if email == "" {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("No email in request")
		return
	}

	userAuth := &models.UserAuth{
		Email: email,
	}

	err := userAuth.GetByEmail()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = sendEmailVerification(userAuth)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func sendEmailVerification(userAuth *models.UserAuth) error {
	expirationTime := time.Now().Add(time.Hour).Unix()
	claims := &models.EmailVerification{
		Email: userAuth.Email,
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
		return err
	}

	err = mail.SendMailVerification(userAuth, tokenString)
	if err != nil {
		log.Println("Error sending email to ", userAuth.Email)
		return err
	}
	log.Println("sent email verification")
	return nil
}

// PasswordResetRequest sends an email with a link with token used in PasswordReset
func PasswordResetRequest(w http.ResponseWriter, r *http.Request) {
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

	if data["email"] == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userAuth := &models.UserAuth{Email: data["email"]}

	if err := userAuth.Get(); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
	}

	expirationTime := time.Now().Add(time.Hour).Unix()

	claims := &models.EmailVerification{
		Email: data["email"],
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

	err = mail.SendPasswordVerification(userAuth, tokenString)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error sending email to ", userAuth.Email)
		log.Println(err)
		return
	}
	log.Println("Sent email verification to ", userAuth.Username)
	w.WriteHeader(http.StatusOK)
}

// PasswordReset resets password from signed jwt in 'token' url param
func PasswordReset(w http.ResponseWriter, r *http.Request) {
	tokenString := chi.URLParam(r, "token")

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

	if data["password"] == "" {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Missing password in json body")
		return
	}

	log.Println("tokenString : ", tokenString)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(viper.GetString("tny.auth.key")), nil
	})

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		log.Println("Error parsing verification token\n", err)
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		email := claims["Email"].(string)

		userAuth := &models.UserAuth{Email: email}
		err = userAuth.Get()
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println("User not found\n", err)
			return
		}

		if err = userAuth.ChangePassword(data["password"]); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("Error changing password\n", err)
			return
		}

		w.WriteHeader(http.StatusOK)
		log.Println("User ", userAuth.UID, " changed password")
		return
	}

	w.WriteHeader(http.StatusUnauthorized)
}
