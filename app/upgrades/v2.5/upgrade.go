package v2_5

import (
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	consensuskeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	icacontrollerkeeper "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/keeper"
	icacontrollertypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/types"
	clientkeeper "github.com/cosmos/ibc-go/v7/modules/core/02-client/keeper"
	ibcexported "github.com/cosmos/ibc-go/v7/modules/core/exported"
	ibctmmigrations "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint/migrations"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	cfg module.Configurator,
	cdc codec.Codec,
	clientKeeper clientkeeper.Keeper,
	paramsKeeper paramskeeper.Keeper,
	consensusParamsKeeper consensuskeeper.Keeper,
	icacontrollerKeeper icacontrollerkeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		// READ: https://github.com/cosmos/cosmos-sdk/blob/v0.47.4/UPGRADING.md#xconsensus
		baseAppLegacySS := paramsKeeper.Subspace(baseapp.Paramspace).
			WithKeyTable(paramstypes.ConsensusParamsKeyTable())
		baseapp.MigrateParams(ctx, baseAppLegacySS, &consensusParamsKeeper)

		// READ: https://github.com/cosmos/ibc-go/blob/v7.2.0/docs/migrations/v6-to-v7.md#chains
		_, err := ibctmmigrations.PruneExpiredConsensusStates(ctx, cdc, clientKeeper)
		if err != nil {
			return nil, err
		}

		// READ: https://github.com/cosmos/ibc-go/blob/v7.2.0/docs/migrations/v7-to-v7_1.md#chains
		params := clientKeeper.GetParams(ctx)
		params.AllowedClients = append(params.AllowedClients, ibcexported.Localhost)
		clientKeeper.SetParams(ctx, params)

		// READ: https://github.com/terra-money/core/issues/166
		icacontrollerKeeper.SetParams(ctx, icacontrollertypes.DefaultParams())

		return mm.RunMigrations(ctx, cfg, fromVM)
	}
}
