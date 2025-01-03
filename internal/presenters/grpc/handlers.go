package GRPCServer

import (
	"context"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Handler struct {
	projectUsecase ProjectUsecase
}

func NewGRPCHandler(projectUsecase ProjectUsecase) *Handler {
	return &Handler{
		projectUsecase: projectUsecase,
	}
}

func (h Handler) CreateDefaultProject(_ context.Context, req *CreateDefaultProjectRequest) (*CreateDefaultProjectResponse, error) {
	if err := h.projectUsecase.DetermineStrategy(req.CustomerUuid, "customer"); err != nil {
		return nil, err
	}
	project, err := h.projectUsecase.CreateDefaultProject(uuid.MustParse(req.CustomerUuid))
	if err != nil {
		return nil, err
	}
	return &CreateDefaultProjectResponse{
		ProjectUuid: project.Uuid.String(),
	}, nil
}
func (h Handler) DeleteProject(_ context.Context, req *DeleteProjectRequest) (*emptypb.Empty, error) {
	err := h.projectUsecase.DeleteProject(req.ProjectUuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (h Handler) mustEmbedUnimplementedProjectServiceServer() {}
