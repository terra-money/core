package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	accountkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	custombankkeeper "github.com/terra-money/alliance/custom/bank/keeper"
	customterratypes "github.com/terra-money/core/v2/custom/bank/types"
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
	// BlockBeforeSend hook should always be called before the TrackBeforeSend hook.
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
		panic(sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "module account %s does not exist", senderModule))
	}

	recipientAcc := k.ak.GetModuleAccount(ctx, recipientModule)
	if recipientAcc == nil {
		panic(sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "module account %s does not exist", recipientModule))
	}

	return k.Keeper.SendCoins(ctx, senderAddr, recipientAcc.GetAddress(), amt)
}
