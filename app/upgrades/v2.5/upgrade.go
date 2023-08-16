package v2_5

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	ibctmmigrations "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint/migrations"
	"github.com/terra-money/core/v2/app"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	cfg module.Configurator,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		// prune expired tendermint consensus states to save storage space
		_, err := ibctmmigrations.PruneExpiredConsensusStates(ctx, app.Codec, app.IBCKeeper.ClientKeeper)
		if err != nil {
			return nil, err
		}

		return app.mm.RunMigrations(ctx, app.configurator, fromVM)
	}
}
