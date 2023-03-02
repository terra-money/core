package v2_3_0

import (
	tokenfactorykeeper "github.com/CosmWasm/wasmd/x/tokenfactory/keeper"
	tokenfactorytypes "github.com/CosmWasm/wasmd/x/tokenfactory/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/terra-money/core/v2/app/config"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	cfg module.Configurator,
	tokenFactoryKeeper tokenfactorykeeper.Keeper) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		currentVm := mm.GetVersionMap()

		// Init token factory with the correct denom
		tokenFactoryKeeper.InitGenesis(ctx, tokenfactorytypes.GenesisState{
			Params: tokenfactorytypes.Params{
				DenomCreationFee: sdk.NewCoins(sdk.NewCoin(config.BondDenom, sdk.NewInt(10_000_000))),
			},
			FactoryDenoms: []tokenfactorytypes.GenesisDenom{},
		})
		fromVM[tokenfactorytypes.ModuleName] = currentVm[tokenfactorytypes.ModuleName]

		return mm.RunMigrations(ctx, cfg, fromVM)
	}
}
