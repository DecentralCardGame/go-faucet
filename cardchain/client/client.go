package client

import (
	"context"
	"log"
	"os"
	"sync"
	"time"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/ignite-hq/cli/ignite/pkg/cosmosclient"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
)

var client *cosmosclient.Client
var once sync.Once

func setConfig() {
	config := sdktypes.GetConfig()
	config.SetBech32PrefixForAccount("cc", "ccpub")
}

func WaitForChain() error {
	rpc, err := rpchttp.New(os.Getenv("RPC_NODE"), "/websocket")
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

	return Init()
}

func Init() error {
	setConfig()

	chainHome := os.Getenv("CHAIN_HOME")
	if chainHome == "" {
		chainHome = "~/.Cardchain"
	}

	localClient, err := cosmosclient.New(
		context.Background(),
		cosmosclient.WithHome(chainHome),
		cosmosclient.WithAddressPrefix("cc"),
		cosmosclient.WithNodeAddress(os.Getenv("RPC_NODE")),
	)

	once.Do(func() {
		client = &localClient
	})

	return err
}

func Get() cosmosclient.Client {
	return *client
}
