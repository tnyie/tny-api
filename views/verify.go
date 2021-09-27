package views

import (
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
	"github.com/spf13/viper"
	"github.com/tnyie/tny-api/models"
)

// VerifyEmail checks verification token and enables a user
func VerifyEmailCheck(w http.ResponseWriter, r *http.Request) {
	log.Println("Verifying email")
	tokenString := chi.URLParam(r, "token")
	if tokenString == "" {
		log.Println("Invalid verification token")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(viper.GetString("tny.auth.key")), nil
	})
	if err != nil {
		log.Println("Error parsing verification token\n", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		email := claims["Email"].(string)
		user := &models.UserAuth{Email: email}
		err = user.GetByEmail()
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println("Couldn't get user")
			return
		}
		err = user.Verify()
		if err != nil {
			log.Println("Error enabling user\n", err)
			return
		}
		http.Redirect(w, r, "https://tny.ie", http.StatusTemporaryRedirect)
		log.Println("Enabled user")
	}
}
