package adapters

import (
	"context"
	"github.com/aerosystems/common-service/clients/grpcclient"
	"github.com/aerosystems/common-service/gen/protobuf/subscription"
	"github.com/aerosystems/project-service/internal/entities"
	"github.com/google/uuid"
	"time"
)

type SubscriptionAdapter struct {
	client subscription.SubscriptionServiceClient
}

func NewSubscriptionAdapter(address string) (*SubscriptionAdapter, error) {
	conn, err := grpcclient.NewGRPCConn(address)
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
