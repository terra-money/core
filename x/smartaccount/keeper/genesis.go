package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/terra-money/core/v2/x/smartaccount/types"
)

// InitGenesis import module genesis
func (k Keeper) InitGenesis(
	ctx sdk.Context,
	genesisState types.GenesisState,
) {
	if err := k.SetParams(ctx, genesisState.Params); err != nil {
		panic(err)
	}

	for _, setting := range genesisState.Settings {
		if err := k.SetSetting(ctx, setting.Owner, *setting); err != nil {
			panic(err)
		}
	}
}

// ExportGenesis export module state
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.KeyPrefixSetting)

	settings := []*types.Setting{}

	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var setting types.Setting
		setting.Unmarshal(iter.Value())
		settings = append(settings, &setting)
	}

	return &types.GenesisState{
		Params:   k.GetParams(ctx),
		Settings: settings,
	}
}
