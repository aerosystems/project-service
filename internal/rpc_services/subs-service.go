package RPCServices

import (
	"github.com/google/uuid"
	"net/rpc"
)

type SubscriptionService interface {
	GetSubscriptionKind(userUuid uuid.UUID) (string, error)
}

type SubscriptionRPC struct {
	rpcClient *rpc.Client
}

func NewSubsRPC(rpcClient *rpc.Client) *SubscriptionRPC {
	return &SubscriptionRPC{
		rpcClient: rpcClient,
	}
}

type SubscriptionRPCPayload struct {
	UserUuid uuid.UUID
	Kind     string
}

func (sr *SubscriptionRPC) GetSubscriptionKind(userUuid uuid.UUID) (string, error) {
	var resSub SubscriptionRPCPayload
	err := sr.rpcClient.Call("SubsServer.GetSubscription", SubscriptionRPCPayload{
		UserUuid: userUuid,
	}, &resSub)
	if err != nil {
		return "", err
	}
	return resSub.Kind, nil
}
