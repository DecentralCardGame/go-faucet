package cardchain

import (
	"github.com/DecentralCardGame/Cardchain/x/cardchain/types"
	"github.com/DecentralCardGame/go-faucet/cardchain/client"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/ignite-hq/cli/ignite/pkg/cosmosclient"
)

func CreateUser(creator string, alias string, userAddressString string) (cosmosclient.Response, error) {
	cosmos := client.Get()

	address, err := cosmos.Address(creator)
	if err != nil {
		return cosmosclient.Response{}, err
	}

	userAddr, err := sdktypes.AccAddressFromBech32(userAddressString)
	if err != nil {
		return cosmosclient.Response{}, err
	}

	msg := types.NewMsgCreateuser(
		address.String(),
		userAddr.String(),
		alias,
	)

	return cosmos.BroadcastTx(creator, msg)
}
