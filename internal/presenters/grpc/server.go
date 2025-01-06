package GRPCServer

import (
	"fmt"
	"github.com/aerosystems/project-service/internal/common/protobuf/project"
	"net"

	grpclogrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpcctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type Server struct {
	log        *logrus.Logger
	addr       string
	grpcServer *grpc.Server
}

func NewGRPCServer(
	port int,
	log *logrus.Logger,
	grpcHandlers project.ProjectServiceServer,
) *Server {
	logrusEntry := logrus.NewEntry(log)
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpcctxtags.UnaryServerInterceptor(grpcctxtags.WithFieldExtractor(grpcctxtags.CodeGenRequestFieldExtractor)),
			grpclogrus.UnaryServerInterceptor(logrusEntry),
		),
		grpc.ChainStreamInterceptor(
			grpcctxtags.StreamServerInterceptor(grpcctxtags.WithFieldExtractor(grpcctxtags.CodeGenRequestFieldExtractor)),
			grpclogrus.StreamServerInterceptor(logrusEntry),
		),
	)
	project.RegisterProjectServiceServer(grpcServer, grpcHandlers)
	return &Server{
		log:        log,
		addr:       fmt.Sprintf(":%d", port),
		grpcServer: grpcServer,
	}
}

func (s *Server) Run() error {
	listen, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", s.addr, err)
	}
	s.log.Infof("gRPC server running on %s", s.addr)
	return s.grpcServer.Serve(listen)
}
