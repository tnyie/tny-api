package authentication

import (
	"context"

	"firebase.google.com/go/auth"
)

// InspectToken returns the token object and/or error
func InspectToken(idToken string) (*auth.Token, error) {
	return authClient.VerifyIDToken(context.Background(), idToken)
}
