package post

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/terra-money/core/v2/x/feeburn/types"
)

type FeeBurnDecorator struct {
	feeBurnKeeper FeeBurnKeeper
	bankkeeper    BankKeeper
}

func NewFeeBurnDecorator(feeBurnKeeper FeeBurnKeeper, bankkeeper BankKeeper) FeeBurnDecorator {
	return FeeBurnDecorator{feeBurnKeeper, bankkeeper}
}

func (fbd FeeBurnDecorator) PostHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, success bool, next sdk.PostHandler) (newCtx sdk.Context, err error) {
	// if the feeburn is not enabled then just continue with the next decorator
	if !fbd.feeBurnKeeper.GetParams(ctx).EnableFeeBurn {
		return next(ctx, tx, simulate, success)
	}

	// Parse sdk.Tx to sdk.FeeTx that way we can get the information
	// about the transaction fees and the fee payer.
	feeTx, ok := tx.(sdk.FeeTx)
	if !ok {
		return ctx, errorsmod.Wrap(sdkerrors.ErrTxDecode, "Tx must be a FeeTx")
	}

	// Get the fees and check if these fees are zero,
	// when fees are zero continue with the next decorator
	// because there is nothing to burn
	txFees := feeTx.GetFee()
	if txFees.IsZero() {
		return next(ctx, tx, simulate, success)
	}

	// Give the gas meter, we can calculate how much gas to refund to the user
	// when the Gas Limit is set to zero it means that it's infinite and we
	// will not refund anything to the user:
	// README: https://docs.cosmos.network/main/learn/beginner/gas-fees#gas-meter
	gasMeter := ctx.GasMeter()
	if gasMeter.Limit() == 0 {
		return next(ctx, tx, simulate, success)
	}
	gasLimit := math.LegacyNewDecFromInt(sdk.NewInt(int64(gasMeter.Limit())))
	remainingGas := math.LegacyNewDecFromInt(sdk.NewInt(int64(gasMeter.GasRemaining())))
	// Percentage of unused gas for this specific denom
	burnRate := remainingGas.Quo(gasLimit)
	fmt.Printf("remainingGas: %s\n", remainingGas)
	fmt.Printf("gasLimit: %s\n", gasLimit)
	fmt.Printf("burnRate: %s\n", burnRate)

	var toBurn sdk.Coins
	// Iterate over the transaction fees and calculate the proportional part
	// of the unused fees denominated in the tokens used to pay fo the fees
	// and add it to the toBurn variable that will be sent to the user.
	for _, txFee := range txFees {

		// Given the percentage of unused gas, calculate the
		// proportional part of the fees that will be refunded
		// to the user in this specific denom.
		unusedFees := math.LegacyNewDecFromInt(txFee.Amount).
			Mul(burnRate).
			TruncateInt()

		fmt.Printf("math.LegacyNewDecFromInt(txFee.Amount) %s\n", math.LegacyNewDecFromInt(txFee.Amount))

		// When the unused fees are positive it means that the user
		// will receive a refund in this specific denom to its wallet.
		if unusedFees.IsPositive() {
			toBurn = append(toBurn, sdk.NewCoin(txFee.Denom, unusedFees))
		}
	}
	if toBurn.IsZero() {
		return next(ctx, tx, simulate, success)
	}

	// Execute the refund to the user, if there is an error
	// return the error otherwise continue with the execution
	err = fbd.bankkeeper.BurnCoins(ctx, authtypes.FeeCollectorName, toBurn)
	if err != nil {
		return ctx, err
	}

	// Emit an event to be able to track the fees burned
	// because there can be a little bit of discrepancy
	// between fees used and burned, because the proportional
	// part of the fees is truncated to an integer.
	err = ctx.EventManager().EmitTypedEvent(
		&types.FeeBurnEvent{
			FeesBurn: toBurn,
			BurnRate: burnRate,
		},
	)
	if err != nil {
		return ctx, err
	}
	return next(ctx, tx, simulate, success)
}
