package keepers

import (

	// #nosec G702

	"path/filepath"

	ibctransfer "github.com/cosmos/ibc-go/v7/modules/apps/transfer"
	ibcclient "github.com/cosmos/ibc-go/v7/modules/core/02-client"
	ibcclienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	"github.com/spf13/cast"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	evidencekeeper "github.com/cosmos/cosmos-sdk/x/evidence/keeper"
	feegrantkeeper "github.com/cosmos/cosmos-sdk/x/feegrant/keeper"
	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	"github.com/cosmos/cosmos-sdk/x/upgrade"

	"github.com/cosmos/cosmos-sdk/x/feegrant"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"

	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"

	consensusparamkeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	consensusparamtypes "github.com/cosmos/cosmos-sdk/x/consensus/types"

	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"

	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradekeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	"github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v7/packetforward"
	packetforwardkeeper "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v7/packetforward/keeper"
	packetforwardtypes "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v7/packetforward/types"

	icq "github.com/cosmos/ibc-apps/modules/async-icq/v7"
	icacontrollertypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/types"
	icahostkeeper "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/host/keeper"
	icahosttypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/host/types"
	icatypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/types"
	ibcfeekeeper "github.com/cosmos/ibc-go/v7/modules/apps/29-fee/keeper"
	ibcfeetypes "github.com/cosmos/ibc-go/v7/modules/apps/29-fee/types"
	ibctransferkeeper "github.com/cosmos/ibc-go/v7/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	porttypes "github.com/cosmos/ibc-go/v7/modules/core/05-port/types"
	ibcexported "github.com/cosmos/ibc-go/v7/modules/core/exported"

	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"

	icqkeeper "github.com/cosmos/ibc-apps/modules/async-icq/v7/keeper"
	icqtypes "github.com/cosmos/ibc-apps/modules/async-icq/v7/types"

	ibcfee "github.com/cosmos/ibc-go/v7/modules/apps/29-fee"
	ibckeeper "github.com/cosmos/ibc-go/v7/modules/core/keeper"

	ibchooks "github.com/cosmos/ibc-apps/modules/ibc-hooks/v7"
	ibchookskeeper "github.com/cosmos/ibc-apps/modules/ibc-hooks/v7/keeper"
	ibchookstypes "github.com/cosmos/ibc-apps/modules/ibc-hooks/v7/types"

	"github.com/CosmWasm/wasmd/x/wasm"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	icahost "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/host"
	customwasmkeeper "github.com/terra-money/core/v2/x/wasm/keeper"

	icacontroller "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller"
	icacontrollerkeeper "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/keeper"
	tokenfactorykeeper "github.com/terra-money/core/v2/x/tokenfactory/keeper"
	tokenfactorytypes "github.com/terra-money/core/v2/x/tokenfactory/types"

	"github.com/terra-money/alliance/x/alliance"
	alliancekeeper "github.com/terra-money/alliance/x/alliance/keeper"
	alliancetypes "github.com/terra-money/alliance/x/alliance/types"
	custombankkeeper "github.com/terra-money/core/v2/x/bank/keeper"
	feesharekeeper "github.com/terra-money/core/v2/x/feeshare/keeper"
	feesharetypes "github.com/terra-money/core/v2/x/feeshare/types"

	terraappconfig "github.com/terra-money/core/v2/app/config"
	// unnamed import of statik for swagger UI support
	_ "github.com/terra-money/core/v2/client/docs/statik"
)

// module account permissions
var maccPerms = map[string][]string{
	authtypes.FeeCollectorName:     nil,
	distrtypes.ModuleName:          nil,
	icatypes.ModuleName:            nil,
	minttypes.ModuleName:           {authtypes.Minter},
	stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
	stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
	govtypes.ModuleName:            {authtypes.Burner},
	ibctransfertypes.ModuleName:    {authtypes.Minter, authtypes.Burner},
	ibcfeetypes.ModuleName:         nil,
	icqtypes.ModuleName:            nil,
	wasmtypes.ModuleName:           {authtypes.Burner},
	tokenfactorytypes.ModuleName:   {authtypes.Burner, authtypes.Minter},
	alliancetypes.ModuleName:       {authtypes.Burner, authtypes.Minter},
	alliancetypes.RewardsPoolName:  nil,
}

type TerraAppKeepers struct {
	// Stores Keys
	keys    map[string]*storetypes.KVStoreKey
	tkeys   map[string]*storetypes.TransientStoreKey
	memKeys map[string]*storetypes.MemoryStoreKey

	// keepers
	AccountKeeper         authkeeper.AccountKeeper
	BankKeeper            custombankkeeper.Keeper
	CapabilityKeeper      *capabilitykeeper.Keeper
	StakingKeeper         *stakingkeeper.Keeper
	SlashingKeeper        slashingkeeper.Keeper
	MintKeeper            mintkeeper.Keeper
	DistrKeeper           distrkeeper.Keeper
	GovKeeper             govkeeper.Keeper
	CrisisKeeper          crisiskeeper.Keeper
	UpgradeKeeper         *upgradekeeper.Keeper
	ParamsKeeper          paramskeeper.Keeper
	ConsensusParamsKeeper consensusparamkeeper.Keeper
	IBCKeeper             *ibckeeper.Keeper // IBC Keeper must be a pointer in the app, so we can SetRouter on it correctly
	EvidenceKeeper        evidencekeeper.Keeper
	TransferKeeper        ibctransferkeeper.Keeper
	AuthzKeeper           authzkeeper.Keeper
	FeeGrantKeeper        feegrantkeeper.Keeper
	ICAControllerKeeper   icacontrollerkeeper.Keeper
	ICAHostKeeper         icahostkeeper.Keeper
	IBCFeeKeeper          ibcfeekeeper.Keeper
	PacketForwardKeeper   packetforwardkeeper.Keeper
	TokenFactoryKeeper    tokenfactorykeeper.Keeper
	AllianceKeeper        alliancekeeper.Keeper
	FeeShareKeeper        feesharekeeper.Keeper
	ICQKeeper             icqkeeper.Keeper

	// IBC hooks
	IBCHooksKeeper   *ibchookskeeper.Keeper
	TransferStack    porttypes.Middleware
	Ics20WasmHooks   *ibchooks.WasmHooks
	HooksICS4Wrapper ibchooks.ICS4Middleware

	// make scoped keepers public for test purposes
	ScopedIBCKeeper           capabilitykeeper.ScopedKeeper
	ScopedTransferKeeper      capabilitykeeper.ScopedKeeper
	ScopedICAControllerKeeper capabilitykeeper.ScopedKeeper
	ScopedICAHostKeeper       capabilitykeeper.ScopedKeeper
	ScopedICQKeeper           capabilitykeeper.ScopedKeeper

	WasmKeeper       customwasmkeeper.Keeper
	scopedWasmKeeper capabilitykeeper.ScopedKeeper
}

func NewTerraAppKeepers(
	appCodec codec.Codec,
	baseApp *baseapp.BaseApp,
	cdc *codec.LegacyAmino,
	appOpts servertypes.AppOptions,
	wasmOpts []wasmkeeper.Option,
	homePath string,
) (keepers TerraAppKeepers) {
	// Set keys KVStoreKey, TransientStoreKey, MemoryStoreKey
	keepers.GenerateKeys()
	keys := keepers.GetKVStoreKey()
	tkeys := keepers.GetTransientStoreKey()

	govModuleAddress := authtypes.NewModuleAddress(govtypes.ModuleName).String()

	keepers.ParamsKeeper = keepers.initParamsKeeper(
		appCodec,
		cdc,
		keys[paramstypes.StoreKey],
		tkeys[paramstypes.TStoreKey],
	)

	keepers.ConsensusParamsKeeper = consensusparamkeeper.NewKeeper(
		appCodec,
		keys[consensusparamtypes.StoreKey],
		govModuleAddress,
	)

	// Create capability keeper and grant capabilities for the modules
	keepers.CapabilityKeeper = capabilitykeeper.NewKeeper(
		appCodec,
		keys[capabilitytypes.StoreKey],
		keepers.memKeys[capabilitytypes.MemStoreKey],
	)
	keepers.ScopedIBCKeeper = keepers.CapabilityKeeper.ScopeToModule(ibcexported.ModuleName)
	keepers.ScopedTransferKeeper = keepers.CapabilityKeeper.ScopeToModule(ibctransfertypes.ModuleName)
	keepers.ScopedICAControllerKeeper = keepers.CapabilityKeeper.ScopeToModule(icacontrollertypes.SubModuleName)
	keepers.ScopedICAHostKeeper = keepers.CapabilityKeeper.ScopeToModule(icahosttypes.SubModuleName)
	keepers.scopedWasmKeeper = keepers.CapabilityKeeper.ScopeToModule(wasmtypes.ModuleName)
	keepers.ScopedICQKeeper = keepers.CapabilityKeeper.ScopeToModule(icqtypes.ModuleName)

	keepers.AccountKeeper = authkeeper.NewAccountKeeper(
		appCodec,
		keys[authtypes.StoreKey],
		authtypes.ProtoBaseAccount,
		maccPerms,
		terraappconfig.AccountAddressPrefix,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	keepers.BankKeeper = custombankkeeper.NewBaseKeeper(
		appCodec,
		keys[banktypes.StoreKey],
		keepers.AccountKeeper,
		keepers.ModuleAccountAddrs(),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	keepers.StakingKeeper = stakingkeeper.NewKeeper(
		appCodec,
		keys[stakingtypes.StoreKey],
		keepers.AccountKeeper,
		keepers.BankKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	keepers.MintKeeper = mintkeeper.NewKeeper(
		appCodec,
		keys[minttypes.StoreKey],
		keepers.StakingKeeper,
		keepers.AccountKeeper,
		keepers.BankKeeper,
		authtypes.FeeCollectorName,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	keepers.DistrKeeper = distrkeeper.NewKeeper(
		appCodec,
		keys[distrtypes.StoreKey],
		keepers.AccountKeeper,
		keepers.BankKeeper,
		keepers.StakingKeeper,
		authtypes.FeeCollectorName,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	keepers.SlashingKeeper = slashingkeeper.NewKeeper(
		appCodec,
		cdc,
		keys[slashingtypes.StoreKey],
		keepers.StakingKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	keepers.CrisisKeeper = *crisiskeeper.NewKeeper(
		appCodec,
		keys[crisistypes.StoreKey],
		cast.ToUint(appOpts.Get(server.FlagInvCheckPeriod)),
		keepers.BankKeeper,
		authtypes.FeeCollectorName,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// Create upgrade keeper with skip upgrade height and homepath from appOpts
	skipUpgradeHeights := map[int64]bool{}
	for _, h := range cast.ToIntSlice(appOpts.Get(server.FlagUnsafeSkipUpgrades)) {
		skipUpgradeHeights[int64(h)] = true
	}
	keepers.UpgradeKeeper = upgradekeeper.NewKeeper(
		skipUpgradeHeights,
		keys[upgradetypes.StoreKey],
		appCodec,
		homePath,
		baseApp,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	keepers.AllianceKeeper = alliancekeeper.NewKeeper(
		appCodec,
		keys[alliancetypes.StoreKey],
		keepers.AccountKeeper,
		keepers.BankKeeper,
		keepers.StakingKeeper,
		keepers.DistrKeeper,
		authtypes.FeeCollectorName,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	keepers.BankKeeper.RegisterKeepers(keepers.AllianceKeeper, keepers.StakingKeeper)

	// NOTE: stakingKeeper above is passed by reference, so that it will contain these hooks
	keepers.StakingKeeper.SetHooks(
		stakingtypes.NewMultiStakingHooks(
			keepers.DistrKeeper.Hooks(),
			keepers.SlashingKeeper.Hooks(),
			keepers.AllianceKeeper.StakingHooks(),
		),
	)
	keepers.TokenFactoryKeeper = tokenfactorykeeper.NewKeeper(
		keys[tokenfactorytypes.StoreKey],
		keepers.AccountKeeper,
		keepers.BankKeeper,
		keepers.DistrKeeper,
		appCodec,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	keepers.IBCKeeper = ibckeeper.NewKeeper(
		appCodec,
		keys[ibcexported.StoreKey],
		keepers.GetSubspace(ibcexported.ModuleName),
		keepers.StakingKeeper,
		keepers.UpgradeKeeper,
		keepers.ScopedIBCKeeper,
	)
	keepers.FeeGrantKeeper = feegrantkeeper.NewKeeper(
		appCodec,
		keys[feegrant.StoreKey],
		keepers.AccountKeeper,
	)
	keepers.AuthzKeeper = authzkeeper.NewKeeper(
		keys[authzkeeper.StoreKey],
		appCodec,
		baseApp.MsgServiceRouter(),
		keepers.AccountKeeper,
	)

	// register the proposal types
	govRouter := govtypesv1beta1.NewRouter()
	govRouter.
		AddRoute(govtypes.RouterKey, govtypesv1beta1.ProposalHandler).
		AddRoute(paramproposal.RouterKey, params.NewParamChangeProposalHandler(keepers.ParamsKeeper)).
		AddRoute(upgradetypes.RouterKey, upgrade.NewSoftwareUpgradeProposalHandler(keepers.UpgradeKeeper)).
		AddRoute(ibcclienttypes.RouterKey, ibcclient.NewClientProposalHandler(keepers.IBCKeeper.ClientKeeper)).
		AddRoute(alliancetypes.RouterKey, alliance.NewAllianceProposalHandler(keepers.AllianceKeeper))

	// Configure the hooks keeper
	hooksKeeper := ibchookskeeper.NewKeeper(
		keys[ibchookstypes.StoreKey],
	)
	keepers.IBCHooksKeeper = &hooksKeeper
	wasmHooks := ibchooks.NewWasmHooks(&hooksKeeper, nil, terraappconfig.AccountAddressPrefix) // The contract keeper needs to be set later
	keepers.Ics20WasmHooks = &wasmHooks
	keepers.HooksICS4Wrapper = ibchooks.NewICS4Middleware(
		keepers.IBCKeeper.ChannelKeeper,
		keepers.Ics20WasmHooks,
	)

	keepers.TransferKeeper = ibctransferkeeper.NewKeeper(
		appCodec,
		keys[ibctransfertypes.StoreKey],
		keepers.GetSubspace(ibctransfertypes.ModuleName),
		keepers.HooksICS4Wrapper,
		keepers.IBCKeeper.ChannelKeeper,
		&keepers.IBCKeeper.PortKeeper,
		keepers.AccountKeeper,
		keepers.BankKeeper,
		keepers.ScopedTransferKeeper,
	)
	transferIBCModule := ibctransfer.NewIBCModule(keepers.TransferKeeper)

	hooksTransferStack := ibchooks.NewIBCMiddleware(&transferIBCModule, &keepers.HooksICS4Wrapper)
	keepers.PacketForwardKeeper = *packetforwardkeeper.NewKeeper(
		appCodec,
		keepers.keys[packetforwardtypes.StoreKey],
		keepers.TransferKeeper,
		keepers.IBCKeeper.ChannelKeeper,
		keepers.DistrKeeper,
		keepers.BankKeeper,
		keepers.IBCKeeper.ChannelKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	keepers.TransferStack = packetforward.NewIBCMiddleware(
		hooksTransferStack,
		&keepers.PacketForwardKeeper,
		5,
		packetforwardkeeper.DefaultForwardTransferPacketTimeoutTimestamp,
		packetforwardkeeper.DefaultRefundTransferPacketTimeoutTimestamp,
	)
	keepers.ICQKeeper = icqkeeper.NewKeeper(
		appCodec,
		keepers.keys[icqtypes.StoreKey],
		keepers.GetSubspace(icqtypes.ModuleName),
		keepers.IBCKeeper.ChannelKeeper,
		keepers.IBCKeeper.ChannelKeeper,
		&keepers.IBCKeeper.PortKeeper,
		keepers.ScopedICQKeeper,
		baseApp.GRPCQueryRouter(),
	)
	keepers.IBCFeeKeeper = ibcfeekeeper.NewKeeper(
		appCodec,
		keys[ibcfeetypes.StoreKey],
		keepers.IBCKeeper.ChannelKeeper,
		keepers.IBCKeeper.ChannelKeeper,
		&keepers.IBCKeeper.PortKeeper,
		keepers.AccountKeeper,
		keepers.BankKeeper,
	)
	keepers.ICAControllerKeeper = icacontrollerkeeper.NewKeeper(
		appCodec,
		keys[icacontrollertypes.StoreKey],
		keepers.GetSubspace(icacontrollertypes.SubModuleName),
		keepers.IBCFeeKeeper,
		keepers.IBCKeeper.ChannelKeeper, &keepers.IBCKeeper.PortKeeper,
		keepers.ScopedICAControllerKeeper,
		baseApp.MsgServiceRouter(),
	)
	keepers.ICAHostKeeper = icahostkeeper.NewKeeper(
		appCodec, keys[icahosttypes.StoreKey],
		keepers.GetSubspace(icahosttypes.SubModuleName),
		keepers.IBCFeeKeeper,
		keepers.IBCKeeper.ChannelKeeper,
		&keepers.IBCKeeper.PortKeeper,
		keepers.AccountKeeper,
		keepers.ScopedICAHostKeeper,
		baseApp.MsgServiceRouter(),
	)

	var icaControllerStack porttypes.IBCModule
	icaControllerStack = icacontroller.NewIBCMiddleware(icaControllerStack, keepers.ICAControllerKeeper)
	icaControllerStack = ibcfee.NewIBCMiddleware(icaControllerStack, keepers.IBCFeeKeeper)

	icaHostIBCModule := icahost.NewIBCModule(keepers.ICAHostKeeper)
	icaHostStack := ibcfee.NewIBCMiddleware(icaHostIBCModule, keepers.IBCFeeKeeper)

	// Create evidence Keeper for to register the IBC light client misbehaviour evidence route
	evidenceKeeper := evidencekeeper.NewKeeper(
		appCodec, keys[evidencetypes.StoreKey], keepers.StakingKeeper, keepers.SlashingKeeper,
	)
	// If evidence needs to be handled for the app, set routes in router here and seal
	keepers.EvidenceKeeper = *evidenceKeeper

	wasmDir := filepath.Join(homePath, "data")

	// The last arguments can contain custom message handlers, and custom query handlers,
	// if we want to allow any custom callbacks
	wasmConfig, err := wasm.ReadWasmConfig(appOpts)
	if err != nil {
		panic("error while reading wasm config: " + err.Error())
	}
	availableCapabilities := "iterator,staking,stargate,cosmwasm_1_1,cosmwasm_1_2,cosmwasm_1_3,cosmwasm_1_4,cosmwasm_1_5,token_factory"
	keepers.WasmKeeper = customwasmkeeper.NewKeeper(
		appCodec,
		keys[wasmtypes.StoreKey],
		keepers.AccountKeeper,
		keepers.BankKeeper,
		keepers.StakingKeeper,
		distrkeeper.NewQuerier(keepers.DistrKeeper),
		keepers.IBCFeeKeeper, // ISC4 Wrapper: fee IBC middleware
		keepers.IBCKeeper.ChannelKeeper,
		&keepers.IBCKeeper.PortKeeper,
		keepers.scopedWasmKeeper,
		keepers.TransferKeeper,
		baseApp.MsgServiceRouter(),
		baseApp.GRPCQueryRouter(),
		wasmDir,
		wasmConfig,
		availableCapabilities,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		wasmOpts...,
	)

	keepers.Ics20WasmHooks.ContractKeeper = keepers.WasmKeeper.Keeper
	// Setup the contract keepers.WasmKeeper before the
	// hook for the BankKeeper othrwise the WasmKeeper
	// will be nil inside the hooks.
	keepers.TokenFactoryKeeper.SetContractKeeper(keepers.WasmKeeper)
	keepers.BankKeeper.SetHooks(
		custombankkeeper.NewMultiBankHooks(
			keepers.TokenFactoryKeeper.Hooks(),
		),
	)

	// Create fee enabled wasm ibc Stack
	var wasmStack porttypes.IBCModule
	wasmStack = wasm.NewIBCHandler(keepers.WasmKeeper, keepers.IBCKeeper.ChannelKeeper, keepers.IBCFeeKeeper)
	wasmStack = ibcfee.NewIBCMiddleware(wasmStack, keepers.IBCFeeKeeper)

	icqModule := icq.NewIBCModule(keepers.ICQKeeper)
	// Create static IBC router, add transfer route, then set and seal it
	ibcRouter := porttypes.NewRouter().
		AddRoute(icacontrollertypes.SubModuleName, icaControllerStack).
		AddRoute(icahosttypes.SubModuleName, icaHostStack).
		AddRoute(ibctransfertypes.ModuleName, keepers.TransferStack).
		AddRoute(wasmtypes.ModuleName, wasmStack).
		AddRoute(icqtypes.ModuleName, icqModule)

	keepers.IBCKeeper.SetRouter(ibcRouter)

	govKeeper := govkeeper.NewKeeper(
		appCodec,
		keys[govtypes.StoreKey],
		keepers.AccountKeeper,
		keepers.BankKeeper,
		keepers.StakingKeeper,
		baseApp.MsgServiceRouter(),
		govtypes.Config{
			MaxMetadataLen: 5100, // define the length of the governance proposal's title and description
		},
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	keepers.FeeShareKeeper = feesharekeeper.NewKeeper(
		appCodec,
		keys[feesharetypes.StoreKey],
		keepers.BankKeeper,
		keepers.WasmKeeper,
		keepers.AccountKeeper,
		authtypes.FeeCollectorName,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// Set legacy router for backwards compatibility with gov v1beta1
	govKeeper.SetLegacyRouter(govRouter)
	keepers.GovKeeper = *govKeeper.SetHooks(
		govtypes.NewMultiGovHooks(
		// register the governance hooks
		),
	)

	return keepers
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *TerraAppKeepers) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)

	/* #nosec */
	for acc := range maccPerms {
		modAccAddrs[authtypes.NewModuleAddress(acc).String()] = true
	}

	delete(modAccAddrs, authtypes.NewModuleAddress(alliancetypes.ModuleName).String())

	return modAccAddrs
}

// initParamsKeeper init params keeper and its subspaces
func (app *TerraAppKeepers) initParamsKeeper(appCodec codec.BinaryCodec, legacyAmino *codec.LegacyAmino, key, tkey storetypes.StoreKey) paramskeeper.Keeper {
	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)

	// Cosmo SDK Base Modules
	paramsKeeper.Subspace(authtypes.ModuleName).WithKeyTable(authtypes.ParamKeyTable())
	paramsKeeper.Subspace(banktypes.ModuleName).WithKeyTable(banktypes.ParamKeyTable())
	paramsKeeper.Subspace(stakingtypes.ModuleName).WithKeyTable(stakingtypes.ParamKeyTable())
	paramsKeeper.Subspace(minttypes.ModuleName).WithKeyTable(minttypes.ParamKeyTable())
	paramsKeeper.Subspace(distrtypes.ModuleName).WithKeyTable(distrtypes.ParamKeyTable())
	paramsKeeper.Subspace(slashingtypes.ModuleName).WithKeyTable(slashingtypes.ParamKeyTable())
	paramsKeeper.Subspace(govtypes.ModuleName).WithKeyTable(govtypesv1.ParamKeyTable())
	paramsKeeper.Subspace(crisistypes.ModuleName).WithKeyTable(crisistypes.ParamKeyTable())

	// IBC Modules
	paramsKeeper.Subspace(ibctransfertypes.ModuleName)
	paramsKeeper.Subspace(ibcexported.ModuleName)
	paramsKeeper.Subspace(icahosttypes.SubModuleName)
	paramsKeeper.Subspace(icacontrollertypes.SubModuleName)
	paramsKeeper.Subspace(packetforwardtypes.ModuleName).WithKeyTable(packetforwardtypes.ParamKeyTable())
	paramsKeeper.Subspace(icqtypes.ModuleName)

	// Custom Modules
	paramsKeeper.Subspace(wasmtypes.ModuleName).WithKeyTable(wasmtypes.ParamKeyTable())
	paramsKeeper.Subspace(tokenfactorytypes.ModuleName).WithKeyTable(tokenfactorytypes.ParamKeyTable())
	paramsKeeper.Subspace(feesharetypes.ModuleName).WithKeyTable(feesharetypes.ParamKeyTable())
	paramsKeeper.Subspace(alliancetypes.ModuleName).WithKeyTable(alliancetypes.ParamKeyTable())

	return paramsKeeper
}

// GetSubspace returns a param subspace for a given module name.
func (appKeepers *TerraAppKeepers) GetSubspace(moduleName string) paramstypes.Subspace {
	subspace, _ := appKeepers.ParamsKeeper.GetSubspace(moduleName)
	return subspace
}

// GetMaccPerms returns a copy of the module account permissions
func ModuleAccountAddrs() map[string]bool {
	dupMaccPerms := make(map[string]bool)
	for acc := range maccPerms {
		dupMaccPerms[authtypes.NewModuleAddress(acc).String()] = true
	}
	delete(dupMaccPerms, authtypes.NewModuleAddress(alliancetypes.ModuleName).String())

	return dupMaccPerms
}
