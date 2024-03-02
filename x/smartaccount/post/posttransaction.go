package post

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/terra-money/core/v2/x/smartaccount/ante"
	"github.com/terra-money/core/v2/x/smartaccount/types"
)

// PostTransactionHookDecorator does authentication for smart accounts
type PostTransactionHookDecorator struct {
	smartaccountKeeper SmartAccountKeeper
	wasmKeeper         WasmKeeper
}

func NewPostTransactionHookDecorator(sak SmartAccountKeeper, wk WasmKeeper) PostTransactionHookDecorator {
	return PostTransactionHookDecorator{
		smartaccountKeeper: sak,
		wasmKeeper:         wk,
	}
}

// FeeSharePostHandler if the smartaccount module is enabled
// takes the total fees paid for each transaction and
// split these fees equally between all the contacts
// involved in the transaction based on the module params.
func (pth PostTransactionHookDecorator) PostHandle(
	ctx sdk.Context,
	tx sdk.Tx,
	simulate bool,
	success bool,
	next sdk.PostHandler,
) (newCtx sdk.Context, err error) {
	setting, ok := ctx.Value(types.ModuleName).(types.Setting)
	if !ok {
		return next(ctx, tx, simulate, success)
	}

	if setting.PostTransaction != nil && len(setting.PostTransaction) > 0 {
		for _, postTx := range setting.PostTransaction {
			contractAddr, err := sdk.AccAddressFromBech32(postTx)
			if err != nil {
				return ctx, err
			}
			data, err := BuildPostTransactionHookMsg(tx)
			if err != nil {
				return ctx, err
			}
			_, err = pth.wasmKeeper.Sudo(ctx, contractAddr, data)
			if err != nil {
				return ctx, err
			}
		}
	}
	return next(ctx, tx, simulate, success)
}

func BuildPostTransactionHookMsg(tx sdk.Tx) ([]byte, error) {
	return ante.BuildPrePostTransactionHookMsg(tx, false)
}
