package data

import (
	"context"
	"database/sql"
	"fmt"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

type Firebase struct {
	App  *firebase.App
	Auth *auth.Client
}

type FcmToken struct {
	UserID   string `json:"userID" validate:"required"`
	FcmToken string `json:"fcmtoken" validate:"required"`
}

type FcmModel struct {
	DB *sql.DB
}

func NewFirebaseIntegration() (*Firebase, error) {
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

	return &Firebase{
		App:  app,
		Auth: authClient,
	}, err
}

func (f FcmModel) Create(fcm *FcmToken) error {
	query := `
	INSERT INTO fcm_tokens (user_id, token)
	VALUES($1, $2)
	RETURNING id
	`

	ctx, cancel := handlectx()
	defer cancel()

	err := f.DB.QueryRowContext(ctx, query, fcm.UserID, fcm.FcmToken).Scan(&fcm.UserID)
	if err != nil {
		return err
	}

	return nil
}

func (f FcmModel) Get(userID string) ([]string, error) {
	query := `SELECT token FROM fcm_tokens WHERE user_id=$1`

	ctx, cancel := handlectx()
	defer cancel()

	rows, err := f.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tokens []string
	for rows.Next() {
		var token string

		err := rows.Scan(&token)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}

	return tokens, nil
}

func (f FcmModel) Delete(userID string, token string) error {
	query := `DELETE FROM fcm_tokens WHERE user_id=$1 AND token=$2`

	ctx, cancel := handlectx()
	defer cancel()

	_, err := f.DB.ExecContext(ctx, query, userID, token)
	if err != nil {
		return err
	}

	return nil
}
