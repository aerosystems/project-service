package GRPCServer

import (
	"context"
	"github.com/aerosystems/common-service/gen/protobuf/project"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ProjectService struct {
	projectUsecase ProjectUsecase
	project.UnimplementedProjectServiceServer
}

func NewProjectService(projectUsecase ProjectUsecase) *ProjectService {
	return &ProjectService{
		projectUsecase: projectUsecase,
	}
}

func (ps ProjectService) CreateDefaultProject(ctx context.Context, req *project.CreateDefaultProjectRequest) (*project.CreateDefaultProjectResponse, error) {
	defaultProject, err := ps.projectUsecase.CreateDefaultProject(ctx, uuid.MustParse(req.CustomerUuid))
	if err != nil {
		return nil, err
	}
	return &project.CreateDefaultProjectResponse{
		ProjectUuid:  defaultProject.Uuid.String(),
		ProjectToken: defaultProject.Token,
	}, nil
}

func (ps ProjectService) DeleteProject(ctx context.Context, req *project.DeleteProjectRequest) (*emptypb.Empty, error) {
	err := ps.projectUsecase.DeleteProject(ctx, req.ProjectUuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
