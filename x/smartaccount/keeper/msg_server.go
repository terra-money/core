package keeper

import (
	"github.com/terra-money/core/v2/x/smartaccount/types"
)

var _ types.MsgServer = &Keeper{}
