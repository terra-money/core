package v3

import (
	"github.com/terra-money/core/v2/x/tokenfactory/exported"
	"github.com/terra-money/core/v2/x/tokenfactory/types"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MigrateStore migrates the x/tokenfactory module state from the consensus version 2 to
// version 3. Specifically, it takes the parameters that are currently stored
// and managed by the x/params module and stores them directly into the x/tokenfactory
// module state.
func MigrateStore(ctx sdk.Context, storeKey storetypes.StoreKey, legacySubspace exported.Subspace, cdc codec.BinaryCodec) error {
	store := ctx.KVStore(storeKey)
	var currParams types.Params
	legacySubspace.GetParamSetIfExists(ctx, &currParams)

	if err := currParams.Validate(); err != nil {
		return err
	}

	bz, err := cdc.Marshal(&currParams)
	if err != nil {
		return err
	}

	store.Set(types.ParamsKey, bz)

	return nil
}
