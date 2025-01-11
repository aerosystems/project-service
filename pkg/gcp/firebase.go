package gcp

import (
	"context"
	"errors"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"fmt"
	"google.golang.org/api/option"
)

type FirebaseApp struct {
	Client *auth.Client
}

func NewFirebaseApp(projectId string, serviceAccountFilePath string) (*FirebaseApp, error) {
	var opts []option.ClientOption
	if file := serviceAccountFilePath; file != "" {
		opts = append(opts, option.WithCredentialsFile(file))
	}

	config := &firebase.Config{ProjectID: projectId}
	firebaseApp, err := firebase.NewApp(context.Background(), config, opts...)
	if err != nil {
		return nil, fmt.Errorf("error initializing app: %w\n", err)
	}

	authClient, err := firebaseApp.Auth(context.Background())
	if err != nil {
		return nil, errors.New("unable to create firebase Auth client")
	}
	return &FirebaseApp{
		Client: authClient,
	}, nil
}
