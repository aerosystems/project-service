package gcp

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"google.golang.org/api/option"
	"os"
)

type PubSubClient struct {
	Client *pubsub.Client
	Ctx    context.Context
}

type Config struct {
	Type                    string `json:"type"`
	ProjectId               string `json:"project_id"`
	PrivateKeyId            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`
	ClientEmail             string `json:"client_email"`
	ClientId                string `json:"client_id"`
	AuthURI                 string `json:"auth_uri"`
	TokenURI                string `json:"token_uri"`
	AuthProviderX509CertURL string `json:"auth_provider_x509_cert_url"`
	ClientX509CertURL       string `json:"client_x509_cert_url"`
	UniverseDomain          string `json:"universe_domain"`
}

func NewPubSubClient(projectId string) (*PubSubClient, error) {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectId)
	if err != nil {
		return nil, err
	}

	return &PubSubClient{
		Client: client,
		Ctx:    ctx,
	}, nil

}

func NewPubSubClientWithAuth(credentialsPath string) (*PubSubClient, error) {
	configData, err := os.ReadFile(credentialsPath)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(configData, &config)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, config.ProjectId, option.WithCredentialsFile(credentialsPath))
	if err != nil {
		return nil, err
	}

	return &PubSubClient{
		Client: client,
		Ctx:    ctx,
	}, nil
}
