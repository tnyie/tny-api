package middleware

import (
	"context"
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
)

type AuthCtx struct{}

func CheckToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// fetch authorization heade
		parts := strings.Split(r.Header.Get("Authorization"), " ")
		if parts[0] != "Bearer" {
			next.ServeHTTP(w, r)
			return
		}
		if parts[1] != "" {
			token, err := jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
				return []byte(viper.GetString("tny.auth.key")), nil
			})
			if err != nil {
				log.Println("Error parsing token")
				next.ServeHTTP(w, r)
				return
			}

			log.Println(token.Valid)
			log.Println(reflect.TypeOf(token.Claims))

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				log.Println("Claims are valid")
				log.Println("Valid token, claims: ", claims)
				next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), AuthCtx{}, claims)))
				return
			}
			log.Println("Invalid claims")
			next.ServeHTTP(w, r)
		}
		log.Println("No authentication header")
		// TODO
		next.ServeHTTP(w, r)
	})
}
