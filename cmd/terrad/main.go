package main

import (
	"os"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/starport/starport/pkg/cosmoscmd"
	"github.com/terra-money/core/app"
)

func main() {
	config := sdk.GetConfig()
	config.SetCoinType(app.CoinType)

	cmdOptions := GetWasmCmdOptions()
	rootCmd, _ := cosmoscmd.NewRootCmd(
		app.AppName,
		app.AccountAddressPrefix,
		app.DefaultNodeHome,
		app.AppName,
		app.ModuleBasics,
		app.NewTerraApp,
		// this line is used by starport scaffolding # root/arguments
		cmdOptions...,
	)
	if err := svrcmd.Execute(rootCmd, app.DefaultNodeHome); err != nil {
		os.Exit(1)
	}
}
