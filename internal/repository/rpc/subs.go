package RPCServices

import (
	"github.com/aerosystems/project-service/internal/models"
	"github.com/aerosystems/project-service/pkg/rpc_client"
	"github.com/google/uuid"
	"time"
)

type SubscriptionRPC struct {
	rpcClient *RPCClient.ReconnectRPCClient
}

func NewSubsRPC(rpcClient *RPCClient.ReconnectRPCClient) *SubscriptionRPC {
	return &SubscriptionRPC{
		rpcClient: rpcClient,
	}
}

type SubsRPCPayload struct {
	UserUuid   uuid.UUID
	Kind       models.KindSubscription
	AccessTime time.Time
}

func (sr *SubscriptionRPC) GetSubscription(userUuid uuid.UUID) (models.KindSubscription, time.Time, error) {
	var resSub SubsRPCPayload
	err := sr.rpcClient.Call("SubsServer.GetSubscription", userUuid, &resSub)
	if err != nil {
		return "", time.Time{}, err
	}
	return resSub.Kind, resSub.AccessTime, nil
}
