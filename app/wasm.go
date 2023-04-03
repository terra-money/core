package app

import (
	"encoding/json"

	tokenfactorybindings "github.com/CosmWasm/wasmd/x/tokenfactory/bindings"
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/cast"
	alliancebindings "github.com/terra-money/alliance/x/alliance/bindings"

	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func CustomQuerier(app *TerraApp) func(ctx sdk.Context, request json.RawMessage) ([]byte, error) {
	allianceQueryPlugin := alliancebindings.NewAllianceQueryPlugin(app.AllianceKeeper)
	allianceQuerier := alliancebindings.CustomQuerier(allianceQueryPlugin)
	wasmQueryPlugin := tokenfactorybindings.NewQueryPlugin(&app.BankKeeper.BaseKeeper, &app.TokenFactoryKeeper)
	tokenfactoryQuerier := tokenfactorybindings.CustomQuerier(wasmQueryPlugin)
	return func(ctx sdk.Context, request json.RawMessage) ([]byte, error) {
		res, err := allianceQuerier(ctx, request)
		if err != nil {
			return nil, err
		}
		if res != nil {
			return res, nil
		}
		res, err = tokenfactoryQuerier(ctx, request)
		return res, err
	}
}

func RegisterCustomPlugins(
	app *TerraApp,
) []wasmkeeper.Option {
	queryPluginOpt := wasmkeeper.WithQueryPlugins(&wasmkeeper.QueryPlugins{
		Custom: CustomQuerier(app),
	})
	messengerDecoratorOpt := wasmkeeper.WithMessageHandlerDecorator(
		tokenfactorybindings.CustomMessageDecorator(&app.BankKeeper.BaseKeeper, &app.TokenFactoryKeeper),
	)

	return []wasm.Option{
		queryPluginOpt,
		messengerDecoratorOpt,
	}
}

// GetWasmOpts build wasm options
func GetWasmOpts(app *TerraApp, appOpts servertypes.AppOptions) []wasm.Option {
	var wasmOpts []wasm.Option
	if cast.ToBool(appOpts.Get("telemetry.enabled")) {
		wasmOpts = append(wasmOpts, wasmkeeper.WithVMCacheMetrics(prometheus.DefaultRegisterer))
	}
	wasmOpts = append(wasmOpts, RegisterCustomPlugins(app)...)

	return wasmOpts
}
