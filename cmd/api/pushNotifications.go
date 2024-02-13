package main

import (
	"context"

	"firebase.google.com/go/messaging"
)

func (app *application) NotifyOne(ctx context.Context, title string, body string, token string) (string, error) {
	response, err := app.fcmClient.Send(ctx, &messaging.Message{
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Token: token,
	})

	if err != nil {
		return "", err
	}

	return response, nil
}

func (app *application) NotifyMany(ctx context.Context, title string, body string, tokens []string) (*messaging.BatchResponse, error) {
	response, err := app.fcmClient.SendMulticast(ctx, &messaging.MulticastMessage{
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Tokens: tokens,
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}
