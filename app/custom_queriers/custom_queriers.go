package custom_queriers

import (
	"encoding/json"
	"fmt"
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	alliancebindings "github.com/terra-money/alliance/x/alliance/bindings"
	alliancekeeper "github.com/terra-money/alliance/x/alliance/keeper"
	tokenfactorybindings "github.com/terra-money/core/v2/x/tokenfactory/bindings"
	tokenfactorykeeper "github.com/terra-money/core/v2/x/tokenfactory/keeper"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Querier func(ctx sdk.Context, request json.RawMessage) ([]byte, error)

func CustomQueriers(queriers ...Querier) func(ctx sdk.Context, request json.RawMessage) ([]byte, error) {
	return func(ctx sdk.Context, request json.RawMessage) ([]byte, error) {
		for _, querier := range queriers {
			res, err := querier(ctx, request)
			if err == nil || !strings.Contains(err.Error(), "unknown query") {
				return res, err
			}
		}
		return nil, fmt.Errorf("unknown query")
	}
}

func RegisterCustomPlugins(
	bank *bankkeeper.BaseKeeper,
	tokenFactory *tokenfactorykeeper.Keeper,
	allianceKeeper *alliancekeeper.Keeper,
) []wasmkeeper.Option {
	tfQuerier := tokenfactorybindings.CustomQuerier(tokenfactorybindings.NewQueryPlugin(bank, tokenFactory))
	allianceQuerier := alliancebindings.CustomQuerier(alliancebindings.NewAllianceQueryPlugin(allianceKeeper))
	queriers := CustomQueriers(tfQuerier, allianceQuerier)

	queryPluginOpt := wasmkeeper.WithQueryPlugins(&wasmkeeper.QueryPlugins{
		Custom: queriers,
	})
	messengerDecoratorOpt := wasmkeeper.WithMessageHandlerDecorator(
		tokenfactorybindings.CustomMessageDecorator(bank, tokenFactory),
	)

	return []wasm.Option{
		queryPluginOpt,
		messengerDecoratorOpt,
	}
}
