package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// SmartAccountCheckDecorator does authentication for smart accounts
type SmartAccountCheckDecorator struct {
	smartaccountKeeper SmartAccountKeeper
}

func NewFeeMarketCheckDecorator(sak SmartAccountKeeper) SmartAccountCheckDecorator {
	return SmartAccountCheckDecorator{
		smartaccountKeeper: sak,
	}
}

// AnteHandle checks if the tx provides sufficient fee to cover the required fee from the fee market.
func (sad SmartAccountCheckDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	// check if the tx is from a smart account
	setting, err := sad.smartaccountKeeper.GetSetting(ctx, tx.GetMsgs()[0].GetSigners()[0].String())
	if sdkerrors.ErrKeyNotFound.Is(err) {
		return next(ctx, tx, simulate)
	} else if err != nil {
		return ctx, err
	}

	if len(setting.Authorization) > 0 {
		for _, auth := range setting.Authorization {
			_ = auth
			// TODO: add code that calls authorization on contracts
		}
	}

	if len(setting.PreTransaction) > 0 {
		for _, preTx := range setting.PreTransaction {
			_ = preTx
			// TODO: add code that calls pre-transaction on contracts
		}
	}

	return next(ctx, tx, simulate)
}
