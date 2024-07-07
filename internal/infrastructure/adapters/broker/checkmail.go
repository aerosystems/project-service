package broker

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"fmt"
	PubSub "github.com/aerosystems/customer-service/pkg/pubsub"
	"time"
)

const (
	defaultTimeout = 2 * time.Second
)

type CheckmailEventsAdapter struct {
	pubsubClient         *PubSub.Client
	topicId              string
	subName              string
	createAccessEndpoint string
}

func NewCheckmailEventsAdapter(pubsubClient *PubSub.Client, topicId, subName, createAccessEndpoint string) *CheckmailEventsAdapter {
	return &CheckmailEventsAdapter{
		pubsubClient:         pubsubClient,
		topicId:              topicId,
		subName:              subName,
		createAccessEndpoint: createAccessEndpoint,
	}
}

type CreateAccessEvent struct {
	Token            string    `json:"token"`
	SubscriptionType string    `json:"subscriptionType"`
	AccessTime       time.Time `json:"accessTime"`
}

func (s CheckmailEventsAdapter) PublishCreateAccessEvent(token, subscriptionType string, accessTime time.Time) error {
	event := CreateAccessEvent{
		Token:            token,
		SubscriptionType: subscriptionType,
		AccessTime:       accessTime,
	}
	eventData, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal create access event: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	topic := s.pubsubClient.Client.Topic(s.topicId)
	ok, err := topic.Exists(ctx)
	defer topic.Stop()
	if err != nil {
		return fmt.Errorf("failed to check if topic exists: %w", err)
	}
	if !ok {
		if _, err := s.pubsubClient.Client.CreateTopic(ctx, s.topicId); err != nil {
			return fmt.Errorf("failed to create topic: %w", err)
		}
	}

	sub := s.pubsubClient.Client.Subscription(s.subName)
	ok, err = sub.Exists(ctx)
	if err != nil {
		return fmt.Errorf("failed to check if subscription exists: %w", err)
	}
	if !ok {
		if _, err := s.pubsubClient.Client.CreateSubscription(ctx, s.subName, pubsub.SubscriptionConfig{
			Topic:       topic,
			AckDeadline: 10 * time.Second,
			PushConfig: pubsub.PushConfig{
				Endpoint: s.createAccessEndpoint,
			},
		}); err != nil {
			return fmt.Errorf("failed to make API access: %w", err)
		}
	}

	result := topic.Publish(ctx, &pubsub.Message{
		Data: eventData,
	})
	if _, err := result.Get(ctx); err != nil {
		return fmt.Errorf("failed to publish API access event: %w", err)
	}

	return nil
}
