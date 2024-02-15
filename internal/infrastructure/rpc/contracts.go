package RPCServer

import "github.com/google/uuid"

type ProjectUsecase interface {
	DetermineStrategy(userUuid string, role string) error
	CreateDefaultProject(userUuid uuid.UUID) error
	GetProjectByToken(token string) (*ProjectRPCPayload, error)
	GetProjectListByUserUuid(userUuid uuid.UUID, filterUserUuid uuid.UUID) ([]ProjectRPCPayload, error)
}
