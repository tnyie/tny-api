package authentication

import (
	"context"
	"log"
	"path/filepath"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

var authClient *auth.Client
var client *firebase.App

// InitFirebase creates an app client and authentication client
func InitFirebase() {
	serviceAccountKeyFilePath, err := filepath.Abs("/config/serviceAccountKey.json")
	if err != nil {
		log.Fatal("Unable to load firebase config\n", err)
	}

	opt := option.WithCredentialsFile(serviceAccountKeyFilePath)

	// Firebase admin SDK init
	client, err = firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatal("Firebase load error\n", err)
	}
	authClient, err = client.Auth(context.Background())
	if err != nil {
		log.Fatal("Firebase authentication load error\n", err)
	}
}
