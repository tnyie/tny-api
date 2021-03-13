package util

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/tnyie/tny-api/middleware"
	"github.com/tnyie/tny-api/models"
)

// CheckLogin checks a request for authentication context, and compares the uid with the required one
func CheckLogin(r *http.Request, requiredID string) (*models.UserAuth, bool) {
	if claims, ok := r.Context().Value(middleware.AuthCtx{}).(jwt.MapClaims); ok {
		if claims["UserID"] == requiredID {
			// check if user is enabled
			user := &models.UserAuth{UID: claims["UserID"].(string)}
			user.Get()
			if user.Enabled == true {
				return user, true
			}
		}
	}
	return nil, false
}
