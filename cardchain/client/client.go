package client

import (
	"context"
	"log"
	"sync"
	"time"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/ignite/cli/ignite/pkg/cosmosclient"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
)

var client *cosmosclient.Client
var once sync.Once

func setConfig() {
	config := sdktypes.GetConfig()
	config.SetBech32PrefixForAccount("cc", "ccpub")
}

func WaitForChain(config Config) error {
	rpc, err := rpchttp.New(config.RPCNode, "/websocket")
	if err != nil {
		return err
	}

	for {
		resp, err := rpc.Status(context.Background())
		if err == nil {
			log.Printf("Found chain with id: %s", resp.NodeInfo.Network)
			break
		}

		log.Print("Waiting for blockchain...")
		time.Sleep(time.Second)
	}

	return Init(config)
}

func Init(config Config) error {
	setConfig()

	localClient, err := cosmosclient.New(
		context.Background(),
		cosmosclient.WithHome(config.ChainHome),
		cosmosclient.WithAddressPrefix("cc"),
		cosmosclient.WithGas("600000"),
		cosmosclient.WithNodeAddress(config.RPCNode),
	)

	once.Do(func() {
		client = &localClient
	})

	return err
}

func Get() cosmosclient.Client {
	return *client
}
