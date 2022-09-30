package payload

import (
	"fmt"
	"net/http"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
)

type Payload struct {
	Token   string
	Address string
	Alias   string
}

func (p Payload) Verify(w http.ResponseWriter) bool {
	_, err := sdktypes.AccAddressFromBech32(p.Address)

	if p.Token == "" {
		http.Error(w, "Field token is empty", http.StatusBadRequest)
		return false
	} else if err != nil {
		http.Error(
			w,
			fmt.Sprintf("Address is invalid: %s", err.Error()),
			http.StatusBadRequest,
		)
		return false
	}
	return true
}
