package RPCClient

import (
	"log"
	"net/rpc"
	"time"
)

func NewClient(network, address string) *rpc.Client {
	count := 0

	for {
		client, err := rpc.Dial(network, address)
		if err != nil {
			log.Printf("RPC Server %s not ready...", address)
			count++
		} else {
			return client
		}

		if count > 10 {
			log.Println(err)
			return nil
		}

		log.Println("Backing off for two seconds...")
		time.Sleep(2 * time.Second)
		continue
	}
}
