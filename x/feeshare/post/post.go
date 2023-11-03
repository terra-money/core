package ante

import (
	"encoding/json"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/authz"

	"github.com/terra-money/core/v2/x/feeshare/types"
	feeshare "github.com/terra-money/core/v2/x/feeshare/types"
)

type FeeSharePayoutDecorator struct {
	feesharekeeper FeeShareKeeper
	bankKeeper     BankKeeper
}

func NewFeeSharePayoutDecorator(fs FeeShareKeeper, bk BankKeeper) FeeSharePayoutDecorator {
	return FeeSharePayoutDecorator{
		feesharekeeper: fs,
		bankKeeper:     bk,
	}
}

// FeeSharePostHandler if the feeshare module is neabled
// takes the total fees paid for each transaction and
// split these fees equaly between all the contacts
// involved in the transactin based on the module params.
func (fsd FeeSharePayoutDecorator) PostHandle(
	ctx sdk.Context,
	tx sdk.Tx, simulate,
	success bool,
	next sdk.PostHandler,
) (newCtx sdk.Context, err error) {
	// Check if fee share is enabled
	params := fsd.feesharekeeper.GetParams(ctx)
	if !params.EnableFeeShare {
		return ctx, nil
	}
	// Parse the transactions to FeeTx
	// to get the total fees paid in the future
	feeTx, ok := tx.(sdk.FeeTx)
	if !ok {
		return ctx, errorsmod.Wrap(sdkerrors.ErrTxDecode, "Tx must be a FeeTx")
	}

	err = fsd.feeSharePayout(ctx, fsd.bankKeeper, fsd.feesharekeeper, feeTx)
	if err != nil {
		return ctx, errorsmod.Wrapf(sdkerrors.ErrInsufficientFunds, err.Error())
	}

	return next(ctx, tx, simulate, success)
}
func (fsd FeeSharePayoutDecorator) feeSharePayout(
	ctx sdk.Context,
	bankKeeper BankKeeper,
	feesharekeeper FeeShareKeeper,
	msgs sdk.FeeTx,
) error {
	// Get valid withdraw addresses from contracts
	toPay := make([]sdk.AccAddress, 0)

	// validAuthz checks if the msg is an authz exec msg and if so, call the validExecuteMsg for each
	// inner msg. If it is a CosmWasm execute message, that logic runs for nested functions.
	validAuthz := func(execMsg *authz.MsgExec) error {
		for _, v := range execMsg.Msgs {
			var innerMsg sdk.Msg
			if err := json.Unmarshal(v.Value, &innerMsg); err != nil {
				return errorsmod.Wrapf(sdkerrors.ErrUnauthorized, "cannot unmarshal authz exec msgs")
			}

			err := addNewFeeSharePayoutsForMsg(ctx, feesharekeeper, &toPay, innerMsg)
			if err != nil {
				return err
			}
		}

		return nil
	}

	for _, m := range msgs {
		if msg, ok := m.(*authz.MsgExec); ok {
			if err := validAuthz(msg); err != nil {
				return err
			}
			continue
		}

		if err := addNewFeeSharePayoutsForMsg(ctx, feesharekeeper, &toPay, m); err != nil {
			return err
		}
	}

	// Do nothing if no one needs payment
	if len(toPay) == 0 {
		return nil
	}

	// Get only allowed governance fees to be paid (helps for taxes)
	var fees sdk.Coins
	if len(params.AllowedDenoms) == 0 {
		// If empty, we allow all denoms to be used as payment
		fees = txFees
	} else {
		for _, fee := range txFees.Sort() {
			for _, allowed := range params.AllowedDenoms {
				if fee.Denom == allowed {
					fees = fees.Add(fee)
				}
			}
		}
	}

	numPairs := len(toPay)

	feesPaidOutput := make([]types.FeePayoutEvent, numPairs)
	if numPairs > 0 {
		govPercent := params.DeveloperShares
		splitFees := FeePayLogic(fees, govPercent, numPairs)

		// pay fees evenly between all withdraw addresses
		for i, withdrawAddr := range toPay {
			err := bankKeeper.SendCoinsFromModuleToAccount(ctx, authtypes.FeeCollectorName, withdrawAddr, splitFees)
			feesPaidOutput[i] = types.FeePayoutEvent{
				WithdrawAddress: withdrawAddr.String(),
				FeesPaid:        splitFees,
			}

			if err != nil {
				return errorsmod.Wrapf(feeshare.ErrFeeSharePayment, "failed to pay fees to contract developer: %s", err.Error())
			}
		}
	}

	bz, err := json.Marshal(feesPaidOutput)
	if err != nil {
		return errorsmod.Wrapf(feeshare.ErrFeeSharePayment, "failed to marshal feesPaidOutput: %s", err.Error())
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			feeshare.EventTypePayoutFeeShare,
			sdk.NewAttribute(feeshare.AttributeWithdrawPayouts, string(bz))),
	)

	return nil
}

func FeePayLogic(fees sdk.Coins, govPercent sdk.Dec, numPairs int) sdk.Coins {
	var splitFees sdk.Coins
	for _, c := range fees.Sort() {
		rewardAmount := govPercent.MulInt(c.Amount).QuoInt64(int64(numPairs)).RoundInt()
		if !rewardAmount.IsZero() {
			splitFees = splitFees.Add(sdk.NewCoin(c.Denom, rewardAmount))
		}
	}
	return splitFees
}

func addNewFeeSharePayoutsForMsg(ctx sdk.Context, feesharekeeper FeeShareKeeper, toPay *[]sdk.AccAddress, m sdk.Msg) error {
	if msg, ok := m.(*wasmtypes.MsgExecuteContract); ok {
		contractAddr, err := sdk.AccAddressFromBech32(msg.Contract)
		if err != nil {
			return err
		}

		shareData, _ := feesharekeeper.GetFeeShare(ctx, contractAddr)

		withdrawAddr := shareData.GetWithdrawerAddr()
		if withdrawAddr != nil && !withdrawAddr.Empty() {
			*toPay = append(*toPay, withdrawAddr)
		}
	}

	return nil
}
