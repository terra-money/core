package post

import (
	feeburnpost "github.com/terra-money/core/v2/x/feeburn/post"
	feesharepost "github.com/terra-money/core/v2/x/feeshare/post"
	customwasmkeeper "github.com/terra-money/core/v2/x/wasm/keeper"
	wasmpost "github.com/terra-money/core/v2/x/wasm/post"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type HandlerOptions struct {
	FeeShareKeeper feesharepost.FeeShareKeeper
	BankKeeper     feesharepost.BankKeeper
	WasmKeeper     customwasmkeeper.Keeper
	FeeBurnKeeper  feeburnpost.FeeBurnKeeper
}

func NewPostHandler(options HandlerOptions) sdk.PostHandler {

	// NOTE: feesharepost handler MUST always run before the wasmpost handler because
	// feeshare will distribute the fees between the contracts enabled with feeshare
	// and then wasmpost will clean the list of executed contracts in the block.
	postDecorators := []sdk.PostDecorator{
		feesharepost.NewFeeSharePayoutDecorator(options.FeeShareKeeper, options.BankKeeper, options.WasmKeeper),
		wasmpost.NewWasmdDecorator(options.WasmKeeper),
		feeburnpost.NewFeeBurnDecorator(options.FeeBurnKeeper, options.BankKeeper),
	}

	return sdk.ChainPostDecorators(postDecorators...)
}
