package RPCServer

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
)

func Listen(rpcPort string) error {
	log.Printf("starting RPC server on 0.0.0.0:%s", rpcPort)
	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", rpcPort))
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
