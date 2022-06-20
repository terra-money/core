package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/terra-money/core/v2/app/ante/types"
)

// keeper of the staking store
type Keeper struct {
	storeKey   sdk.StoreKey
	cdc        codec.BinaryCodec
	paramstore paramtypes.Subspace
}

// NewKeeper creates a new staking Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec, key sdk.StoreKey, ps paramtypes.Subspace,
) Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		storeKey:   key,
		cdc:        cdc,
		paramstore: ps,
	}
}

// set the minimum commission
func (k Keeper) SetMinimumCommission(ctx sdk.Context, minimumCommission sdk.Dec) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&sdk.DecProto{Dec: minimumCommission})
	store.Set(types.MinimumCommissionKey, b)
}

// get the minimum commission
func (k Keeper) GetMinimumCommission(ctx sdk.Context) sdk.Dec {
	store := ctx.KVStore(k.storeKey)

	var minimumCommissionProto sdk.DecProto
	b := store.Get(types.MinimumCommissionKey)
	k.cdc.MustUnmarshal(b, &minimumCommissionProto)

	return minimumCommissionProto.Dec
}

// InitGenesis initialize a GenesisState for a given context and keeper
func (k Keeper) InitGenesis(ctx sdk.Context, data *types.GenesisState) {
	k.SetParams(ctx, data.GetParams())
	k.SetMinimumCommission(ctx, data.MinimumCommission)
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	return &types.GenesisState{
		Params:            k.GetParams(ctx),
		MinimumCommission: k.GetMinimumCommission(ctx),
	}
}
