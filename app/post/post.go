package post

import (
	feesharepost "github.com/terra-money/core/v2/x/feeshare/post"
	customwasmkeeper "github.com/terra-money/core/v2/x/wasm/keeper"
	wasmpost "github.com/terra-money/core/v2/x/wasm/post"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type HandlerOptions struct {
	FeeShareKeeper feesharepost.FeeShareKeeper
	BankKeeper     feesharepost.BankKeeper
	WasmKeeper     customwasmkeeper.Keeper
}

func NewPostHandler(options HandlerOptions) sdk.PostHandler {

	postDecorators := []sdk.PostDecorator{
		feesharepost.NewFeeSharePayoutDecorator(options.FeeShareKeeper, options.BankKeeper, options.WasmKeeper),
		wasmpost.NewWasmdDecorator(options.WasmKeeper),
	}

	return sdk.ChainPostDecorators(postDecorators...)
}
