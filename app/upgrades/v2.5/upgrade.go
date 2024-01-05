package v2_5

import (
	"time"

	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"

	sdkerrors "cosmossdk.io/errors"
	icacontrollerkeeper "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/keeper"
	icacontrollertypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/types"
	clientkeeper "github.com/cosmos/ibc-go/v7/modules/core/02-client/keeper"
	ibcclienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	ibcexported "github.com/cosmos/ibc-go/v7/modules/core/exported"
	ibctm "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	consensuskeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	cfg module.Configurator,
	cdc codec.Codec,
	clientKeeper clientkeeper.Keeper,
	paramsKeeper paramskeeper.Keeper,
	consensusParamsKeeper consensuskeeper.Keeper,
	icacontrollerKeeper icacontrollerkeeper.Keeper,
	authKeeper authkeeper.AccountKeeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		// READ: https://github.com/cosmos/cosmos-sdk/blob/v0.47.4/UPGRADING.md#xconsensus
		baseAppLegacySS := paramsKeeper.Subspace(baseapp.Paramspace).
			WithKeyTable(paramstypes.ConsensusParamsKeyTable())
		baseapp.MigrateParams(ctx, baseAppLegacySS, &consensusParamsKeeper)

		// READ: https://github.com/cosmos/ibc-go/blob/v7.2.0/docs/migrations/v6-to-v7.md#chains
		// _, err := ibctmmigrations.PruneExpiredConsensusStates(ctx, cdc, clientKeeper)
		// if err != nil {
		// 	return nil, err
		// }
		err := increaseUnbondingPeriod(ctx, cdc, clientKeeper)
		if err != nil {
			return nil, err
		}

		// READ: https://github.com/cosmos/ibc-go/blob/v7.2.0/docs/migrations/v7-to-v7_1.md#chains
		params := clientKeeper.GetParams(ctx)
		params.AllowedClients = append(params.AllowedClients, ibcexported.Localhost)
		clientKeeper.SetParams(ctx, params)

		// READ: https://github.com/terra-money/core/issues/166
		icacontrollerKeeper.SetParams(ctx, icacontrollertypes.DefaultParams())
		vm, err := mm.RunMigrations(ctx, cfg, fromVM)
		if err != nil {
			return nil, err
		}

		return vm, nil
	}
}

// Iterate all IBC clients and increase unbonding period for all atlantic-2 clients
func increaseUnbondingPeriod(ctx sdk.Context, cdc codec.BinaryCodec, clientKeeper clientkeeper.Keeper) error {
	var clientIDs []string
	clientKeeper.IterateClientStates(ctx, []byte(ibcexported.Tendermint), func(clientID string, _ ibcexported.ClientState) bool {
		clientIDs = append(clientIDs, clientID)
		return false
	})

	var totalUpdated int

	for _, clientID := range clientIDs {
		clientState, ok := clientKeeper.GetClientState(ctx, clientID)
		if !ok {
			return sdkerrors.Wrapf(ibcclienttypes.ErrClientNotFound, "clientID %s", clientID)
		}

		tmClientState, ok := clientState.(*ibctm.ClientState)
		if !ok {
			return sdkerrors.Wrap(ibcclienttypes.ErrInvalidClient, "client state is not tendermint even though client id contains 07-tendermint")
		}

		// ATLANTIC 2 blockchain changed the unbonding period on their side,
		// we take advantage of having to upgrade the chain to also increase
		// the unbonding priod on our side.
		if tmClientState.GetChainID() == "atlantic-2" {
			tmClientState.UnbondingPeriod = time.Hour * 24 * 5

			clientKeeper.SetClientState(ctx, clientID, tmClientState)
		}
	}

	clientLogger := clientKeeper.Logger(ctx)
	clientLogger.Info("total ibc clients updated: ", totalUpdated)

	return nil
}
