package v2_11

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	ibctransferkeeper "github.com/cosmos/ibc-go/v7/modules/apps/transfer/keeper"
	transfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	custombankkeeper "github.com/terra-money/core/v2/x/bank/keeper"
)

type EscrowUpdate struct {
	EscrowAddress sdk.AccAddress
	Assets        []sdk.Coin
}

func CreateUpgradeHandler(
	mm *module.Manager,
	cfg module.Configurator,
	bankKeeper custombankkeeper.Keeper,
	transferKeeper ibctransferkeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		if ctx.ChainID() != "phoenix-1" {
			return mm.RunMigrations(ctx, cfg, vm)
		}
		// This slice is initialized with objects that describe the escrow account address and the coins that need to be minted to fix the discrepancy.
		// To find the escrow account and address have been used the following software:
		// https://github.com/strangelove-ventures/escrow-checker/commit/adf0d867e2210c9ff0a27d8dff1c74ed0c8a00dc
		updates := []EscrowUpdate{
			{
				EscrowAddress: sdk.AccAddress("terra1s308jav50mgct9x4f87u23w2tfe8q6qe45y7s4"),
				Assets:        []sdk.Coin{sdk.NewCoin("ibc/815FC81EB6BD612206BD9A9909A02F7691D24A5B97CDFE2124B1BDCA9D4AB14C", sdk.NewInt(1000000000))},
			},
		}

		for _, update := range updates {
			for _, coin := range update.Assets {
				coins := sdk.NewCoins(coin)

				if err := bankKeeper.MintCoins(ctx, transfertypes.ModuleName, coins); err != nil {
					return nil, err
				}

				if err := bankKeeper.SendCoinsFromModuleToAccount(ctx, transfertypes.ModuleName, update.EscrowAddress, coins); err != nil {
					return nil, err
				}

				// For ibc-go v7+ you will also need to update the transfer module's store for the total escrow amounts.
				currentTotalEscrow := transferKeeper.GetTotalEscrowForDenom(ctx, coin.GetDenom())
				newTotalEscrow := currentTotalEscrow.Add(coin)
				transferKeeper.SetTotalEscrowForDenom(ctx, newTotalEscrow)
			}
		}

		return mm.RunMigrations(ctx, cfg, vm)
	}
}
