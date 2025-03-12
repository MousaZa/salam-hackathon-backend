package db

import (
	"context"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/api/option"
)

type Firestore struct {
	Client *firestore.Client
	Logger *hclog.Logger
}

func NewConnection() (*Firestore, error) {
	l := hclog.Default()
	co := option.WithCredentialsFile("firebase-secret.json")
	ctx := context.Background()

	app, err := firebase.NewApp(ctx, nil, co)
	if err != nil {
		l.Error("Error connecting to Firestore", "error", err)
		return nil, err
	}

	client, err := app.Firestore(ctx)

	if err != nil {
		l.Error("Failed to create Firestore client", "error", err)
		return nil, err
	}

	return &Firestore{
		Client: client,
		Logger: &l,
	}, nil

}
