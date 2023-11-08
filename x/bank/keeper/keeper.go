package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	accountkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	errorsmod "cosmossdk.io/errors"
	custombankkeeper "github.com/terra-money/alliance/custom/bank/keeper"
	customterratypes "github.com/terra-money/core/v2/x/bank/types"
)

type Keeper struct {
	custombankkeeper.Keeper
	hooks customterratypes.BankHooks
	ak    accountkeeper.AccountKeeper
}

var _ bankkeeper.Keeper = Keeper{}

func NewBaseKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	ak accountkeeper.AccountKeeper,
	blockedAddrs map[string]bool,
	authority string,
) Keeper {
	keeper := Keeper{
		Keeper: custombankkeeper.NewBaseKeeper(cdc, storeKey, ak, blockedAddrs, authority),
		hooks:  nil,
		ak:     ak,
	}

	return keeper
}

// Set the bank hooks
func (k *Keeper) SetHooks(bh customterratypes.BankHooks) *Keeper {
	if k.hooks != nil {
		panic("cannot set bank hooks twice")
	}

	k.hooks = bh

	return k
}

// SendCoins transfers amt coins from a sending account to a receiving account.
// An error is returned upon failure.
func (k Keeper) SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error {
	err := k.BlockBeforeSend(ctx, fromAddr, toAddr, amt)
	if err != nil {
		return err
	}
	k.TrackBeforeSend(ctx, fromAddr, toAddr, amt)

	return k.Keeper.SendCoins(ctx, fromAddr, toAddr, amt)
}

// SendCoinsFromModuleToManyAccounts transfers coins from a ModuleAccount to multiple AccAddresses.
// It will panic if the module account does not exist. An error is returned if
// the recipient address is black-listed or if sending the tokens fails.
func (k Keeper) SendCoinsFromModuleToModule(ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coins) error {
	senderAddr := k.ak.GetModuleAddress(senderModule)
	if senderAddr == nil {
		panic(errorsmod.Wrapf(customterratypes.ErrUnknownAddress, "senderModule address %s is nil", senderModule))
	}
	recipientAcc := k.ak.GetModuleAccount(ctx, recipientModule)
	if recipientAcc == nil {
		panic(errorsmod.Wrapf(customterratypes.ErrUnknownAddress, "recipientModule address %s is nil", recipientModule))
	}

	return k.Keeper.SendCoins(ctx, senderAddr, recipientAcc.GetAddress(), amt)
}

// UndelegateCoins performs undelegation by crediting amt coins to an account with
// address addr. For vesting accounts, undelegation amounts are tracked for both
// vesting and vested coins. The coins are then transferred from a ModuleAccount
// address to the delegator address. If any of the undelegation amounts are
// negative, an error is returned.
func (k Keeper) UndelegateCoins(ctx sdk.Context, moduleAccAddr, delegatorAddr sdk.AccAddress, amt sdk.Coins) error {
	err := k.BlockBeforeSend(ctx, moduleAccAddr, delegatorAddr, amt)
	if err != nil {
		return err
	}
	k.TrackBeforeSend(ctx, moduleAccAddr, delegatorAddr, amt)

	return k.Keeper.UndelegateCoins(ctx, moduleAccAddr, delegatorAddr, amt)
}

// DelegateCoins performs delegation by deducting amt coins from an account with
// address addr. For vesting accounts, delegations amounts are tracked for both
// vesting and vested coins. The coins are then transferred from the delegator
// address to a ModuleAccount address. If any of the delegation amounts are negative,
// an error is returned.
func (k Keeper) DelegateCoins(ctx sdk.Context, delegatorAddr, moduleAccAddr sdk.AccAddress, amt sdk.Coins) error {
	err := k.BlockBeforeSend(ctx, delegatorAddr, moduleAccAddr, amt)
	if err != nil {
		return err
	}
	k.TrackBeforeSend(ctx, delegatorAddr, moduleAccAddr, amt)

	return k.Keeper.DelegateCoins(ctx, delegatorAddr, moduleAccAddr, amt)
}

// InputOutputCoins performs multi-send functionality. It accepts a series of
// inputs that correspond to a series of outputs. It returns an error if the
// inputs and outputs don't line up or if any single transfer of tokens fails.
func (k Keeper) InputOutputCoins(ctx sdk.Context, inputs []banktypes.Input, outputs []banktypes.Output) error {
	// Only 1 input is allowed for all outputs check the following url:
	// https://github.com/terra-money/cosmos-sdk/blob/release/v0.47.x/x/bank/types/msgs.go#L87-L89
	//
	// This if statement is added here too so we know
	// when multiple inputs are allowed in the future
	// because ErrMultipleSenders will fail to import
	// because will be removed from the code.
	if len(inputs) != 1 {
		return banktypes.ErrMultipleSenders
	}
	input := inputs[0]
	inputaddress := sdk.MustAccAddressFromBech32(input.Address)

	for _, output := range outputs {
		outputaddress := sdk.MustAccAddressFromBech32(output.Address)

		err := k.BlockBeforeSend(ctx, inputaddress, outputaddress, output.Coins)
		if err != nil {
			return err
		}
		k.TrackBeforeSend(ctx, inputaddress, outputaddress, output.Coins)
	}

	return k.Keeper.InputOutputCoins(ctx, inputs, outputs)
}
