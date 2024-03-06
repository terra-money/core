package post

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	feemarketpost "github.com/skip-mev/feemarket/x/feemarket/post"
	feesharepost "github.com/terra-money/core/v2/x/feeshare/post"
	customwasmkeeper "github.com/terra-money/core/v2/x/wasm/keeper"
	wasmpost "github.com/terra-money/core/v2/x/wasm/post"
)

type HandlerOptions struct {
	FeeShareKeeper feesharepost.FeeShareKeeper
	BankKeeper     feesharepost.BankKeeper
	WasmKeeper     customwasmkeeper.Keeper

	AccountKeeper   feemarketpost.AccountKeeper
	FeeMarketKeeper feemarketpost.FeeMarketKeeper
	FeeGrantKeeper  feemarketpost.FeeGrantKeeper
}

func NewPostHandler(options HandlerOptions) (sdk.PostHandler, error) {

	if options.AccountKeeper == nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "account keeper is required for post builder")
	}

	if options.BankKeeper == nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "bank keeper is required for post builder")
	}

	if options.FeeMarketKeeper == nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "feemarket keeper is required for post builder")
	}

	postDecorators := []sdk.PostDecorator{
		feemarketpost.NewFeeMarketDeductDecorator(
			options.AccountKeeper,
			options.BankKeeper,
			options.FeeGrantKeeper,
			options.FeeMarketKeeper,
		),
		feesharepost.NewFeeSharePayoutDecorator(options.FeeShareKeeper, options.BankKeeper, options.WasmKeeper),
		wasmpost.NewWasmdDecorator(options.WasmKeeper),
	}

	return sdk.ChainPostDecorators(postDecorators...), nil
}
