package cardchain

import (
	"context"
	"os"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/ignite-hq/cli/ignite/pkg/cosmosclient"
)

func SetConfig() {
	config := sdktypes.GetConfig()
	config.SetBech32PrefixForAccount("cc", "ccpub")
}

func getClient() (cosmosclient.Client, error) {
	chainHome := os.Getenv("CHAIN_HOME")
	if chainHome == "" {
		chainHome = "~/.Cardchain"
	}

	return cosmosclient.New(
		context.Background(),
		cosmosclient.WithHome(chainHome),
		cosmosclient.WithAddressPrefix("cc"),
		cosmosclient.WithNodeAddress(os.Getenv("RPC_NODE")),
	)
}

func broadcastMsg(
	cosmos cosmosclient.Client,
	creator string,
	msg sdktypes.Msg,
) (*cosmosclient.Response, error) {
	resp, err := cosmos.BroadcastTx(creator, msg)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
