package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type SmartAccountsDecorator struct{}

func NewFeeSharePayoutDecorator() SmartAccountsDecorator {
	return SmartAccountsDecorator{}
}

// FeeSharePostHandler if the smartaccounts module is enabled
// takes the total fees paid for each transaction and
// split these fees equally between all the contacts
// involved in the transaction based on the module params.
func (sad SmartAccountsDecorator) PostHandle(
	ctx sdk.Context,
	tx sdk.Tx,
	simulate bool,
	success bool,
	next sdk.PostHandler,
) (newCtx sdk.Context, err error) {

	return next(ctx, tx, simulate, success)
}
