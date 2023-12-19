package test_helpers

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	dbm "github.com/cometbft/cometbft-db"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	ibcfee "github.com/cosmos/ibc-go/v7/modules/apps/29-fee"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/terra-money/alliance/x/alliance"
	"github.com/terra-money/core/v2/x/feeshare"
	"github.com/terra-money/core/v2/x/tokenfactory"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	mocktestutils "github.com/cosmos/cosmos-sdk/testutil/mock"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	tmtypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v7/router"
	icq "github.com/cosmos/ibc-apps/modules/async-icq/v7"
	ibchooks "github.com/cosmos/ibc-apps/modules/ibc-hooks/v7"
	ica "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts"
	"github.com/cosmos/ibc-go/v7/modules/apps/transfer"
	ibc "github.com/cosmos/ibc-go/v7/modules/core"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
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

	"github.com/CosmWasm/wasmd/x/wasm"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	terra_app "github.com/terra-money/core/v2/app"
	appparams "github.com/terra-money/core/v2/app/params"
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

func (ats *AppTestSuite) TestSimAppExportAndBlockedAddrs(t *testing.T) {
	encCfg := terra_app.MakeEncodingConfig()
	db := dbm.NewMemDB()
	app := terra_app.NewTerraApp(
		log.NewTMLogger(log.NewSyncWriter(os.Stdout)),
		db, nil, true, map[int64]bool{}, terra_app.DefaultNodeHome, 0, encCfg,
		simtestutil.EmptyAppOptions{}, wasmtypes.DefaultWasmConfig())

	// generate validator private/public key
	privVal := mocktestutils.NewPV()
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
	app2 := terra_app.NewTerraApp(
		log.NewTMLogger(log.NewSyncWriter(os.Stdout)),
		db, nil, true, map[int64]bool{}, terra_app.DefaultNodeHome, 0,
		encCfg, simtestutil.EmptyAppOptions{}, wasmtypes.DefaultWasmConfig())
	_, err = app2.ExportAppStateAndValidators(false, []string{}, []string{})
	require.NoError(t, err, "ExportAppStateAndValidators should not have an error")
}

func TestInitGenesisOnMigration(t *testing.T) {
	db := dbm.NewMemDB()
	encCfg := terra_app.MakeEncodingConfig()
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	app := terra_app.NewTerraApp(
		logger, db, nil, true, map[int64]bool{},
		terra_app.DefaultNodeHome, 0, encCfg, simtestutil.EmptyAppOptions{}, wasmtypes.DefaultWasmConfig())

	ctx := app.NewContext(true, tmproto.Header{Height: app.LastBlockHeight()})

	// Create a mock module. This module will serve as the new module we're
	// adding during a migration.
	mockCtrl := gomock.NewController(t)
	t.Cleanup(mockCtrl.Finish)
	mockModule := mocktestutils.NewMockAppModuleWithAllExtensions(mockCtrl)
	mockDefaultGenesis := json.RawMessage(`{"key": "value"}`)
	mockModule.EXPECT().DefaultGenesis(gomock.Eq(app.AppCodec())).Times(1).Return(mockDefaultGenesis)
	mockModule.EXPECT().InitGenesis(gomock.Eq(ctx), gomock.Eq(app.AppCodec()), gomock.Eq(mockDefaultGenesis)).Times(1).Return(nil)
	mockModule.EXPECT().ConsensusVersion().Times(1).Return(uint64(0))

	app.GetModuleManager().Modules["mock"] = mockModule

	// Run migrations only for "mock" module. We exclude it from
	// the VersionMap to simulate upgrading with a new module.
	res, err := app.GetModuleManager().RunMigrations(ctx, app.GetConfigurator(),
		module.VersionMap{
			"alliance":               alliance.AppModule{}.ConsensusVersion(),
			"auth":                   auth.AppModule{}.ConsensusVersion(),
			"authz":                  authzmodule.AppModule{}.ConsensusVersion(),
			"bank":                   bank.AppModule{}.ConsensusVersion(),
			"capability":             capability.AppModule{}.ConsensusVersion(),
			"crisis":                 crisis.AppModule{}.ConsensusVersion(),
			"distribution":           distribution.AppModule{}.ConsensusVersion(),
			"evidence":               evidence.AppModule{}.ConsensusVersion(),
			"feegrant":               feegrantmodule.AppModule{}.ConsensusVersion(),
			"feeshare":               feeshare.AppModule{}.ConsensusVersion(),
			"feeibc":                 ibcfee.AppModule{}.ConsensusVersion(),
			"genutil":                genutil.AppModule{}.ConsensusVersion(),
			"gov":                    gov.AppModule{}.ConsensusVersion(),
			"ibc":                    ibc.AppModule{}.ConsensusVersion(),
			"interchainquery":        icq.AppModule{}.ConsensusVersion(),
			"ibchooks":               ibchooks.AppModule{}.ConsensusVersion(),
			"interchainaccounts":     ica.AppModule{}.ConsensusVersion(),
			"mint":                   mint.AppModule{}.ConsensusVersion(),
			"packetfowardmiddleware": router.AppModule{}.ConsensusVersion(),
			"params":                 params.AppModule{}.ConsensusVersion(),
			"slashing":               slashing.AppModule{}.ConsensusVersion(),
			"staking":                staking.AppModule{}.ConsensusVersion(),
			"tokenfactory":           tokenfactory.AppModule{}.ConsensusVersion(),
			"transfer":               transfer.AppModule{}.ConsensusVersion(),
			"upgrade":                upgrade.AppModule{}.ConsensusVersion(),
			"vesting":                vesting.AppModule{}.ConsensusVersion(),
			"wasm":                   wasm.AppModule{}.ConsensusVersion(),
		},
	)
	require.NoError(t, err)
	require.Equal(t, res, module.VersionMap{
		"alliance":               5,
		"auth":                   4,
		"authz":                  2,
		"bank":                   4,
		"capability":             1,
		"consensus":              1,
		"crisis":                 2,
		"distribution":           3,
		"evidence":               1,
		"feegrant":               2,
		"feeshare":               2,
		"feeibc":                 1,
		"genutil":                1,
		"gov":                    4,
		"ibc":                    4,
		"ibchooks":               1,
		"interchainaccounts":     2,
		"interchainquery":        1,
		"mint":                   2,
		"mock":                   0,
		"packetfowardmiddleware": 1,
		"params":                 1,
		"slashing":               3,
		"staking":                4,
		"tokenfactory":           3,
		"transfer":               3,
		"upgrade":                2,
		"vesting":                1,
		"wasm":                   4,
	})
}

func (ats *AppTestSuite) TestCodecs(t *testing.T) {
	ats.Setup()

	encCfg := terra_app.MakeEncodingConfig()

	// Vlidate that the tests contain the correct encoding configuration
	require.Equal(t, encCfg.Amino, ats.App.LegacyAmino())
	require.Equal(t, encCfg.Marshaler, ats.App.AppCodec())
	require.Equal(t, encCfg.InterfaceRegistry, ats.App.InterfaceRegistry())
	require.NotEmpty(t, ats.App.GetKey(banktypes.StoreKey))
	require.NotEmpty(t, ats.App.GetTKey(paramstypes.TStoreKey))
	require.NotEmpty(t, ats.App.GetMemKey(capabilitytypes.MemStoreKey))
}

func TestSimAppEnforceStakingForVestingTokens(t *testing.T) {
	encCfg := terra_app.MakeEncodingConfig()
	db := dbm.NewMemDB()
	app := terra_app.NewTerraApp(
		log.NewTMLogger(log.NewSyncWriter(os.Stdout)),
		db, nil, true, map[int64]bool{}, terra_app.DefaultNodeHome, 0, encCfg,
		simtestutil.EmptyAppOptions{}, wasmtypes.DefaultWasmConfig(),
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
	privVal := mocktestutils.NewPV()
	pubKey, err := privVal.GetPubKey()
	require.NoError(t, err, "PubKey should not have an error")
	validator := tmtypes.NewValidator(pubKey, 1)
	valSet := tmtypes.NewValidatorSet([]*tmtypes.Validator{validator})

	genesisState := SetupGenesisValSet(valSet, genAccounts, nil, app, encCfg, balances...)
	ctx := app.NewContext(true, tmproto.Header{Height: app.LastBlockHeight()})

	genesisState[authtypes.ModuleName] = app.GetAppCodec().MustMarshalJSON(authtypes.NewGenesisState(authtypes.DefaultParams(), genAccounts))
	delegations := app.Keepers.StakingKeeper.GetAllDelegations(ctx)
	sharePerValidators := make(map[string]sdk.Dec)

	for _, del := range delegations {
		if val, found := sharePerValidators[del.ValidatorAddress]; !found {
			sharePerValidators[del.ValidatorAddress] = del.GetShares()
		} else {
			sharePerValidators[del.ValidatorAddress] = val.Add(del.GetShares())
		}
	}

	/* #nosec */
	for _, share := range sharePerValidators {
		require.Equal(t, sdk.NewDec(3_500_001_000_000), share)
	}
}

func SetupGenesisValSet(
	valSet *tmtypes.ValidatorSet,
	genAccs []authtypes.GenesisAccount,
	opts []wasm.Option,
	app *terra_app.TerraApp,
	encCfg appparams.EncodingConfig,
	balances ...banktypes.Balance,
) terra_app.GenesisState {
	genesisState := terra_app.NewDefaultGenesisState(encCfg.Marshaler)
	// set genesis accounts
	authGenesis := authtypes.NewGenesisState(authtypes.DefaultParams(), genAccs)
	genesisState[authtypes.ModuleName] = app.AppCodec().MustMarshalJSON(authGenesis)

	validators := make([]stakingtypes.Validator, 0, len(valSet.Validators))
	delegations := make([]stakingtypes.Delegation, 0, len(valSet.Validators))

	bondAmt := sdk.NewInt(1000000)
	totalSupply := sdk.NewCoins()

	for _, val := range valSet.Validators {
		pk, err := cryptocodec.FromTmPubKeyInterface(val.PubKey)
		if err != nil {
			panic(err)
		}

		pkAny, err := codectypes.NewAnyWithValue(pk)
		if err != nil {
			panic(err)
		}
		validator := stakingtypes.Validator{
			OperatorAddress:   sdk.ValAddress(val.Address).String(),
			ConsensusPubkey:   pkAny,
			Jailed:            false,
			Status:            stakingtypes.Bonded,
			Tokens:            bondAmt,
			DelegatorShares:   sdk.OneDec(),
			Description:       stakingtypes.Description{},
			UnbondingHeight:   int64(0),
			UnbondingTime:     time.Unix(0, 0).UTC(),
			Commission:        stakingtypes.NewCommission(sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec()),
			MinSelfDelegation: sdk.ZeroInt(),
		}
		validators = append(validators, validator)
		delegations = append(delegations, stakingtypes.NewDelegation(genAccs[0].GetAddress(), val.Address.Bytes(), sdk.OneDec()))

	}

	// set validators and delegations
	stakingGenesis := stakingtypes.NewGenesisState(stakingtypes.DefaultParams(), validators, delegations)
	genesisState[stakingtypes.ModuleName] = app.AppCodec().MustMarshalJSON(stakingGenesis)

	// add bonded amount to bonded pool module account
	balances = append(balances, banktypes.Balance{
		Address: authtypes.NewModuleAddress(stakingtypes.BondedPoolName).String(),
		Coins:   sdk.Coins{sdk.NewCoin(sdk.DefaultBondDenom, bondAmt)},
	})

	for _, b := range balances {
		// add genesis acc tokens and delegated tokens to total supply
		totalSupply = totalSupply.Add(b.Coins...)
	}
	// update total supply
	bankGenesis := banktypes.NewGenesisState(banktypes.DefaultGenesisState().Params, balances, totalSupply, []banktypes.Metadata{}, []banktypes.SendEnabled{})
	genesisState[banktypes.ModuleName] = app.AppCodec().MustMarshalJSON(bankGenesis)

	return genesisState
}
