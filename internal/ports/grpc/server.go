package GRPCServer

import (
	"github.com/aerosystems/common-service/gen/protobuf/project"
	"github.com/aerosystems/common-service/presenters/grpcserver"
	"github.com/sirupsen/logrus"
)

type Server struct {
	grpcServer *grpcserver.Server
}

func NewGRPCServer(
	cfg *grpcserver.Config,
	log *logrus.Logger,
	projectService *ProjectService,
) *Server {
	server := grpcserver.NewGRPCServer(cfg, log)

	server.RegisterService(project.ProjectService_ServiceDesc, projectService)

	return &Server{
		grpcServer: server,
	}
}

func (s *Server) Run() error {
	return s.grpcServer.Run()
}
