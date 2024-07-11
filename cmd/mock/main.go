package main

import (
	"encoding/json"
	"fmt"
	"github.com/aerosystems/project-service/internal/models"
	"github.com/google/uuid"
	"time"
)

type EventBody struct {
	Message struct {
		Data []byte `json:"data"`
	} `json:"message"`
	Subscription string `json:"subscription"`
}

type CreateProjectEvent struct {
	CustomerUuid     string    `json:"customerUuid"`
	SubscriptionType string    `json:"subscriptionType"`
	AccessTime       time.Time `json:"accessTime"`
}

func main() {
	customer := CreateProjectEvent{
		CustomerUuid:     uuid.New().String(),
		SubscriptionType: models.TrialSubscription.String(),
		AccessTime:       time.Now().Add(time.Hour * 24 * 7),
	}

	data, err := json.Marshal(customer)
	if err != nil {
		panic(err)
	}

	var event EventBody
	event.Message.Data = data
	event.Subscription = "create-project"

	var jsonData []byte
	jsonData, err = json.Marshal(event)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(jsonData))
}
