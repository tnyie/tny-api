package util

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/tnyie/tny-api/middleware"
	"github.com/tnyie/tny-api/models"
)

// CheckLogin checks a request for authentication context,
// and compares the uid with the required one
func CheckLogin(r *http.Request, requiredID string) (*models.UserAuth, bool) {

	if claims, ok := r.Context().Value(middleware.BearerCtx{}).(jwt.MapClaims); ok {
		// if claims["UserID"] == requiredID || requiredID == "" {
		// check if user is enabled
		user := &models.UserAuth{UID: claims["UserID"].(string)}
		user.Get()
		if user.Enabled && user.UID != "" {
			if requiredID == "" {
				return user, true
			} else {
				if user.UID == requiredID {
					return user, true
				} else {
					// allow an admin to perform any operation if UID doesn't match
					if user.Admin {
						return nil, true
					}
				}
			}
		}

		return user, false
	} else if userID, ok := r.Context().Value(middleware.KeyCtx{}).(string); ok {
		// if using APIKey
		user := &models.UserAuth{UID: userID}
		user.Get()
		if user.Enabled && user.UID != "" {
			// return user, true
			if requiredID == "" {
				return user, true
			} else {
				if user.UID == requiredID {
					return user, true
				} else {
					// allow an admin to perform any user operation
					if user.Admin {
						return nil, true
					}
				}
			}
		}
	}
	return nil, false
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
