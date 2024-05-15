package RpcRepo

import (
	"github.com/aerosystems/project-service/internal/models"
	"github.com/aerosystems/project-service/pkg/rpc_client"
	"github.com/google/uuid"
	"time"
)

type SubsRepo struct {
	rpcClient *RpcClient.ReconnectRpcClient
}

func NewSubsRepo(rpcClient *RpcClient.ReconnectRpcClient) *SubsRepo {
	return &SubsRepo{
		rpcClient: rpcClient,
	}
}

type SubsRPCPayload struct {
	UserUuid   uuid.UUID
	Kind       string
	AccessTime time.Time
}

func (sr *SubsRepo) GetSubscription(userUuid uuid.UUID) (models.KindSubscription, time.Time, error) {
	var resSub SubsRPCPayload
	err := sr.rpcClient.Call("Server.GetSubscription", userUuid, &resSub)
	if err != nil {
		return models.KindSubscription{}, time.Time{}, err
	}
	return models.NewKindSubscription(resSub.Kind), resSub.AccessTime, nil
}
