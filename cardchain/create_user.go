package cardchain

import (
	"github.com/DecentralCardGame/Cardchain/x/cardchain/types"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
)

func CreateUser(creator string, alias string, userAddressString string) error {
	cosmos, err := getClient()
	if err != nil {
		return err
	}

	address, err := getAddr(cosmos, creator)
	if err != nil {
		return err
	}

	userAddr, err := sdktypes.AccAddressFromBech32(userAddressString)
	if err != nil {
		return err
	}

	msg := types.NewMsgCreateuser(
		address.String(),
		userAddr.String(),
		alias,
	)

	err = broadcastMsg(cosmos, creator, msg)

	return err
}
