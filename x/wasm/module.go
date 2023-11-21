package wasm

import (
	"github.com/cosmos/cosmos-sdk/baseapp"

	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/CosmWasm/wasmd/x/wasm"
	"github.com/CosmWasm/wasmd/x/wasm/exported"
	"github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/CosmWasm/wasmd/x/wasm/simulation"
	"github.com/CosmWasm/wasmd/x/wasm/types"

	customwasmkeeper "github.com/terra-money/core/v2/x/wasm/keeper"

	"github.com/cosmos/cosmos-sdk/types/module"
)

// AppModule implements an application module for the wasm module.
type AppModule struct {
	wasm.AppModule
	keeper         *customwasmkeeper.Keeper
	legacySubspace exported.Subspace
	msgServer      types.MsgServer
}

// NewAppModule creates a new AppModule object
func NewAppModule(
	cdc codec.Codec,
	keeper *customwasmkeeper.Keeper,
	validatorSetSource keeper.ValidatorSetSource,
	ak types.AccountKeeper,
	bk simulation.BankKeeper,
	router *baseapp.MsgServiceRouter,
	ss exported.Subspace,
) AppModule {
	appModule := wasm.NewAppModule(cdc, keeper.Keeper, validatorSetSource, ak, bk, router, ss)
	msgServer := customwasmkeeper.NewCustomMsgServerImpl(keeper)

	return AppModule{
		AppModule:      appModule,
		keeper:         keeper,
		legacySubspace: ss,
		msgServer:      msgServer,
	}
}

func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), am.msgServer)
	types.RegisterQueryServer(cfg.QueryServer(), keeper.Querier(am.keeper.Keeper))

	m := keeper.NewMigrator(*am.keeper.Keeper, am.legacySubspace)
	err := cfg.RegisterMigration(types.ModuleName, 1, m.Migrate1to2)
	if err != nil {
		panic(err)
	}
	err = cfg.RegisterMigration(types.ModuleName, 2, m.Migrate2to3)
	if err != nil {
		panic(err)
	}
	err = cfg.RegisterMigration(types.ModuleName, 3, m.Migrate3to4)
	if err != nil {
		panic(err)
	}
}
