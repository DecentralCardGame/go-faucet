package cardchain

import (
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/ignite-hq/cli/ignite/pkg/cosmosclient"
)

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
