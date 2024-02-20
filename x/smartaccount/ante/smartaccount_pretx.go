package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// SmartAccountPreTxDecorator does authentication for smart accounts
type SmartAccountPreTxDecorator struct {
	sak SmartAccountKeeper
	wk  WasmKeeper
}

func NewSmartAccountPreTxDecorator(
	sak SmartAccountKeeper,
	wk WasmKeeper,
) SmartAccountPreTxDecorator {
	return SmartAccountPreTxDecorator{
		sak: sak,
		wk:  wk,
	}
}

// AnteHandle checks if the tx provides sufficient fee to cover the required fee from the fee market.
func (sad SmartAccountPreTxDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	// check if the tx is from a smart account
	setting, err := sad.sak.GetSetting(ctx, tx.GetMsgs()[0].GetSigners()[0].String())
	if sdkerrors.ErrKeyNotFound.Is(err) {
		return next(ctx, tx, simulate)
	} else if err != nil {
		return ctx, err
	}

	if setting.PreTransaction != nil && len(setting.PreTransaction) > 0 {
		for _, preTx := range setting.PreTransaction {
			_ = preTx
			// TODO: add code that calls pre-transaction on contracts
		}
	}

	return next(ctx, tx, simulate)
}
