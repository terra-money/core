package post

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	customwasmkeeper "github.com/terra-money/core/v2/x/wasm/keeper"
)

type WasmdDecorator struct {
	wasmKeeper customwasmkeeper.Keeper
}

func NewWasmdDecorator(wk customwasmkeeper.Keeper) WasmdDecorator {
	return WasmdDecorator{
		wasmKeeper: wk,
	}
}

// MUST: this should always be the latest decorator executed on tx processing
// as it clears the executed contract addresses
func (fsd WasmdDecorator) PostHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, success bool, next sdk.PostHandler) (newCtx sdk.Context, err error) {
	fsd.wasmKeeper.DeleteExecutedContractAddresses(ctx)
	return next(ctx, tx, simulate, success)
}
