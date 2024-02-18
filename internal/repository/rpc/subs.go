package RpcRepo

import (
	"github.com/aerosystems/project-service/internal/models"
	"github.com/aerosystems/project-service/pkg/rpc_client"
	"github.com/google/uuid"
	"time"
)

type SubsRepo struct {
	rpcClient *RPCClient.ReconnectRPCClient
}

func NewSubsRepo(rpcClient *RPCClient.ReconnectRPCClient) *SubsRepo {
	return &SubsRepo{
		rpcClient: rpcClient,
	}
}

type SubsRPCPayload struct {
	UserUuid   uuid.UUID
	Kind       models.KindSubscription
	AccessTime time.Time
}

func (sr *SubsRepo) GetSubscription(userUuid uuid.UUID) (models.KindSubscription, time.Time, error) {
	var resSub SubsRPCPayload
	err := sr.rpcClient.Call("Server.GetSubscription", userUuid, &resSub)
	if err != nil {
		return "", time.Time{}, err
	}
	return resSub.Kind, resSub.AccessTime, nil
}
