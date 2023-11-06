package ante

import (
	"slices"

	errorsmod "cosmossdk.io/errors"

	"github.com/cosmos/cosmos-sdk/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

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
// split these fees equaly between all the contacts involved in the
// transaction based on the module params.
func (fsd FeeSharePayoutDecorator) FeeSharePayout(ctx sdk.Context, txFees sdk.Coins, devShares types.Dec, allowedDenoms []string) (err error) {
	events := ctx.EventManager().Events()
	contractAddresses, err := ExtractContractAddrs(events)
	if err != nil {
		return err
	}
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
		ctx.EventManager().EmitTypedEvent(
			&feeshare.FeePayoutEvent{
				WithdrawAddress: withdrawerAddrs.String(),
				FeesPaid:        feeToBePaid,
			},
		)
	}

	return err
}

// Iterate the events and search for the execute event then iterate the
// attributes in search for _contract_address and get the value which
// is the contract address to search for all the beneficiaries, info:
// https://github.com/CosmWasm/wasmd/blob/main/EVENTS.md#validation-rules
func ExtractContractAddrs(events sdk.Events) ([]string, error) {
	contractAddresses := []string{}
	for _, ev := range events {
		if ev.Type != "execute" {
			continue
		}

		for _, attr := range ev.Attributes {
			if attr.Key != "_contract_address" {
				continue
			}
			// if the contract address has already been
			// added just skip it to avoid duplicates
			if slices.Contains(contractAddresses, attr.Value) {
				continue
			}

			contractAddresses = append(contractAddresses, attr.Value)
		}
	}

	return contractAddresses, nil
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
func CalculateFee(fees sdk.Coins, govPercent sdk.Dec, pairs int, allowedDenoms []string) sdk.Coins {
	var alloedFeesDenoms sdk.Coins
	if len(allowedDenoms) == 0 {
		alloedFeesDenoms = fees
	} else {
		for _, fee := range fees {
			if slices.Contains(allowedDenoms, fee.Denom) {
				alloedFeesDenoms = alloedFeesDenoms.Add(fee)
			}
		}
	}

	var splitFees sdk.Coins
	for _, c := range alloedFeesDenoms.Sort() {
		rewardAmount := govPercent.MulInt(c.Amount).QuoInt64(int64(pairs)).RoundInt()
		if !rewardAmount.IsZero() {
			splitFees = splitFees.Add(sdk.NewCoin(c.Denom, rewardAmount))
		}
	}
	return splitFees
}
