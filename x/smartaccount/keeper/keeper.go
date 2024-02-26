package keeper

import (
	"fmt"

	"github.com/cometbft/cometbft/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/terra-money/core/v2/x/smartaccount/types"
)

// Keeper of this module maintains collections of smartaccount for contracts
// registered to receive transaction fees.
type Keeper struct {
	storeKey storetypes.StoreKey
	cdc      codec.BinaryCodec
}

// NewKeeper creates new instances of the fees Keeper
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
) Keeper {

	return Keeper{
		storeKey: storeKey,
		cdc:      cdc,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetSetting returns the smart account setting for the ownerAddr
func (k Keeper) GetSetting(ctx sdk.Context, ownerAddr string) (*types.Setting, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetKeyPrefixSetting(ownerAddr))
	if bz == nil {
		return nil, sdkerrors.ErrKeyNotFound.Wrapf("setting not found for ownerAddr: %s", ownerAddr)
	}

	var setting types.Setting
	if err := setting.Unmarshal(bz); err != nil {
		return nil, err
	}

	return &setting, nil
}

// SetSetting sets the smart account setting for the ownerAddr
func (k Keeper) SetSetting(ctx sdk.Context, setting types.Setting) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := setting.Marshal()
	if err != nil {
		return err
	}
	if setting.Owner == "" {
		return sdkerrors.ErrInvalidRequest.Wrap("owner cannot be empty")
	}
	store.Set(types.GetKeyPrefixSetting(setting.Owner), bz)
	return nil
}

// DeleteSetting deletes the smart account setting for the ownerAddr
func (k Keeper) DeleteSetting(ctx sdk.Context, ownerAddr string) error {
	store := ctx.KVStore(k.storeKey)
	if ownerAddr == "" {
		return sdkerrors.ErrInvalidRequest.Wrap("owner cannot be empty")
	}
	bz := store.Get(types.GetKeyPrefixSetting(ownerAddr))
	if bz == nil {
		return sdkerrors.ErrKeyNotFound.Wrapf("setting not found for ownerAddr: %s", ownerAddr)
	}
	store.Delete(types.GetKeyPrefixSetting(ownerAddr))
	return nil
}
