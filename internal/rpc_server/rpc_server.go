package RPCServer

import (
	"fmt"
	"github.com/aerosystems/project-service/internal/services"
	"github.com/sirupsen/logrus"
	"net"
	"net/rpc"
)

type ProjectServer struct {
	rpcPort        int
	log            *logrus.Logger
	projectService services.ProjectService
}

func NewProjectServer(
	rpcPort int,
	log *logrus.Logger,
	projectService services.ProjectService,
) *ProjectServer {
	return &ProjectServer{
		rpcPort:        rpcPort,
		log:            log,
		projectService: projectService,
	}
}

func (ps *ProjectServer) Listen(rpcPort int) error {
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
