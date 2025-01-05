package GRPCServer

import (
	"context"
	"github.com/aerosystems/project-service/internal/common/protobuf/project"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ProjectHandler struct {
	projectUsecase ProjectUsecase
	project.UnimplementedProjectServiceServer
}

func NewProjectHandler(projectUsecase ProjectUsecase) *ProjectHandler {
	return &ProjectHandler{
		projectUsecase: projectUsecase,
	}
}

func (h ProjectHandler) CreateDefaultProject(_ context.Context, req *project.CreateDefaultProjectRequest) (*project.CreateDefaultProjectResponse, error) {
	defaultProject, err := h.projectUsecase.CreateDefaultProject(uuid.MustParse(req.CustomerUuid))
	if err != nil {
		return nil, err
	}
	return &project.CreateDefaultProjectResponse{
		ProjectUuid: defaultProject.Uuid.String(),
	}, nil
}

func (h ProjectHandler) DeleteProject(_ context.Context, req *project.DeleteProjectRequest) (*emptypb.Empty, error) {
	err := h.projectUsecase.DeleteProject(req.ProjectUuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
