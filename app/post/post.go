package post

import (
	feesharepost "github.com/terra-money/core/v2/x/feeshare/post"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type HandlerOptions struct {
	FeeShareKeeper feesharepost.FeeShareKeeper
	BankKeeper     feesharepost.BankKeeper
}

func NewPostHandler(options HandlerOptions) sdk.PostHandler {

	postDecorators := []sdk.PostDecorator{
		feesharepost.NewFeeSharePayoutDecorator(options.FeeShareKeeper, options.BankKeeper),
	}

	return sdk.ChainPostDecorators(postDecorators...)
}
