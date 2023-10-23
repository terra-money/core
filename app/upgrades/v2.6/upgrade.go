package v2_6

import (
	"time"

	sdkerrors "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	clientkeeper "github.com/cosmos/ibc-go/v7/modules/core/02-client/keeper"
	ibcclienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	ibcexported "github.com/cosmos/ibc-go/v7/modules/core/exported"
	ibctm "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"
	pobtypes "github.com/skip-mev/pob/x/builder/types"
	feesharekeeper "github.com/terra-money/core/v2/x/feeshare/keeper"
	feesharetypes "github.com/terra-money/core/v2/x/feeshare/types"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	cfg module.Configurator,
	cdc codec.Codec,
	clientKeeper clientkeeper.Keeper,
	authKeeper authkeeper.AccountKeeper,
	feesharekeeper feesharekeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		// feeshare module is a new module added in v2.6,
		// we need to set the default params
		err := feesharekeeper.SetParams(ctx, feesharetypes.DefaultParams())
		if err != nil {
			return nil, err
		}
		// overwrite pob account to a module account for pisco-1
		overwritePobModuleAccount(ctx, authKeeper)

		// Increase the unbonding period for atlantic-2
		err = increaseUnbondingPeriod(ctx, cdc, clientKeeper)
		if err != nil {
			return nil, err
		}
		return mm.RunMigrations(ctx, cfg, fromVM)
	}
}

// Overwrite the module account for pisco-1
func overwritePobModuleAccount(ctx sdk.Context, authKeeper authkeeper.AccountKeeper) {
	if ctx.ChainID() == "pisco-1" {
		macc := authtypes.NewEmptyModuleAccount(pobtypes.ModuleName)
		pobaccount := authKeeper.GetAccount(ctx, macc.GetAddress())
		// if pob account exists, overwrite it
		// if not, create a new one
		if pobaccount != nil {
			macc.AccountNumber = pobaccount.GetAccountNumber()
			maccI := (authKeeper.NewAccount(ctx, macc)).(authtypes.ModuleAccountI)
			authKeeper.SetModuleAccount(ctx, maccI)
		} else {
			maccI := (authKeeper.NewAccount(ctx, macc)).(authtypes.ModuleAccountI)
			authKeeper.SetModuleAccount(ctx, maccI)
		}
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
			tmClientState.UnbondingPeriod = time.Hour * 24 * 21

			clientKeeper.SetClientState(ctx, clientID, tmClientState)
		}
	}

	clientLogger := clientKeeper.Logger(ctx)
	clientLogger.Info("total ibc clients updated: ", totalUpdated)

	return nil
}
