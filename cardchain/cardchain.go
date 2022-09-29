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
	return cosmosclient.New(
		context.Background(),
		cosmosclient.WithAddressPrefix("cc"),
		cosmosclient.WithNodeAddress(os.Getenv("RPC_NODE")),
	)
}

func getAddr(cosmos cosmosclient.Client, user string) (sdktypes.AccAddress, error) {
	address, err := cosmos.Address(user)
	if err != nil {
		return nil, err
	}
	return address, nil
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
