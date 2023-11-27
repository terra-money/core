package app

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"reflect" // #nosec G702

	"github.com/prometheus/client_golang/prometheus"

	authsims "github.com/cosmos/cosmos-sdk/x/auth/simulation"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/terra-money/core/v2/app/keepers"
	"github.com/terra-money/core/v2/app/post"
	"github.com/terra-money/core/v2/app/rpc"
	tokenfactorybindings "github.com/terra-money/core/v2/x/tokenfactory/bindings"

	dbm "github.com/cometbft/cometbft-db"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/libs/log"
	tmos "github.com/cometbft/cometbft/libs/os"
	"github.com/gorilla/mux"
	"github.com/rakyll/statik/fs"
	"github.com/spf13/cast"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec/types"

	custombankmodule "github.com/terra-money/core/v2/x/bank"
	customwasmodule "github.com/terra-money/core/v2/x/wasm"

	"github.com/cosmos/cosmos-sdk/baseapp"
	nodeservice "github.com/cosmos/cosmos-sdk/client/grpc/node"
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	cosmosante "github.com/cosmos/cosmos-sdk/x/auth/ante"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	vestingexported "github.com/cosmos/cosmos-sdk/x/auth/vesting/exported"
	authzmodule "github.com/cosmos/cosmos-sdk/x/authz/module"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/capability"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	feegrantmodule "github.com/cosmos/cosmos-sdk/x/feegrant/module"

	"github.com/cosmos/cosmos-sdk/x/gov"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/cosmos/cosmos-sdk/x/mint"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"

	"github.com/cosmos/cosmos-sdk/x/params"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradeclient "github.com/cosmos/cosmos-sdk/x/upgrade/client"

	"github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v7/router"

	ica "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts"
	ibcfee "github.com/cosmos/ibc-go/v7/modules/apps/29-fee"
	ibctransfer "github.com/cosmos/ibc-go/v7/modules/apps/transfer"
	ibc "github.com/cosmos/ibc-go/v7/modules/core"
	ibcclientclient "github.com/cosmos/ibc-go/v7/modules/core/02-client/client"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"

	"github.com/terra-money/alliance/x/alliance"
	allianceclient "github.com/terra-money/alliance/x/alliance/client"
	alliancetypes "github.com/terra-money/alliance/x/alliance/types"
	feeshare "github.com/terra-money/core/v2/x/feeshare"
	feesharetypes "github.com/terra-money/core/v2/x/feeshare/types"

	pobabci "github.com/skip-mev/pob/abci"
	pobmempool "github.com/skip-mev/pob/mempool"

	tmjson "github.com/cometbft/cometbft/libs/json"

	"github.com/terra-money/core/v2/app/ante"
	terraappconfig "github.com/terra-money/core/v2/app/config"
	terraappparams "github.com/terra-money/core/v2/app/params"

	// unnamed import of statik for swagger UI support
	_ "github.com/terra-money/core/v2/client/docs/statik"
)

func getGovProposalHandlers() []govclient.ProposalHandler {
	var govProposalHandlers []govclient.ProposalHandler

	govProposalHandlers = append(govProposalHandlers,
		paramsclient.ProposalHandler,
		upgradeclient.LegacyProposalHandler,
		upgradeclient.LegacyCancelProposalHandler,
		ibcclientclient.UpdateClientProposalHandler,
		ibcclientclient.UpgradeProposalHandler,
		allianceclient.CreateAllianceProposalHandler,
		allianceclient.UpdateAllianceProposalHandler,
		allianceclient.DeleteAllianceProposalHandler,
	)

	return govProposalHandlers
}

var (
	// DefaultNodeHome default home directories for the application daemon
	DefaultNodeHome string
)

var (
	_ servertypes.Application = (*TerraApp)(nil)
)

func init() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	DefaultNodeHome = filepath.Join(userHomeDir, "."+terraappconfig.AppName)
}

// TerraApp extends an ABCI application, but with most of its parameters exported.
// They are exported for convenience in creating helper functions, as object
// capabilities aren't needed for testing.
type TerraApp struct {
	*baseapp.BaseApp

	cdc               *codec.LegacyAmino
	appCodec          codec.Codec
	interfaceRegistry types.InterfaceRegistry

	// keys to access the substores
	keys    map[string]*storetypes.KVStoreKey
	tkeys   map[string]*storetypes.TransientStoreKey
	memKeys map[string]*storetypes.MemoryStoreKey

	Keepers keepers.TerraAppKeepers

	invCheckPeriod uint

	// Custom checkTx handler
	checkTxHandler pobabci.CheckTx

	// the module manager
	mm           *module.Manager
	basicManager module.BasicManager
	// the configurator
	configurator module.Configurator
}

func (app TerraApp) GetAppCodec() codec.Codec {
	return app.appCodec
}

func (app TerraApp) GetConfigurator() module.Configurator {
	return app.configurator
}

func (app TerraApp) GetModuleManager() *module.Manager {
	return app.mm
}

// NewTerraApp returns a reference to an initialized Terra.
func NewTerraApp(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	loadLatest bool,
	skipUpgradeHeights map[int64]bool,
	homePath string,
	invCheckPeriod uint,
	encodingConfig terraappparams.EncodingConfig,
	appOpts servertypes.AppOptions,
	wasmConfig wasmtypes.WasmConfig,
	baseAppOptions ...func(*baseapp.BaseApp),
) *TerraApp {
	appCodec := encodingConfig.Marshaler
	cdc := encodingConfig.Amino
	interfaceRegistry := encodingConfig.InterfaceRegistry

	bApp := baseapp.NewBaseApp(terraappconfig.AppName, logger, db, encodingConfig.TxConfig.TxDecoder(), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetVersion(version.Version)
	bApp.SetInterfaceRegistry(interfaceRegistry)
	app := &TerraApp{
		BaseApp:           bApp,
		cdc:               cdc,
		appCodec:          appCodec,
		interfaceRegistry: interfaceRegistry,
		invCheckPeriod:    invCheckPeriod,
	}
	app.Keepers = keepers.NewTerraAppKeepers(
		appCodec,
		bApp,
		cdc,
		appOpts,
		app.GetWasmOpts(appOpts),
	)
	app.keys = app.Keepers.GetKVStoreKey()
	app.tkeys = app.Keepers.GetTransientStoreKey()
	app.memKeys = app.Keepers.GetMemoryStoreKey()
	bApp.SetParamStore(&app.Keepers.ConsensusParamsKeeper)

	// upgrade handlers
	app.configurator = module.NewConfigurator(app.appCodec, app.MsgServiceRouter(), app.GRPCQueryRouter())

	/****  Module Options ****/

	// NOTE: we may consider parsing `appOpts` inside module constructors. For the moment
	// we prefer to be more strict in what arguments the modules expect.
	skipGenesisInvariants := cast.ToBool(appOpts.Get(crisis.FlagSkipGenesisInvariants))

	app.mm = module.NewManager(appModules(app, encodingConfig, skipGenesisInvariants)...)

	// NOTE: Any module instantiated in the module manager that is later modified
	// must be passed by reference here.
	app.mm.SetOrderBeginBlockers(beginBlockersOrder...)

	app.mm.SetOrderEndBlockers(endBlockerOrder...)

	// NOTE: The genutils module must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	// NOTE: Capability module must occur first so that it can initialize any capabilities
	// so that other modules that want to create or claim capabilities afterwards in InitChain
	// can do so safely.
	app.mm.SetOrderInitGenesis(initGenesisOrder...)

	app.mm.RegisterInvariants(&app.Keepers.CrisisKeeper)
	app.mm.RegisterServices(app.configurator)

	// initialize stores
	app.MountKVStores(app.keys)
	app.MountTransientStores(app.tkeys)
	app.MountMemoryStores(app.memKeys)

	// register upgrade
	app.RegisterUpgradeHandlers()
	app.RegisterUpgradeStores()

	config := pobmempool.NewDefaultAuctionFactory(encodingConfig.TxConfig.TxDecoder())
	// when maxTx is set as 0, there won't be a limit on the number of txs in this mempool
	pobMempool := pobmempool.NewAuctionMempool(encodingConfig.TxConfig.TxDecoder(), encodingConfig.TxConfig.TxEncoder(), 0, config)

	anteHandler, err := ante.NewAnteHandler(
		ante.HandlerOptions{
			HandlerOptions: cosmosante.HandlerOptions{
				AccountKeeper:   app.Keepers.AccountKeeper,
				BankKeeper:      app.Keepers.BankKeeper,
				FeegrantKeeper:  app.Keepers.FeeGrantKeeper,
				SignModeHandler: encodingConfig.TxConfig.SignModeHandler(),
				SigGasConsumer:  cosmosante.DefaultSigVerificationGasConsumer,
			},
			BankKeeper:        app.Keepers.BankKeeper,
			FeeShareKeeper:    app.Keepers.FeeShareKeeper,
			IBCkeeper:         app.Keepers.IBCKeeper,
			TxCounterStoreKey: app.keys[wasmtypes.StoreKey],
			WasmConfig:        wasmConfig,
			PobBuilderKeeper:  app.Keepers.BuilderKeeper,
			TxConfig:          encodingConfig.TxConfig,
			PobMempool:        pobMempool,
		},
	)
	if err != nil {
		panic(err)
	}
	postHandler := post.NewPostHandler(
		post.HandlerOptions{
			FeeShareKeeper: app.Keepers.FeeShareKeeper,
			BankKeeper:     app.Keepers.BankKeeper,
			WasmKeeper:     app.Keepers.WasmKeeper,
		},
	)

	// Create the proposal handler that will be used to build and validate blocks.
	handler := pobabci.NewProposalHandler(
		pobMempool,
		bApp.Logger(),
		anteHandler,
		encodingConfig.TxConfig.TxEncoder(),
		encodingConfig.TxConfig.TxDecoder(),
	)
	app.SetPrepareProposal(handler.PrepareProposalHandler())
	app.SetProcessProposal(handler.ProcessProposalHandler())

	// Set the custom CheckTx handler on BaseApp.
	checkTxHandler := pobabci.NewCheckTxHandler(
		app.BaseApp,
		encodingConfig.TxConfig.TxDecoder(),
		pobMempool,
		anteHandler,
		app.ChainID(),
	)

	// initialize BaseApp
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetAnteHandler(anteHandler)
	app.SetPostHandler(postHandler)
	app.SetEndBlocker(app.EndBlocker)
	app.SetMempool(pobMempool)
	app.SetCheckTx(checkTxHandler.CheckTx())

	if loadLatest {
		if err := app.LoadLatestVersion(); err != nil {
			tmos.Exit(err.Error())
		}

		// Initialize and seal the capability keeper so all persistent capabilities
		// are loaded in-memory and prevent any further modules from creating scoped
		// sub-keepers.
		// This must be done during creation of baseapp rather than in InitChain so
		// that in-memory capabilities get regenerated on app restart.
		// Note that since this reads from the store, we can only perform it when
		// `loadLatest` is set to true.
		app.Keepers.CapabilityKeeper.Seal()
	}

	return app
}

// Name returns the name of the App
func (app *TerraApp) Name() string { return app.BaseApp.Name() }

// BeginBlocker application updates every begin block
func (app *TerraApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

// EndBlocker application updates every end block
func (app *TerraApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

// InitChainer application update at chain initialization
func (app *TerraApp) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState GenesisState
	if err := tmjson.Unmarshal(req.AppStateBytes, &genesisState); err != nil {
		panic(err)
	}
	res := app.mm.InitGenesis(ctx, app.appCodec, genesisState)
	app.Keepers.UpgradeKeeper.SetModuleVersionMap(ctx, app.mm.GetVersionMap())

	// stake all vesting tokens
	app.enforceStakingForVestingTokens(ctx, genesisState)

	return res
}

// LoadHeight loads a particular height
func (app *TerraApp) LoadHeight(height int64) error {
	return app.LoadVersion(height)
}

// LegacyAmino returns SimApp's amino codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *TerraApp) LegacyAmino() *codec.LegacyAmino {
	return app.cdc
}

// AppCodec returns Terra's app codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *TerraApp) AppCodec() codec.Codec {
	return app.appCodec
}

// InterfaceRegistry returns Terra's InterfaceRegistry
func (app *TerraApp) InterfaceRegistry() types.InterfaceRegistry {
	return app.interfaceRegistry
}

// GetKey returns the KVStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *TerraApp) GetKey(storeKey string) *storetypes.KVStoreKey {
	return app.keys[storeKey]
}

// GetTKey returns the TransientStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *TerraApp) GetTKey(storeKey string) *storetypes.TransientStoreKey {
	return app.tkeys[storeKey]
}

// GetMemKey returns the MemStoreKey for the provided mem key.
//
// NOTE: This is solely used for testing purposes.
func (app *TerraApp) GetMemKey(storeKey string) *storetypes.MemoryStoreKey {
	return app.memKeys[storeKey]
}

// GetSubspace returns a param subspace for a given module name.
//
// NOTE: This is solely to be used for testing purposes.
func (app *TerraApp) GetSubspace(moduleName string) paramstypes.Subspace {
	subspace, found := app.Keepers.ParamsKeeper.GetSubspace(moduleName)
	if !found {
		panic("Module with '" + moduleName + "' name does not exist")
	}

	return subspace
}

// RegisterAPIRoutes registers all application module routes with the provided
// API server.
func (app *TerraApp) RegisterAPIRoutes(apiSvr *api.Server, apiConfig config.APIConfig) {
	clientCtx := apiSvr.ClientCtx

	rpc.RegisterHealthcheckRoute(clientCtx, apiSvr.Router)
	// Register new tx routes from grpc-gateway.
	authtx.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	// Register new tendermint queries routes from grpc-gateway.
	tmservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	// Register node gRPC service for grpc-gateway.
	nodeservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	// Register legacy and grpc-gateway routes for all modules.
	ModuleBasics.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// register swagger API from root so that other applications can override easily
	if apiConfig.Swagger {
		RegisterSwaggerAPI(apiSvr.Router)
	}
}

// RegisterTxService implements the Application.RegisterTxService method.
func (app *TerraApp) RegisterTxService(clientCtx client.Context) {
	authtx.RegisterTxService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.BaseApp.Simulate, app.interfaceRegistry)
}

// RegisterTendermintService implements the Application.RegisterTendermintService method.
func (app *TerraApp) RegisterTendermintService(clientCtx client.Context) {
	tmservice.RegisterTendermintService(
		clientCtx,
		app.BaseApp.GRPCQueryRouter(),
		app.interfaceRegistry,
		app.Query,
	)
}

// RegisterSwaggerAPI registers swagger route with API Server
func RegisterSwaggerAPI(rtr *mux.Router) {
	statikFS, err := fs.New()
	if err != nil {
		panic(err)
	}

	staticServer := http.FileServer(statikFS)
	rtr.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger/", staticServer))
}

// enforceStakingForVestingTokens enforce vesting tokens to be staked
// CONTRACT: validator's gentx account must not be a vesting account
func (app *TerraApp) enforceStakingForVestingTokens(ctx sdk.Context, genesisState GenesisState) {

	var authState authtypes.GenesisState
	app.appCodec.MustUnmarshalJSON(genesisState[authtypes.ModuleName], &authState)

	allValidators := app.Keepers.StakingKeeper.GetAllValidators(ctx)

	// Filter out validators which have huge max commission than 20%
	var validators []stakingtypes.Validator
	maxCommissionCondition := sdk.NewDecWithPrec(20, 2)
	for _, val := range allValidators {
		if val.Commission.CommissionRates.MaxRate.LTE(maxCommissionCondition) {
			validators = append(validators, val)
		}
	}

	validatorLen := len(validators)

	// ignore when validator len is zero
	if validatorLen == 0 {
		return
	}

	i := 0
	stakeSplitCondition := sdk.NewInt(1_000_000_000_000)
	powerReduction := app.Keepers.StakingKeeper.PowerReduction(ctx)
	for _, acc := range authState.GetAccounts() {
		var account authtypes.AccountI
		if err := app.InterfaceRegistry().UnpackAny(acc, &account); err != nil {
			panic(err)
		}

		if vestingAcc, ok := account.(vestingexported.VestingAccount); ok {
			amt := vestingAcc.GetOriginalVesting().AmountOf(app.Keepers.StakingKeeper.BondDenom(ctx))

			// to prevent staking multiple times over the same validator
			// adjust split amount for the whale account
			splitAmt := stakeSplitCondition
			if amt.GT(stakeSplitCondition.MulRaw(int64(validatorLen))) {
				splitAmt = amt.QuoRaw(int64(validatorLen))
			}

			// if a vesting account has more staking token than `stakeSplitCondition`,
			// split staking balance to distribute staking power evenly
			// Ex) 2_200_000_000_000
			// stake 1_000_000_000_000 to val1
			// stake 1_000_000_000_000 to val2
			// stake 200_000_000_000 to val3
			for ; amt.GTE(powerReduction); amt = amt.Sub(splitAmt) {
				validator := validators[i%validatorLen]
				if _, err := app.Keepers.StakingKeeper.Delegate(
					ctx,
					vestingAcc.GetAddress(),
					sdk.MinInt(amt, splitAmt),
					stakingtypes.Unbonded,
					validator,
					true,
				); err != nil {
					panic(err)
				}

				// reload validator to avoid power index problem
				validator, _ = app.Keepers.StakingKeeper.GetValidator(ctx, validator.GetOperator())
				validators[i%validatorLen] = validator

				// increase index only when staking happened
				i++
			}

		}
	}
}

func (app *TerraApp) SimulationManager() *module.SimulationManager {
	appCodec := app.appCodec
	// create the simulation manager and define the order of the modules for deterministic simulations
	sm := module.NewSimulationManager(
		auth.NewAppModule(appCodec, app.Keepers.AccountKeeper, authsims.RandomGenesisAccounts, app.Keepers.GetSubspace(authtypes.ModuleName)),
		authzmodule.NewAppModule(appCodec, app.Keepers.AuthzKeeper, app.Keepers.AccountKeeper, app.Keepers.BankKeeper, app.interfaceRegistry),
		custombankmodule.NewAppModule(appCodec, app.Keepers.BankKeeper, app.Keepers.AccountKeeper, app.Keepers.GetSubspace(banktypes.ModuleName)),
		capability.NewAppModule(appCodec, *app.Keepers.CapabilityKeeper, false),
		feegrantmodule.NewAppModule(appCodec, app.Keepers.AccountKeeper, app.Keepers.BankKeeper, app.Keepers.FeeGrantKeeper, app.interfaceRegistry),
		gov.NewAppModule(appCodec, &app.Keepers.GovKeeper, app.Keepers.AccountKeeper, app.Keepers.BankKeeper, app.Keepers.GetSubspace(govtypes.ModuleName)),
		mint.NewAppModule(appCodec, app.Keepers.MintKeeper, app.Keepers.AccountKeeper, nil, app.Keepers.GetSubspace(minttypes.ModuleName)),
		staking.NewAppModule(appCodec, app.Keepers.StakingKeeper, app.Keepers.AccountKeeper, app.Keepers.BankKeeper, app.Keepers.GetSubspace(stakingtypes.ModuleName)),
		distr.NewAppModule(appCodec, app.Keepers.DistrKeeper, app.Keepers.AccountKeeper, app.Keepers.BankKeeper, app.Keepers.StakingKeeper, app.Keepers.GetSubspace(distrtypes.ModuleName)),
		slashing.NewAppModule(appCodec, app.Keepers.SlashingKeeper, app.Keepers.AccountKeeper, app.Keepers.BankKeeper, app.Keepers.StakingKeeper, app.Keepers.GetSubspace(stakingtypes.ModuleName)),
		params.NewAppModule(app.Keepers.ParamsKeeper),
		evidence.NewAppModule(app.Keepers.EvidenceKeeper),
		ibc.NewAppModule(app.Keepers.IBCKeeper),
		ibctransfer.NewAppModule(app.Keepers.TransferKeeper),
		ibcfee.NewAppModule(app.Keepers.IBCFeeKeeper),
		ica.NewAppModule(&app.Keepers.ICAControllerKeeper, &app.Keepers.ICAHostKeeper),
		router.NewAppModule(&app.Keepers.RouterKeeper),
		customwasmodule.NewAppModule(appCodec, &app.Keepers.WasmKeeper, app.Keepers.StakingKeeper, app.Keepers.AccountKeeper, app.Keepers.BankKeeper, app.BaseApp.MsgServiceRouter(), app.Keepers.GetSubspace(wasmtypes.ModuleName)),
		alliance.NewAppModule(appCodec, app.Keepers.AllianceKeeper, app.Keepers.StakingKeeper, app.Keepers.AccountKeeper, app.Keepers.BankKeeper, app.interfaceRegistry, app.Keepers.GetSubspace(alliancetypes.ModuleName)),
		feeshare.NewAppModule(app.Keepers.FeeShareKeeper, app.Keepers.AccountKeeper, app.GetSubspace(feesharetypes.ModuleName)),
	)

	sm.RegisterStoreDecoders()
	return sm
}

// DefaultGenesis returns a default genesis from the registered AppModuleBasic's.
func (a *TerraApp) DefaultGenesis() map[string]json.RawMessage {
	return a.basicManager.DefaultGenesis(a.appCodec)
}

func (app *TerraApp) RegisterNodeService(clientCtx client.Context) {
	nodeservice.RegisterNodeService(clientCtx, app.GRPCQueryRouter())
}

// ChainID gets chainID from private fields of BaseApp
// Should be removed once SDK 0.50.x will be adopted
func (app *TerraApp) ChainID() string {
	field := reflect.ValueOf(app.BaseApp).Elem().FieldByName("chainID")
	return field.String()
}

// CheckTx will check the transaction with the provided checkTxHandler. We override the default
// handler so that we can verify bid transactions before they are inserted into the mempool.
// With the POB CheckTx, we can verify the bid transaction and all of the bundled transactions
// before inserting the bid transaction into the mempool.
func (app *TerraApp) CheckTx(req abci.RequestCheckTx) abci.ResponseCheckTx {
	return app.checkTxHandler(req)
}

// SetCheckTx sets the checkTxHandler for the app.
func (app *TerraApp) SetCheckTx(handler pobabci.CheckTx) {
	app.checkTxHandler = handler
}

func (app *TerraApp) GetWasmOpts(appOpts servertypes.AppOptions) []wasmkeeper.Option {
	var wasmOpts []wasmkeeper.Option
	if cast.ToBool(appOpts.Get("telemetry.enabled")) {
		wasmOpts = append(wasmOpts, wasmkeeper.WithVMCacheMetrics(prometheus.DefaultRegisterer))
	}

	wasmOpts = append(wasmOpts, tokenfactorybindings.RegisterCustomPlugins(
		&app.Keepers.BankKeeper.BaseKeeper,
		&app.Keepers.TokenFactoryKeeper)...,
	)

	return wasmOpts
}
