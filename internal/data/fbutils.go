package data

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

type FirebaseUtils struct {
	App  *firebase.App
	Auth *auth.Client
}

func NewFirebaseUtil() (*FirebaseUtils, error) {
	ctx := context.Background()
	opt := []option.ClientOption{option.WithCredentialsFile("./serviceAccountKey.json")}
	app, err := firebase.NewApp(context.Background(), nil, opt...)
	if err != nil {
		return nil, fmt.Errorf("error initializing app: %v", err)
	}

	authClient, err := app.Auth(ctx)
	if err != nil {
		return nil, fmt.Errorf("error initializing auth: %v", err)
	}

	return &FirebaseUtils{
		App:  app,
		Auth: authClient,
	}, err
}
