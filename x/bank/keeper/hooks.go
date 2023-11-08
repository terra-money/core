package keeper

import (
	customterratypes "github.com/terra-money/core/v2/x/bank/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Implements StakingHooks interface
var _ customterratypes.BankHooks = Keeper{}

// TrackBeforeSend executes the TrackBeforeSend hook if registered.
func (k Keeper) TrackBeforeSend(ctx sdk.Context, from, to sdk.AccAddress, amount sdk.Coins) {
	if k.hooks != nil {
		k.hooks.TrackBeforeSend(ctx, from, to, amount)
	}
}

// BlockBeforeSend executes the BlockBeforeSend hook if registered.
func (k Keeper) BlockBeforeSend(ctx sdk.Context, from, to sdk.AccAddress, amount sdk.Coins) error {
	if k.hooks != nil {
		return k.hooks.BlockBeforeSend(ctx, from, to, amount)
	}
	return nil
}

// MultiBankHooks combine multiple bank hooks, all hook functions are run in array sequence
type MultiBankHooks []customterratypes.BankHooks

// NewMultiBankHooks takes a list of BankHooks and returns a MultiBankHooks
func NewMultiBankHooks(hooks ...customterratypes.BankHooks) MultiBankHooks {
	return hooks
}

// TrackBeforeSend runs the TrackBeforeSend hooks in order for each BankHook in a MultiBankHooks struct
func (h MultiBankHooks) TrackBeforeSend(ctx sdk.Context, from, to sdk.AccAddress, amount sdk.Coins) {
	for i := range h {
		h[i].TrackBeforeSend(ctx, from, to, amount)
	}
}

// BlockBeforeSend runs the BlockBeforeSend hooks in order for each BankHook in a MultiBankHooks struct
func (h MultiBankHooks) BlockBeforeSend(ctx sdk.Context, from, to sdk.AccAddress, amount sdk.Coins) error {
	for i := range h {
		err := h[i].BlockBeforeSend(ctx, from, to, amount)
		if err != nil {
			return err
		}
	}
	return nil
}
