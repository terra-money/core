package post

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/terra-money/core/v2/x/feeburn/types"
)

// Define the expected keeper interface
type BankKeeper interface {
	BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
}

type FeeBurnKeeper interface {
	GetParams(ctx sdk.Context) types.Params
}
