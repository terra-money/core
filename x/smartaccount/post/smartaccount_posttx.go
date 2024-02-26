package post

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// SmartAccountCheckDecorator does authentication for smart accounts
type SmartAccountPostTxDecorator struct {
	smartaccountKeeper SmartAccountKeeper
}

func NewSmartAccountPostTxDecorator(sak SmartAccountKeeper) SmartAccountPostTxDecorator {
	return SmartAccountPostTxDecorator{
		smartaccountKeeper: sak,
	}
}

// FeeSharePostHandler if the smartaccount module is enabled
// takes the total fees paid for each transaction and
// split these fees equally between all the contacts
// involved in the transaction based on the module params.
func (sad SmartAccountPostTxDecorator) PostHandle(
	ctx sdk.Context,
	tx sdk.Tx,
	simulate bool,
	success bool,
	next sdk.PostHandler,
) (newCtx sdk.Context, err error) {

	setting, err := sad.smartaccountKeeper.GetSetting(ctx, tx.GetMsgs()[0].GetSigners()[0].String())
	if sdkerrors.ErrKeyNotFound.Is(err) {
		return next(ctx, tx, simulate, success)
	} else if err != nil {
		return ctx, err
	}

	if setting.PostTransaction != nil && len(setting.PostTransaction) > 0 {
		for _, postTx := range setting.PostTransaction {
			_ = postTx
			// TODO: add code that calls post-transaction on contracts
		}
	}
	return next(ctx, tx, simulate, success)
}
