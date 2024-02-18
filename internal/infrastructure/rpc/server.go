package RpcServer

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net"
	"net/rpc"
)

const rpcPort = 5001

type Server struct {
	log            *logrus.Logger
	projectUsecase ProjectUsecase
}

func NewServer(
	log *logrus.Logger,
	projectUsecase ProjectUsecase,
) *Server {
	return &Server{
		log:            log,
		projectUsecase: projectUsecase,
	}
}

func (s Server) Run() error {
	if err := rpc.Register(s); err != nil {
		return err
	}
	s.log.Infof("starting project-service RPC server on port %d\n", rpcPort)
	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", rpcPort))
	if err != nil {
		return err
	}
	defer listen.Close()

	for {
		rpcConn, err := listen.Accept()
		if err != nil {
			continue
		}
		go rpc.ServeConn(rpcConn)
	}
}
