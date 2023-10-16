package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/terra-money/core/v2/x/feeshare/types"
)

// InitGenesis import module genesis
func (k Keeper) InitGenesis(
	ctx sdk.Context,
	data types.GenesisState,
) {
	if err := k.SetParams(ctx, data.Params); err != nil {
		panic(err)
	}

	for _, share := range data.FeeShare {
		contract := share.GetContractAddr()
		deployer := share.GetDeployerAddr()
		withdrawer := share.GetWithdrawerAddr()

		// Set initial contracts receiving transaction fees
		k.SetFeeShare(ctx, share)
		k.SetDeployerMap(ctx, deployer, contract)

		if len(withdrawer) != 0 {
			k.SetWithdrawerMap(ctx, withdrawer, contract)
		}
	}
}

// ExportGenesis export module state
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	return &types.GenesisState{
		Params:   k.GetParams(ctx),
		FeeShare: k.GetFeeShares(ctx),
	}
}
