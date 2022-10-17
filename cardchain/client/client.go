package client

import (
	"context"
	"os"
	"sync"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/ignite-hq/cli/ignite/pkg/cosmosclient"
)

var client *cosmosclient.Client
var once sync.Once

func setConfig() {
	config := sdktypes.GetConfig()
	config.SetBech32PrefixForAccount("cc", "ccpub")
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
