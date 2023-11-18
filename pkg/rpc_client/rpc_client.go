package RPCClient

import (
	"errors"
	"net/rpc"
	"sync"
)

type ReconnectRPCClient struct {
	mutex   sync.Mutex
	rpc     *rpc.Client
	network string
	address string
}

func NewClient(protocol, address string) *ReconnectRPCClient {
	return &ReconnectRPCClient{
		network: protocol,
		address: address,
	}
}

// Close closes the underlying socket file descriptor.
func (rpcClient *ReconnectRPCClient) Close() error {
	rpcClient.mutex.Lock()
	defer rpcClient.mutex.Unlock()
	// If rpc client has not connected yet there is nothing to close.
	if rpcClient.rpc == nil {
		return nil
	}
	// Reset rpcClient.rpc to allow for subsequent calls to use a new
	// (socket) connection.
	client := rpcClient.rpc
	rpcClient.rpc = nil
	return client.Close()
}

// Call makes RPC call to the remote endpoint using the default codec, namely encoding/gob.
func (rpcClient *ReconnectRPCClient) Call(serviceMethod string, args interface{}, reply interface{}) (err error) {
	rpcClient.mutex.Lock()
	defer rpcClient.mutex.Unlock()
	dialCall := func() error {
		// If the rpc.Client is nil, we attempt to (re)connect with the remote endpoint.
		if rpcClient.rpc == nil {
			client, err := rpc.Dial(rpcClient.network, rpcClient.address)
			if err != nil {
				return err
			}
			rpcClient.rpc = client
		}
		// If the RPC fails due to a network-related error, then we reset
		// rpc.Client for a subsequent reconnect.
		return rpcClient.rpc.Call(serviceMethod, args, reply)
	}
	if err = dialCall(); errors.Is(err, rpc.ErrShutdown) {
		rpcClient.rpc.Close()
		rpcClient.rpc = nil
		err = dialCall()
	}
	return err
}
