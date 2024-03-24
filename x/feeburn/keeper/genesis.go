package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/terra-money/core/v2/x/feeburn/types"
)

// InitGenesis import module genesis
func (k Keeper) InitGenesis(ctx sdk.Context, data types.GenesisState) {
	if err := k.SetParams(ctx, data.Params); err != nil {
		panic(err)
	}
}

// ExportGenesis export module state
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	return &types.GenesisState{
		Params: k.GetParams(ctx),
	}
}
