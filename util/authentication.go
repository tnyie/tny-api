package util

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/tnyie/tny-api/middleware"
)

// CheckLogin checks a request for authentication context, and compares the uid with the required one
func CheckLogin(r *http.Request, requiredID string) bool {
	if claims, ok := r.Context().Value(middleware.AuthCtx{}).(jwt.MapClaims); ok {
		if claims["UserID"] == requiredID {
			return true
		}
	}
	return false
}
