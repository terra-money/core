package ante

import (
	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	feeshare "github.com/terra-money/core/v2/x/feeshare/types"
	customwasmkeeper "github.com/terra-money/core/v2/x/wasm/keeper"
)

type FeeSharePayoutDecorator struct {
	feesharekeeper FeeShareKeeper
	bankKeeper     BankKeeper
	wasmKeeper     customwasmkeeper.Keeper
}

func NewFeeSharePayoutDecorator(fs FeeShareKeeper, bk BankKeeper, wk customwasmkeeper.Keeper) FeeSharePayoutDecorator {
	return FeeSharePayoutDecorator{
		feesharekeeper: fs,
		bankKeeper:     bk,
		wasmKeeper:     wk,
	}
}

// FeeSharePostHandler if the feeshare module is enabled
// takes the total fees paid for each transaction and
// split these fees equally between all the contacts
// involved in the transaction based on the module params.
func (fsd FeeSharePayoutDecorator) PostHandle(
	ctx sdk.Context,
	tx sdk.Tx,
	simulate bool,
	success bool,
	next sdk.PostHandler,
) (newCtx sdk.Context, err error) {
	// Check if fee share is enabled and shares
	// to be distribute are greater than zero.
	params := fsd.feesharekeeper.GetParams(ctx)
	if !params.EnableFeeShare || params.DeveloperShares.IsZero() {
		return next(ctx, tx, simulate, success)
	}

	// Parse the transactions to FeeTx to get the total fees paid
	feeTx, ok := tx.(sdk.FeeTx)
	if !ok {
		return ctx, errorsmod.Wrap(sdkerrors.ErrTxDecode, "Tx must be a FeeTx")
	}
	if feeTx.GetFee().Empty() || feeTx.GetFee().IsZero() {
		return next(ctx, tx, simulate, success)
	}

	err = fsd.FeeSharePayout(ctx, feeTx.GetFee(), params.DeveloperShares, params.AllowedDenoms)
	if err != nil {
		return ctx, errorsmod.Wrapf(sdkerrors.ErrLogic, err.Error())
	}

	return next(ctx, tx, simulate, success)
}

// FeeSharePayout takes the total fees paid for a transaction and
// split these fees equally between all the contacts involved in the
// transaction based on the module params.
func (fsd FeeSharePayoutDecorator) FeeSharePayout(ctx sdk.Context, txFees sdk.Coins, devShares sdk.Dec, allowedDenoms []string) (err error) {
	executedContracts, found := fsd.wasmKeeper.GetExecutedContractAddresses(ctx)
	if !found {
		return err
	}
	contractAddresses := executedContracts.ContractAddresses
	if len(contractAddresses) == 0 {
		return err
	}

	withdrawerAddrs, err := GetWithdrawalAddressFromContract(ctx, contractAddresses, fsd.feesharekeeper)
	if err != nil {
		return err
	}
	if len(withdrawerAddrs) == 0 {
		return err
	}

	feeToBePaid := CalculateFee(txFees, devShares, len(withdrawerAddrs), allowedDenoms)
	if feeToBePaid.IsZero() {
		return err
	}

	// pay the fees to the withdrawer addresses
	for _, withdrawerAddrs := range withdrawerAddrs {
		err = fsd.bankKeeper.SendCoinsFromModuleToAccount(ctx, authtypes.FeeCollectorName, withdrawerAddrs, feeToBePaid)
		if err != nil {
			return err
		}
		err := ctx.EventManager().EmitTypedEvent(
			&feeshare.FeePayoutEvent{
				WithdrawAddress: withdrawerAddrs.String(),
				FeesPaid:        feeToBePaid,
			},
		)
		if err != nil {
			return err
		}
	}

	return err
}

// Iterate the contract addresses and get the
// withdrawer address from the module store
func GetWithdrawalAddressFromContract(ctx sdk.Context, contractAddresses []string, fsk FeeShareKeeper) ([]sdk.AccAddress, error) {
	var withdrawerAddrs []sdk.AccAddress

	for _, contractAddr := range contractAddresses {
		parsedContractAddr, err := sdk.AccAddressFromBech32(contractAddr)
		if err != nil {
			return nil, err
		}

		shareData, hasfeeshare := fsk.GetFeeShare(ctx, parsedContractAddr)

		if !hasfeeshare {
			continue
		}

		withdrawerAddr := shareData.GetWithdrawerAddr()
		if withdrawerAddr != nil && !withdrawerAddr.Empty() {
			withdrawerAddrs = append(withdrawerAddrs, withdrawerAddr)
		}
	}

	return withdrawerAddrs, nil
}

// CalculateFee takes the total fees paid for a transaction and split
// these fees equaly between all number of pairs considering allwoedDenoms
func CalculateFee(fees sdk.Coins, devShares sdk.Dec, numOfdevs int, allowedDenoms []string) sdk.Coins {
	var allowedFeesDenoms sdk.Coins
	if len(allowedDenoms) == 0 {
		allowedFeesDenoms = fees
	} else {
		for _, fee := range fees {
			for _, allowedDenom := range allowedDenoms {
				if fee.Denom == allowedDenom {
					allowedFeesDenoms = allowedFeesDenoms.Add(fee)
					break
				}
			}
		}
	}

	var splitFees sdk.Coins
	for _, c := range allowedFeesDenoms.Sort() {
		rewardAmount := devShares.MulInt(c.Amount).QuoInt64(int64(numOfdevs)).RoundInt()
		if !rewardAmount.IsZero() {
			splitFees = splitFees.Add(sdk.NewCoin(c.Denom, rewardAmount))
		}
	}
	return splitFees
}
