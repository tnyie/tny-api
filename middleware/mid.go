package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"github.com/tnyie/tny-api/models"
)

type BearerCtx struct{}
type KeyCtx struct{}

func CheckToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// fetch authorization heade
		parts := strings.Split(r.Header.Get("Authorization"), " ")

		if parts[1] != "" {
			if parts[0] == "Bearer" {
				token, err := jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
					return []byte(viper.GetString("tny.auth.key")), nil
				})
				if err != nil {
					log.Println("Error parsing token")
					next.ServeHTTP(w, r)
					return
				}

				if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
					log.Println("Claims are valid")
					log.Println("Valid token, claims: ", claims)
					next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), BearerCtx{}, claims)))
					return
				}
				log.Println("Invalid claims")
				next.ServeHTTP(w, r)
			} else if parts[0] == "Key" {
				keyString := parts[1]
				userID, err := models.ValidAPIKey(keyString)
				if err != nil {
					w.WriteHeader(http.StatusUnauthorized)
					log.Println("API Key invalid\n", err)
					next.ServeHTTP(w, r)
					return
				}
				next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), KeyCtx{}, userID)))
				return
			}
		}
		log.Println("No authentication header")
		// TODO
		next.ServeHTTP(w, r)
	})
}
