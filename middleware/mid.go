package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gal/tny/authentication"
)

type authCtx interface{}

// AuthCtx context key to access authentication response
var AuthCtx authCtx

// FirebaseAuth middleware passes user's id in context
func FirebaseAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tokenString := r.Header.Get("Authorization")
		idToken := strings.TrimSpace(strings.Replace(tokenString, "Bearer", "", 1))
		if idToken == "" {
			// if they are not authenticated, continue without auth, but chekc in handler
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), AuthCtx, "")))
			return
		}
		token, err := authentication.InspectToken(idToken)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), AuthCtx, token.UID)))
		return
	})
}
