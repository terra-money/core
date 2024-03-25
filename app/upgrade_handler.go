package app

import (
	ibchookstypes "github.com/cosmos/ibc-apps/modules/ibc-hooks/v7/types"
	icacontrollertypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/types"
	ibcfeetypes "github.com/cosmos/ibc-go/v7/modules/apps/29-fee/types"
	alliancetypes "github.com/terra-money/alliance/x/alliance/types"
	terraappconfig "github.com/terra-money/core/v2/app/config"
	v2_10 "github.com/terra-money/core/v2/app/upgrades/v2.10"
	v2_11 "github.com/terra-money/core/v2/app/upgrades/v2.11"
	v2_2_0 "github.com/terra-money/core/v2/app/upgrades/v2.2.0"
	v2_3_0 "github.com/terra-money/core/v2/app/upgrades/v2.3.0"
	v2_4 "github.com/terra-money/core/v2/app/upgrades/v2.4"
	v2_5 "github.com/terra-money/core/v2/app/upgrades/v2.5"
	v2_6 "github.com/terra-money/core/v2/app/upgrades/v2.6"
	v2_7 "github.com/terra-money/core/v2/app/upgrades/v2.7"
	v2_8 "github.com/terra-money/core/v2/app/upgrades/v2.8"
	v2_9 "github.com/terra-money/core/v2/app/upgrades/v2.9"
	feesharetypes "github.com/terra-money/core/v2/x/feeshare/types"
	tokenfactorytypes "github.com/terra-money/core/v2/x/tokenfactory/types"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	consensusparamtypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	icqtypes "github.com/cosmos/ibc-apps/modules/async-icq/v7/types"
)

// RegisterUpgradeHandlers returns upgrade handlers
func (app *TerraApp) RegisterUpgradeHandlers() {
	app.Keepers.UpgradeKeeper.SetUpgradeHandler(
		terraappconfig.Upgrade2_2_0,
		v2_2_0.CreateUpgradeHandler(app.GetModuleManager(), app.GetConfigurator()),
	)
	app.Keepers.UpgradeKeeper.SetUpgradeHandler(
		terraappconfig.Upgrade2_3_0,
		v2_3_0.CreateUpgradeHandler(app.GetModuleManager(), app.GetConfigurator(), app.Keepers.TokenFactoryKeeper),
	)
	// This is pisco only since an incorrect plan name was used for the upgrade
	app.Keepers.UpgradeKeeper.SetUpgradeHandler(
		terraappconfig.Upgrade2_4_rc,
		v2_4.CreateUpgradeHandler(app.GetModuleManager(), app.GetConfigurator()),
	)
	app.Keepers.UpgradeKeeper.SetUpgradeHandler(
		terraappconfig.Upgrade2_4,
		v2_4.CreateUpgradeHandler(app.GetModuleManager(), app.GetConfigurator()),
	)
	app.Keepers.UpgradeKeeper.SetUpgradeHandler(
		terraappconfig.Upgrade2_5,
		v2_5.CreateUpgradeHandler(app.GetModuleManager(),
			app.GetConfigurator(),
			app.GetAppCodec(),
			app.Keepers.IBCKeeper.ClientKeeper,
			app.Keepers.ParamsKeeper,
			app.Keepers.ConsensusParamsKeeper,
			app.Keepers.ICAControllerKeeper,
			app.Keepers.AccountKeeper,
		),
	)
	app.Keepers.UpgradeKeeper.SetUpgradeHandler(
		terraappconfig.Upgrade2_6,
		v2_6.CreateUpgradeHandler(app.GetModuleManager(),
			app.GetConfigurator(),
			app.GetAppCodec(),
			app.Keepers.IBCKeeper.ClientKeeper,
			app.Keepers.AccountKeeper,
			app.Keepers.FeeShareKeeper,
		),
	)
	app.Keepers.UpgradeKeeper.SetUpgradeHandler(
		terraappconfig.Upgrade2_7,
		v2_7.CreateUpgradeHandler(
			app.GetModuleManager(),
			app.GetConfigurator(),
			app.GetAppCodec(),
			app.Keepers.ICQKeeper,
		),
	)
	app.Keepers.UpgradeKeeper.SetUpgradeHandler(
		terraappconfig.Upgrade2_8,
		v2_8.CreateUpgradeHandler(
			app.GetModuleManager(),
			app.GetConfigurator(),
			app.GetAppCodec(),
		),
	)
	app.Keepers.UpgradeKeeper.SetUpgradeHandler(
		terraappconfig.Upgrade2_9,
		v2_9.CreateUpgradeHandler(
			app.GetModuleManager(),
			app.GetConfigurator(),
			app.GetAppCodec(),
			app.Keepers.ICQKeeper,
		),
	)
	app.Keepers.UpgradeKeeper.SetUpgradeHandler(
		terraappconfig.Upgrade2_10,
		v2_10.CreateUpgradeHandler(
			app.GetModuleManager(),
			app.GetConfigurator(),
			app.GetAppCodec(),
		),
	)
	app.Keepers.UpgradeKeeper.SetUpgradeHandler(
		terraappconfig.Upgrade2_11,
		v2_11.CreateUpgradeHandler(
			app.GetModuleManager(),
			app.GetConfigurator(),
			app.Keepers.BankKeeper,
			app.Keepers.TransferKeeper,
		),
	)
}

func (app *TerraApp) RegisterUpgradeStores() {
	upgradeInfo, err := app.Keepers.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(err)
	}

	// Add stores for new modules
	if upgradeInfo.Name == terraappconfig.Upgrade2_3_0 && !app.Keepers.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		storeUpgrades := storetypes.StoreUpgrades{
			Added: []string{icacontrollertypes.StoreKey, tokenfactorytypes.StoreKey, ibcfeetypes.StoreKey, ibchookstypes.StoreKey, alliancetypes.StoreKey},
		}
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades))
	} else if upgradeInfo.Name == terraappconfig.Upgrade2_5 && !app.Keepers.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		storeUpgrades := storetypes.StoreUpgrades{
			Added: []string{consensusparamtypes.StoreKey, crisistypes.StoreKey, "builder"},
			// Module intertx removed in v2.5 because it was never used (https://github.com/cosmos/interchain-accounts-demo)
			// The same functionalities are availablein the interchain-accounts under the path
			// integration-tests/src/modules/ica/icav1.test.ts
			Deleted: []string{"intertx"},
		}
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades))
	} else if upgradeInfo.Name == terraappconfig.Upgrade2_6 && !app.Keepers.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		storeUpgrades := storetypes.StoreUpgrades{Added: []string{feesharetypes.StoreKey}}
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades))
	} else if upgradeInfo.Name == terraappconfig.Upgrade2_7 && !app.Keepers.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		storeUpgrades := storetypes.StoreUpgrades{Added: []string{icqtypes.StoreKey}}
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades))
	} else if upgradeInfo.Name == terraappconfig.Upgrade2_9 && !app.Keepers.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		storeUpgrades := storetypes.StoreUpgrades{Deleted: []string{"builder"}, Added: []string{icqtypes.StoreKey}}
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades))
	} else if upgradeInfo.Name == terraappconfig.Upgrade2_10 && !app.Keepers.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		storeUpgrades := storetypes.StoreUpgrades{}
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades))
	} else if upgradeInfo.Name == terraappconfig.Upgrade2_11 && !app.Keepers.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		storeUpgrades := storetypes.StoreUpgrades{}
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades))
	}
}
