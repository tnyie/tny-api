package util

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/tnyie/tny-api/middleware"
	"github.com/tnyie/tny-api/models"
)

// CheckLogin checks a request for authentication context,
// and compares the uid with the required one
func CheckLogin(r *http.Request, requiredID string) (*models.UserAuth, bool, bool) {

	user := getUserAuth(r)

	if user == nil && user.Enabled && user.UID != "" {
		if requiredID == "" {
			return user, true, false
		} else if user.UID == requiredID {
			return user, true, false
		} else if user.Admin {
			return user, false, true
		}
	}

	return nil, false, false
}

func getUserAuth(r *http.Request) *models.UserAuth {
	var user *models.UserAuth
	if claims, ok := r.Context().Value(middleware.BearerCtx{}).(jwt.MapClaims); ok {
		user.UID = claims["UserID"].(string)
	} else if userID, ok := r.Context().Value(middleware.KeyCtx{}).(string); ok {
		user.UID = userID
	} else {
		return nil
	}

	user.Get()
	return user
}

// func IsAdmin(r *http.Request) bool {
// 	if claims, ok := r.Context().Value(middleware.BearerCtx{}).(jwt.MapClaims); ok {
// 		if claims["UserID"] != "" {
// 			userAuth := &models.UserAuth{UID: claims["UserID"].(string)}
// 			userAuth.Get()
// 			if userAuth.Enabled && userAuth.Admin {
// 				return true
// 			} else {
// 				return false
// 			}
// 		}
// 		return false
// 	} else if userID, ok := r.Context().Value(middleware.KeyCtx{}).(string); ok {
// 		userAuth := &models.UserAuth{UID: userID}
// 		userAuth.Get()
// 		if userAuth.Enabled && userAuth.Admin {
// 			return true
// 		} else {
// 			return false
// 		}
// 	}
// 	return false
// }
