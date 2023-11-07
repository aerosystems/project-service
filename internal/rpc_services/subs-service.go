package RPCServices

import (
	"net/rpc"
)

type SubscriptionService interface {
	GetSubscriptionKind(userId int) (string, error)
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
	UserId int
	Kind   string
}

func (sr *SubscriptionRPC) GetSubscriptionKind(userId int) (string, error) {
	var resSub SubscriptionRPCPayload
	err := sr.rpcClient.Call("SubsServer.GetSubscription", SubscriptionRPCPayload{
		UserId: userId,
	}, &resSub)
	if err != nil {
		return "", err
	}
	return resSub.Kind, nil
}
