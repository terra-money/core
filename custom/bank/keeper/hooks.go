package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
	customterratypes "github.com/terra-money/core/v2/custom/bank/types"
)

// Implements StakingHooks interface
var _ customterratypes.BankHooks = types.BaseSendKeeper{}

// TrackBeforeSend executes the TrackBeforeSend hook if registered.
func (k types.BaseSendKeeper) TrackBeforeSend(ctx sdk.Context, from, to sdk.AccAddress, amount sdk.Coins) {
	if k.hooks != nil {
		k.hooks.TrackBeforeSend(ctx, from, to, amount)
	}
}

// BlockBeforeSend executes the BlockBeforeSend hook if registered.
func (k types.BaseSendKeeper) BlockBeforeSend(ctx sdk.Context, from, to sdk.AccAddress, amount sdk.Coins) error {
	if k.hooks != nil {
		return k.hooks.BlockBeforeSend(ctx, from, to, amount)
	}
	return nil
}
