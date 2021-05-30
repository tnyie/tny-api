package oidc

import (
	"context"
	"log"
	"net/http"

	oidc "github.com/coreos/go-oidc"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

var (
	// TODO change so state is per user request
	state        = viper.GetString("oauth.state")
	provider     *oidc.Provider
	oauth2Config *oauth2.Config
	verifier     *oidc.IDTokenVerifier
)

func InitOIDC() {
	var err error
	provider, err = oidc.NewProvider(context.Background(), "http://auth.tny.local:8080/auth/realms/netsoc-cloud")
	if err != nil {
		panic(err)
	}
	log.Println("bruh")
	scopes := viper.GetStringSlice("oauth.scopes")
	scopes = append(scopes, "profile", "email")

	oauth2Config = &oauth2.Config{
		ClientID:     viper.GetString("oauth.client.id"),
		ClientSecret: viper.GetString("oauth.client.secret"),
		RedirectURL:  viper.GetString("oauth.redirect.url"),

		// Discovery returns the OAuth2 endpoints.
		Endpoint: provider.Endpoint(),

		// "openid" is a required scope for OpenID Connect flows.
		Scopes: scopes,
	}
	verifier = provider.Verifier(&oidc.Config{ClientID: viper.GetString("oauth.client.id")})
}

func HandleCallback(w http.ResponseWriter, r *http.Request) {
	oauth2Token, err := oauth2Config.Exchange(context.Background(), r.URL.Query().Get("code"))
	if err != nil {
		log.Println("error handling callback from oauth provider\n", err)
		return
	}

	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		log.Println("Couldn't extract id token from oauth provider\n", err)
		return
	}

	idToken, err := verifier.Verify(context.Background(), rawIDToken)
	if err != nil {
		log.Println("Failed to verify token from oauth provider\n", err)
		return
	}

	var claims struct {
		Email string
		Roles []string
	}

	if err := idToken.Claims(&claims); err != nil {
		log.Println("error\n", err)
	}
	w.Write([]byte(rawIDToken))
}

func HandleRedirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, oauth2Config.AuthCodeURL(state), http.StatusFound)
}
