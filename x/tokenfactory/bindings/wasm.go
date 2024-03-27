package bindings

import (
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	alliancebindings "github.com/terra-money/alliance/x/alliance/bindings"
	alliancekeeper "github.com/terra-money/alliance/x/alliance/keeper"
	tokenfactorykeeper "github.com/terra-money/core/v2/x/tokenfactory/keeper"
	wasm2 "github.com/terra-money/core/v2/x/wasm"

	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
)

func RegisterCustomPlugins(
	bank *bankkeeper.BaseKeeper,
	tokenFactory *tokenfactorykeeper.Keeper,
	allianceKeeper *alliancekeeper.Keeper,
) []wasmkeeper.Option {
	tfQuerier := CustomQuerier(NewQueryPlugin(bank, tokenFactory))
	allianceQuerier := alliancebindings.CustomQuerier(alliancebindings.NewAllianceQueryPlugin(allianceKeeper))
	queriers := wasm2.CustomQueriers(tfQuerier, allianceQuerier)

	queryPluginOpt := wasmkeeper.WithQueryPlugins(&wasmkeeper.QueryPlugins{
		Custom: queriers,
	})
	messengerDecoratorOpt := wasmkeeper.WithMessageHandlerDecorator(
		CustomMessageDecorator(bank, tokenFactory),
	)

	return []wasm.Option{
		queryPluginOpt,
		messengerDecoratorOpt,
	}
}
