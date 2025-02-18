package adapters

import (
	"context"
	"crypto/tls"
	"github.com/aerosystems/common-service/gen/protobuf/subscription"
	"github.com/aerosystems/project-service/internal/entities"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"time"
)

type SubscriptionAdapter struct {
	client subscription.SubscriptionServiceClient
}

func NewSubscriptionAdapter(address string) (*SubscriptionAdapter, error) {
	opts := []grpc.DialOption{
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:    30,
			Timeout: 30,
		}),
	}
	if address[len(address)-4:] == ":443" {
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	conn, err := grpc.NewClient(address, opts...)
	if err != nil {
		return nil, err
	}
	return &SubscriptionAdapter{
		client: subscription.NewSubscriptionServiceClient(conn),
	}, nil
}

func (sa SubscriptionAdapter) GetSubscription(ctx context.Context, customerUUID uuid.UUID) (subscriptionType entities.SubscriptionType, accessTime time.Time, err error) {
	resp, err := sa.client.GetSubscription(ctx, &subscription.GetSubscriptionRequest{
		CustomerUuid: customerUUID.String(),
	})
	if err != nil {
		return entities.UnknownSubscription, time.Time{}, err
	}
	subscriptionType = entities.NewSubscriptionType(resp.SubscriptionType)
	return subscriptionType, resp.AccessTime.AsTime(), nil
}
