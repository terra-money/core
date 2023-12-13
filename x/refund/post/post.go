package post

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/terra-money/core/v2/x/refund/types"
)

type RefundDecorator struct {
	bankkeeper types.BankKeeper
}

func NewRefundDecorator(bankkeeper types.BankKeeper) RefundDecorator {
	return RefundDecorator{bankkeeper}
}

// CONTRACT: this functiona ssumes that feegrant is enabled.
func (rd RefundDecorator) PostHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, success bool, next sdk.PostHandler) (newCtx sdk.Context, err error) {
	// Parse sdk.Tx to sdk.FeeTx that way we can get the information
	// about the transaction fees and the fee payer.
	feeTx, ok := tx.(sdk.FeeTx)
	if !ok {
		return ctx, errorsmod.Wrap(sdkerrors.ErrTxDecode, "Tx must be a FeeTx")
	}

	// Get the fee and check if it is zero continue with the next
	// decorator because there is nothing to refund to the user
	txFees := feeTx.GetFee()
	if txFees.IsZero() {
		return next(ctx, tx, simulate, success)
	}

	// Get the fee payer and the fee granter, assume that the fees are
	// returned to the fee payer unless there is a fee granter which
	// will receive the proportional part of the unused fees.
	refundBeneficiary := feeTx.FeePayer()
	feeGranter := feeTx.FeeGranter()
	if feeGranter != nil {
		refundBeneficiary = feeGranter
	}

	gasMeter := ctx.GasMeter()
	gasConsumed := gasMeter.GasConsumed()
	gasLimit := gasMeter.Limit()
	parsedGasConsumed := math.LegacyNewDecFromInt(sdk.NewInt(int64(gasConsumed)))
	parsedGasLimit := math.LegacyNewDecFromInt(sdk.NewInt(int64(gasLimit)))

	var toRefund sdk.Coins
	for _, txFee := range txFees {
		multiplier := parsedGasLimit.Quo(parsedGasConsumed)

		amountToRefund := txFee.Amount.MulRaw(multiplier.RoundInt().ToLegacyDec().TruncateInt64())

		if amountToRefund.IsPositive() {
			toRefund = append(toRefund, sdk.NewCoin(txFee.Denom, amountToRefund))
		}
	}

	if toRefund.IsZero() {
		return next(ctx, tx, simulate, success)
	}

	err = rd.bankkeeper.SendCoinsFromModuleToAccount(ctx, authtypes.FeeCollectorName, refundBeneficiary, toRefund)
	if err != nil {
		return ctx, err
	}

	return next(ctx, tx, simulate, success)
}
