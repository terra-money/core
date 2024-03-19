package app

import (

	// #nosec G702

	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	"github.com/cosmos/cosmos-sdk/x/authz"
	authzmodule "github.com/cosmos/cosmos-sdk/x/authz/module"
	"github.com/cosmos/cosmos-sdk/x/bank"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/capability"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	feegrantmodule "github.com/cosmos/cosmos-sdk/x/feegrant/module"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"

	"github.com/cosmos/cosmos-sdk/x/gov"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/cosmos/cosmos-sdk/x/mint"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"

	consensus "github.com/cosmos/cosmos-sdk/x/consensus"

	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v7/packetforward"
	packetforwardtypes "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v7/packetforward/types"
	icqtypes "github.com/cosmos/ibc-apps/modules/async-icq/v7/types"
	ibchookstypes "github.com/cosmos/ibc-apps/modules/ibc-hooks/v7/types"
	ica "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts"
	icatypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/types"
	ibcfee "github.com/cosmos/ibc-go/v7/modules/apps/29-fee"
	ibcfeetypes "github.com/cosmos/ibc-go/v7/modules/apps/29-fee/types"
	ibctransfer "github.com/cosmos/ibc-go/v7/modules/apps/transfer"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v7/modules/core"
	alliancetypes "github.com/terra-money/alliance/x/alliance/types"
	feesharetypes "github.com/terra-money/core/v2/x/feeshare/types"
	tokenfactorytypes "github.com/terra-money/core/v2/x/tokenfactory/types"

	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	consensusparamtypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	icq "github.com/cosmos/ibc-apps/modules/async-icq/v7"

	solomachine "github.com/cosmos/ibc-go/v7/modules/light-clients/06-solomachine"
	ibctm "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"

	ibchooks "github.com/cosmos/ibc-apps/modules/ibc-hooks/v7"

	"github.com/CosmWasm/wasmd/x/wasm"
	custombankmodule "github.com/terra-money/core/v2/x/bank"
	customwasmodule "github.com/terra-money/core/v2/x/wasm"

	"github.com/terra-money/core/v2/x/tokenfactory"

	"github.com/terra-money/alliance/x/alliance"
	feeshare "github.com/terra-money/core/v2/x/feeshare"

	ibcexported "github.com/cosmos/ibc-go/v7/modules/core/exported"
	terrappsparams "github.com/terra-money/core/v2/app/params"

	"github.com/cosmos/cosmos-sdk/x/feegrant"
)

// ModuleBasics defines the module BasicManager is in charge of setting up basic,
// non-dependant module elements, such as codec registration
// and genesis verification.
var ModuleBasics = module.NewBasicManager(
	auth.AppModuleBasic{},
	genutil.NewAppModuleBasic(genutiltypes.DefaultMessageValidator),
	bank.AppModuleBasic{},
	capability.AppModuleBasic{},
	staking.AppModuleBasic{},
	mint.AppModuleBasic{},
	distr.AppModuleBasic{},
	gov.NewAppModuleBasic(getGovProposalHandlers()),
	params.AppModuleBasic{},
	crisis.AppModuleBasic{},
	slashing.AppModuleBasic{},
	ibc.AppModuleBasic{},
	ibctm.AppModuleBasic{},
	solomachine.AppModuleBasic{},
	feegrantmodule.AppModuleBasic{},
	upgrade.AppModuleBasic{},
	evidence.AppModuleBasic{},
	ibctransfer.AppModuleBasic{},
	vesting.AppModuleBasic{},
	ica.AppModuleBasic{},
	ibcfee.AppModuleBasic{},
	packetforward.AppModuleBasic{},
	authzmodule.AppModuleBasic{},
	tokenfactory.AppModuleBasic{},
	ibchooks.AppModuleBasic{},
	wasm.AppModuleBasic{},
	consensus.AppModuleBasic{},
	alliance.AppModuleBasic{},
	feeshare.AppModuleBasic{},
	icq.AppModuleBasic{},
)

// NOTE: Any module instantiated in the module manager that is later modified
// must be passed by reference here.
func appModules(app *TerraApp, encodingConfig terrappsparams.EncodingConfig, skipGenesisInvariants bool) []module.AppModule {
	return []module.AppModule{
		genutil.NewAppModule(
			app.Keepers.AccountKeeper, app.Keepers.StakingKeeper, app.BaseApp.DeliverTx,
			encodingConfig.TxConfig,
		),
		auth.NewAppModule(app.appCodec, app.Keepers.AccountKeeper, nil, app.GetSubspace(authtypes.ModuleName)),
		vesting.NewAppModule(app.Keepers.AccountKeeper, app.Keepers.BankKeeper, app.Keepers.DistrKeeper, app.Keepers.StakingKeeper),
		custombankmodule.NewAppModule(app.appCodec, app.Keepers.BankKeeper, app.Keepers.AccountKeeper, app.GetSubspace(banktypes.ModuleName)),
		capability.NewAppModule(app.appCodec, *app.Keepers.CapabilityKeeper, false),
		crisis.NewAppModule(&app.Keepers.CrisisKeeper, skipGenesisInvariants, app.GetSubspace(crisistypes.ModuleName)),
		feegrantmodule.NewAppModule(app.appCodec, app.Keepers.AccountKeeper, app.Keepers.BankKeeper, app.Keepers.FeeGrantKeeper, app.interfaceRegistry),
		gov.NewAppModule(app.appCodec, &app.Keepers.GovKeeper, app.Keepers.AccountKeeper, app.Keepers.BankKeeper, app.GetSubspace(govtypes.ModuleName)),
		mint.NewAppModule(app.appCodec, app.Keepers.MintKeeper, app.Keepers.AccountKeeper, nil, app.GetSubspace(minttypes.ModuleName)),
		slashing.NewAppModule(app.appCodec, app.Keepers.SlashingKeeper, app.Keepers.AccountKeeper, app.Keepers.BankKeeper, app.Keepers.StakingKeeper, app.GetSubspace(slashingtypes.ModuleName)),
		distr.NewAppModule(app.appCodec, app.Keepers.DistrKeeper, app.Keepers.AccountKeeper, app.Keepers.BankKeeper, app.Keepers.StakingKeeper, app.GetSubspace(distrtypes.ModuleName)),
		staking.NewAppModule(app.appCodec, app.Keepers.StakingKeeper, app.Keepers.AccountKeeper, app.Keepers.BankKeeper, app.GetSubspace(stakingtypes.ModuleName)),
		consensus.NewAppModule(app.appCodec, app.Keepers.ConsensusParamsKeeper),
		upgrade.NewAppModule(app.Keepers.UpgradeKeeper),
		evidence.NewAppModule(app.Keepers.EvidenceKeeper),
		ibc.NewAppModule(app.Keepers.IBCKeeper),
		params.NewAppModule(app.Keepers.ParamsKeeper),
		authzmodule.NewAppModule(app.appCodec, app.Keepers.AuthzKeeper, app.Keepers.AccountKeeper, app.Keepers.BankKeeper, app.interfaceRegistry),
		ibctransfer.NewAppModule(app.Keepers.TransferKeeper),
		ibcfee.NewAppModule(app.Keepers.IBCFeeKeeper),
		ica.NewAppModule(&app.Keepers.ICAControllerKeeper, &app.Keepers.ICAHostKeeper),
		packetforward.NewAppModule(&app.Keepers.PacketForwardKeeper, app.GetSubspace(packetforwardtypes.ModuleName)),
		customwasmodule.NewAppModule(app.appCodec, &app.Keepers.WasmKeeper, app.Keepers.StakingKeeper, app.Keepers.AccountKeeper, app.Keepers.BankKeeper, app.MsgServiceRouter(), app.GetSubspace(wasmtypes.ModuleName)),
		ibchooks.NewAppModule(app.Keepers.AccountKeeper),
		tokenfactory.NewAppModule(app.Keepers.TokenFactoryKeeper, app.Keepers.AccountKeeper, app.Keepers.BankKeeper, app.GetSubspace(tokenfactorytypes.ModuleName)),
		alliance.NewAppModule(app.appCodec, app.Keepers.AllianceKeeper, app.Keepers.StakingKeeper, app.Keepers.AccountKeeper, app.Keepers.BankKeeper, app.interfaceRegistry, app.GetSubspace(alliancetypes.ModuleName)),
		feeshare.NewAppModule(app.Keepers.FeeShareKeeper, app.Keepers.AccountKeeper, app.GetSubspace(feesharetypes.ModuleName)),
		icq.NewAppModule(app.Keepers.ICQKeeper),
	}
}

// NOTE: The genutils module must occur after staking so that pools are
// properly initialized with tokens from genesis accounts.
// NOTE: Capability module must occur first so that it can initialize any capabilities
// so that other modules that want to create or claim capabilities afterwards in InitChain
// can do so safely.
var initGenesisOrder = []string{
	capabilitytypes.ModuleName,
	authtypes.ModuleName,
	banktypes.ModuleName,
	distrtypes.ModuleName,
	stakingtypes.ModuleName,
	slashingtypes.ModuleName,
	govtypes.ModuleName,
	minttypes.ModuleName,
	crisistypes.ModuleName,
	genutiltypes.ModuleName,
	evidencetypes.ModuleName,
	authz.ModuleName,
	paramstypes.ModuleName,
	upgradetypes.ModuleName,
	vestingtypes.ModuleName,
	feegrant.ModuleName,
	ibcexported.ModuleName,
	ibctransfertypes.ModuleName,
	icatypes.ModuleName,
	ibcfeetypes.ModuleName,
	packetforwardtypes.ModuleName,
	tokenfactorytypes.ModuleName,
	ibchookstypes.ModuleName,
	wasmtypes.ModuleName,
	alliancetypes.ModuleName,
	feesharetypes.ModuleName,
	consensusparamtypes.ModuleName,
	icqtypes.ModuleName,
}

var beginBlockersOrder = []string{
	upgradetypes.ModuleName,
	capabilitytypes.ModuleName,
	minttypes.ModuleName,
	distrtypes.ModuleName,
	slashingtypes.ModuleName,
	evidencetypes.ModuleName,
	stakingtypes.ModuleName,
	authtypes.ModuleName,
	banktypes.ModuleName,
	govtypes.ModuleName,
	crisistypes.ModuleName,
	genutiltypes.ModuleName,
	authz.ModuleName,
	feegrant.ModuleName,
	paramstypes.ModuleName,
	vestingtypes.ModuleName,
	// additional modules
	ibcexported.ModuleName,
	ibctransfertypes.ModuleName,
	icatypes.ModuleName,
	ibcfeetypes.ModuleName,
	packetforwardtypes.ModuleName,
	ibchookstypes.ModuleName,
	wasmtypes.ModuleName,
	tokenfactorytypes.ModuleName,
	alliancetypes.ModuleName,
	feesharetypes.ModuleName,
	consensusparamtypes.ModuleName,
	icqtypes.ModuleName,
}

var endBlockerOrder = []string{
	crisistypes.ModuleName,
	govtypes.ModuleName,
	stakingtypes.ModuleName,
	capabilitytypes.ModuleName,
	authtypes.ModuleName,
	banktypes.ModuleName,
	distrtypes.ModuleName,
	slashingtypes.ModuleName,
	minttypes.ModuleName,
	genutiltypes.ModuleName,
	evidencetypes.ModuleName,
	authz.ModuleName,
	feegrant.ModuleName,
	paramstypes.ModuleName,
	upgradetypes.ModuleName,
	vestingtypes.ModuleName,
	// additional non simd modules
	ibcexported.ModuleName,
	ibctransfertypes.ModuleName,
	icatypes.ModuleName,
	ibcfeetypes.ModuleName,
	packetforwardtypes.ModuleName,
	ibchookstypes.ModuleName,
	wasmtypes.ModuleName,
	tokenfactorytypes.ModuleName,
	alliancetypes.ModuleName,
	feesharetypes.ModuleName,
	consensusparamtypes.ModuleName,
	icqtypes.ModuleName,
}
