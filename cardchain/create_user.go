package cardchain

import (
	"github.com/DecentralCardGame/Cardchain/x/cardchain/types"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/ignite-hq/cli/ignite/pkg/cosmosclient"
)

func CreateUser(creator string, alias string, userAddressString string) (*cosmosclient.Response, error) {
	cosmos, err := getClient()
	if err != nil {
		return nil, err
	}

	address, err := getAddr(cosmos, creator)
	if err != nil {
		return nil, err
	}

	userAddr, err := sdktypes.AccAddressFromBech32(userAddressString)
	if err != nil {
		return nil, err
	}

	msg := types.NewMsgCreateuser(
		address.String(),
		userAddr.String(),
		alias,
	)

	return broadcastMsg(cosmos, creator, msg)
}
