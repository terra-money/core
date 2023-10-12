package v2_6

import (
	"time"

	sdkerrors "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	clientkeeper "github.com/cosmos/ibc-go/v7/modules/core/02-client/keeper"
	ibcclienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	ibcexported "github.com/cosmos/ibc-go/v7/modules/core/exported"
	ibctm "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	cfg module.Configurator,
	cdc codec.Codec,
	clientKeeper clientkeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		err := increaseUnbondingPeriod(ctx, cdc, clientKeeper)
		if err != nil {
			return nil, err
		}
		return mm.RunMigrations(ctx, cfg, fromVM)
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
