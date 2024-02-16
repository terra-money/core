package keeper

import (
	"github.com/terra-money/core/v2/x/smartaccounts/types"
)

var _ types.MsgServer = &Keeper{}
