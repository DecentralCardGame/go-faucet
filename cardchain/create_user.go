package cardchain

import (
	"context"

	"github.com/DecentralCardGame/cardchain/x/cardchain/types"
	"github.com/DecentralCardGame/go-faucet/cardchain/client"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/ignite/cli/v29/ignite/pkg/cosmosclient"
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

	msg := types.NewMsgUserCreate(
		address,
		userAddr.String(),
		alias,
	)

	account, err := cosmos.Account(creator)
	if err != nil {
		return cosmosclient.Response{}, err
	}

	return cosmos.BroadcastTx(context.Background(), account, msg)
}
