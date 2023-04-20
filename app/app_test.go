package app

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"
	"github.com/terra-money/core/v2/app/wasmconfig"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/tests/mocks"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	authzmodule "github.com/cosmos/cosmos-sdk/x/authz/module"
	"github.com/cosmos/cosmos-sdk/x/bank"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/capability"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	feegrantmodule "github.com/cosmos/cosmos-sdk/x/feegrant/module"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	ica "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts"
	"github.com/cosmos/ibc-go/v6/modules/apps/transfer"
	ibc "github.com/cosmos/ibc-go/v6/modules/core"
	"github.com/cosmos/ibc-go/v6/testing/mock"
	"github.com/strangelove-ventures/packet-forward-middleware/v6/router"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/CosmWasm/wasmd/x/wasm"
)

var (
	priv1 = secp256k1.GenPrivKey()
	priv2 = secp256k1.GenPrivKey()
	priv3 = secp256k1.GenPrivKey()
	priv4 = secp256k1.GenPrivKey()
	pk1   = priv1.PubKey()
	pk2   = priv2.PubKey()
	pk3   = priv3.PubKey()
	pk4   = priv4.PubKey()
	addr1 = sdk.AccAddress(pk1.Address())
	addr2 = sdk.AccAddress(pk2.Address())
	addr3 = sdk.AccAddress(pk3.Address())
	addr4 = sdk.AccAddress(pk4.Address())
)

func TestSimAppExportAndBlockedAddrs(t *testing.T) {
	encCfg := MakeEncodingConfig()
	db := dbm.NewMemDB()
	app := NewTerraApp(
		log.NewTMLogger(log.NewSyncWriter(os.Stdout)),
		db, nil, true, map[int64]bool{}, simapp.DefaultNodeHome, 0, encCfg,
		simapp.EmptyAppOptions{}, wasmconfig.DefaultConfig())

	// generate validator private/public key
	privVal := mock.NewPV()
	pubKey, err := privVal.GetPubKey()
	require.NoError(t, err)

	// create validator set with single validator
	validator := tmtypes.NewValidator(pubKey, 1)
	valSet := tmtypes.NewValidatorSet([]*tmtypes.Validator{validator})

	// generate genesis account
	senderPrivKey := secp256k1.GenPrivKey()
	acc := authtypes.NewBaseAccount(senderPrivKey.PubKey().Address().Bytes(), senderPrivKey.PubKey(), 0, 0)
	balance := banktypes.Balance{
		Address: acc.GetAddress().String(),
		Coins:   sdk.NewCoins(),
	}

	genesisState := SetupGenesisValSet(valSet, []authtypes.GenesisAccount{acc}, nil, app, encCfg, balance)
	stateBytes, err := json.MarshalIndent(genesisState, "", "  ")
	require.NoError(t, err)

	// Initialize the chain
	app.InitChain(
		abci.RequestInitChain{
			Validators:    []abci.ValidatorUpdate{},
			AppStateBytes: stateBytes,
		},
	)
	app.Commit()

	// Making a new app object with the db, so that initchain hasn't been called
	app2 := NewTerraApp(
		log.NewTMLogger(log.NewSyncWriter(os.Stdout)),
		db, nil, true, map[int64]bool{}, simapp.DefaultNodeHome, 0,
		encCfg, simapp.EmptyAppOptions{}, wasmconfig.DefaultConfig())
	_, err = app2.ExportAppStateAndValidators(false, []string{})
	require.NoError(t, err, "ExportAppStateAndValidators should not have an error")
}

func TestGetMaccPerms(t *testing.T) {
	dup := GetMaccPerms()
	require.Equal(t, maccPerms, dup, "duplicated module account permissions differed from actual module account permissions")
}

func TestInitGenesisOnMigration(t *testing.T) {
	db := dbm.NewMemDB()
	encCfg := MakeEncodingConfig()
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	app := NewTerraApp(
		logger, db, nil, true, map[int64]bool{},
		simapp.DefaultNodeHome, 0, encCfg, simapp.EmptyAppOptions{}, wasmconfig.DefaultConfig())

	ctx := app.NewContext(true, tmproto.Header{Height: app.LastBlockHeight()})

	// Create a mock module. This module will serve as the new module we're
	// adding during a migration.
	mockCtrl := gomock.NewController(t)
	t.Cleanup(mockCtrl.Finish)
	mockModule := mocks.NewMockAppModule(mockCtrl)
	mockDefaultGenesis := json.RawMessage(`{"key": "value"}`)
	mockModule.EXPECT().DefaultGenesis(gomock.Eq(app.appCodec)).Times(1).Return(mockDefaultGenesis)
	mockModule.EXPECT().InitGenesis(gomock.Eq(ctx), gomock.Eq(app.appCodec), gomock.Eq(mockDefaultGenesis)).Times(1).Return(nil)
	mockModule.EXPECT().ConsensusVersion().Times(1).Return(uint64(0))

	app.mm.Modules["mock"] = mockModule

	// Run migrations only for "mock" module. We exclude it from
	// the VersionMap to simulate upgrading with a new module.
	_, err := app.mm.RunMigrations(ctx, app.configurator,
		module.VersionMap{
			"bank":                   bank.AppModule{}.ConsensusVersion(),
			"auth":                   auth.AppModule{}.ConsensusVersion(),
			"authz":                  authzmodule.AppModule{}.ConsensusVersion(),
			"staking":                staking.AppModule{}.ConsensusVersion(),
			"mint":                   mint.AppModule{}.ConsensusVersion(),
			"distribution":           distribution.AppModule{}.ConsensusVersion(),
			"slashing":               slashing.AppModule{}.ConsensusVersion(),
			"gov":                    gov.AppModule{}.ConsensusVersion(),
			"params":                 params.AppModule{}.ConsensusVersion(),
			"upgrade":                upgrade.AppModule{}.ConsensusVersion(),
			"feegrant":               feegrantmodule.AppModule{}.ConsensusVersion(),
			"evidence":               evidence.AppModule{}.ConsensusVersion(),
			"crisis":                 crisis.AppModule{}.ConsensusVersion(),
			"genutil":                genutil.AppModule{}.ConsensusVersion(),
			"capability":             capability.AppModule{}.ConsensusVersion(),
			"wasm":                   wasm.AppModule{}.ConsensusVersion(),
			"ibc":                    ibc.AppModule{}.ConsensusVersion(),
			"transfer":               transfer.AppModule{}.ConsensusVersion(),
			"interchainaccounts":     ica.AppModule{}.ConsensusVersion(),
			"packetfowardmiddleware": router.AppModule{}.ConsensusVersion(),
			"vesting":                vesting.AppModule{}.ConsensusVersion(),
		},
	)
	require.NoError(t, err)
}

func TestLegacyAmino(t *testing.T) {
	encCfg := MakeEncodingConfig()
	db := dbm.NewMemDB()
	app := NewTerraApp(
		log.NewTMLogger(log.NewSyncWriter(os.Stdout)),
		db, nil, true, map[int64]bool{}, DefaultNodeHome, 0,
		encCfg, simapp.EmptyAppOptions{}, wasmconfig.DefaultConfig())

	require.Equal(t, encCfg.Amino, app.LegacyAmino())
}

func TestAppCodec(t *testing.T) {
	encCfg := MakeEncodingConfig()
	db := dbm.NewMemDB()
	app := NewTerraApp(
		log.NewTMLogger(log.NewSyncWriter(os.Stdout)),
		db, nil, true, map[int64]bool{}, DefaultNodeHome, 0,
		encCfg, simapp.EmptyAppOptions{}, wasmconfig.DefaultConfig())

	require.Equal(t, encCfg.Marshaler, app.AppCodec())
}

func TestInterfaceRegistry(t *testing.T) {
	encCfg := MakeEncodingConfig()
	db := dbm.NewMemDB()
	app := NewTerraApp(
		log.NewTMLogger(log.NewSyncWriter(os.Stdout)),
		db, nil, true, map[int64]bool{}, DefaultNodeHome, 0,
		encCfg, simapp.EmptyAppOptions{}, wasmconfig.DefaultConfig())

	require.Equal(t, encCfg.InterfaceRegistry, app.InterfaceRegistry())
}

func TestGetKey(t *testing.T) {
	encCfg := MakeEncodingConfig()
	db := dbm.NewMemDB()
	app := NewTerraApp(
		log.NewTMLogger(log.NewSyncWriter(os.Stdout)),
		db, nil, true, map[int64]bool{}, DefaultNodeHome, 0,
		encCfg, simapp.EmptyAppOptions{}, wasmconfig.DefaultConfig())

	require.NotEmpty(t, app.GetKey(banktypes.StoreKey))
	require.NotEmpty(t, app.GetTKey(paramstypes.TStoreKey))
	require.NotEmpty(t, app.GetMemKey(capabilitytypes.MemStoreKey))
}

func TestSimAppEnforceStakingForVestingTokens(t *testing.T) {
	encCfg := MakeEncodingConfig()
	db := dbm.NewMemDB()
	app := NewTerraApp(
		log.NewTMLogger(log.NewSyncWriter(os.Stdout)),
		db, nil, true, map[int64]bool{}, simapp.DefaultNodeHome, 0, encCfg,
		simapp.EmptyAppOptions{}, wasmconfig.DefaultConfig(),
	)
	genAccounts := authtypes.GenesisAccounts{
		vestingtypes.NewContinuousVestingAccount(
			authtypes.NewBaseAccountWithAddress(addr1),
			sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(2_500_000_000_000))),
			1660000000,
			1670000000,
		),
		vestingtypes.NewContinuousVestingAccount(
			authtypes.NewBaseAccountWithAddress(addr2),
			sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(4_500_000_000_000))),
			1660000000,
			1670000000,
		),
		authtypes.NewBaseAccountWithAddress(addr3),
		authtypes.NewBaseAccountWithAddress(addr4),
	}
	balances := []banktypes.Balance{
		{
			Address: addr1.String(),
			Coins:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(2_500_000_000_000))),
		},
		{
			Address: addr2.String(),
			Coins:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(4_500_000_000_000))),
		},
		{
			Address: addr3.String(),
			Coins:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1_000_000))),
		},
		{
			Address: addr4.String(),
			Coins:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1_000_000))),
		},
	}

	// generate validator private/public key
	privVal := mock.NewPV()
	pubKey, err := privVal.GetPubKey()
	require.NoError(t, err, "PubKey should not have an error")
	validator := tmtypes.NewValidator(pubKey, 1)
	valSet := tmtypes.NewValidatorSet([]*tmtypes.Validator{validator})

	genesisState := SetupGenesisValSet(valSet, genAccounts, nil, app, encCfg, balances...)
	ctx := app.NewContext(true, tmproto.Header{Height: app.LastBlockHeight()})

	genesisState[authtypes.ModuleName] = app.appCodec.MustMarshalJSON(authtypes.NewGenesisState(authtypes.DefaultParams(), genAccounts))
	delegations := app.StakingKeeper.GetAllDelegations(ctx)
	sharePerValidators := make(map[string]sdk.Dec)

	for _, del := range delegations {
		if val, found := sharePerValidators[del.ValidatorAddress]; !found {
			sharePerValidators[del.ValidatorAddress] = del.GetShares()
		} else {
			sharePerValidators[del.ValidatorAddress] = val.Add(del.GetShares())
		}
	}

	for _, share := range sharePerValidators {
		require.Equal(t, sdk.NewDec(3_500_001_000_000), share)
	}
}
