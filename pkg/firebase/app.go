package firebaseApp

import (
	"context"
	"errors"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"fmt"
	"google.golang.org/api/option"
)

type App struct {
	Client *auth.Client
}

func NewApp(projectId string, serviceAccountFilePath string) (*App, error) {
	var opts []option.ClientOption
	if file := serviceAccountFilePath; file != "" {
		opts = append(opts, option.WithCredentialsFile(file))
	}

	config := &firebase.Config{ProjectID: projectId}
	firebaseApp, err := firebase.NewApp(context.Background(), config, opts...)
	if err != nil {
		fmt.Errorf("error initializing app: %v\n", err)
	}

	authClient, err := firebaseApp.Auth(context.Background())
	if err != nil {
		errors.New("unable to create firebase Auth client")
	}
	return &App{
		Client: authClient,
	}, nil
}
