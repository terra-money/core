package test_helpers

import (
	"os"
	"time"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	dbm "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/stretchr/testify/suite"
	terra_app "github.com/terra-money/core/v2/app"
	appparams "github.com/terra-money/core/v2/app/params"
	feesharetypes "github.com/terra-money/core/v2/x/feeshare/types"
	tokenfactorytypes "github.com/terra-money/core/v2/x/tokenfactory/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	secp256k1 "github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
)

type AppTestSuite struct {
	suite.Suite

	App            *terra_app.TerraApp
	DB             dbm.DB
	Ctx            sdk.Context
	QueryHelper    *baseapp.QueryServiceTestHelper
	TestAccs       []sdk.AccAddress
	TestAccPrivs   []cryptotypes.PrivKey
	EncodingConfig appparams.EncodingConfig
}

// Setup sets up basic environment for suite (App, Ctx, and test accounts)
func (s *AppTestSuite) Setup() {
	appparams.RegisterAddressesConfig()
	encCfg := terra_app.MakeEncodingConfig()
	genesisState := terra_app.NewDefaultGenesisState(encCfg.Marshaler)
	genesisState.SetDefaultTerraConfig(encCfg.Marshaler)

	db := dbm.NewMemDB()
	s.DB = db
	s.App = terra_app.NewTerraApp(
		log.NewTMLogger(log.NewSyncWriter(os.Stdout)),
		db,
		nil,
		true,
		map[int64]bool{},
		terra_app.DefaultNodeHome,
		0,
		encCfg,
		simtestutil.EmptyAppOptions{},
		wasmtypes.DefaultWasmConfig(),
	)
	s.EncodingConfig = encCfg

	s.Ctx = s.App.NewContext(true, tmproto.Header{Height: 1, Time: time.Now()})
	s.QueryHelper = &baseapp.QueryServiceTestHelper{
		GRPCQueryRouter: s.App.GRPCQueryRouter(),
		Ctx:             s.Ctx,
	}
	err := s.App.Keepers.BankKeeper.SetParams(s.Ctx, banktypes.NewParams(true))
	s.Require().NoError(err)

	err = s.App.Keepers.WasmKeeper.SetParams(s.Ctx, wasmtypes.DefaultParams())
	s.Require().NoError(err)

	err = s.App.Keepers.FeeShareKeeper.SetParams(s.Ctx, feesharetypes.DefaultParams())
	s.Require().NoError(err)

	err = s.App.Keepers.TokenFactoryKeeper.SetParams(s.Ctx, tokenfactorytypes.DefaultParams())
	s.Require().NoError(err)

	err = s.FundModule(authtypes.FeeCollectorName, sdk.NewCoins(sdk.NewCoin("uluna", sdk.NewInt(1000)), sdk.NewCoin("utoken", sdk.NewInt(500))))
	s.Require().NoError(err)

	s.App.Keepers.DistrKeeper.SetFeePool(s.Ctx, distrtypes.InitialFeePool())

	testAccounts, privKeys := s.CreateRandomAccounts(3)
	s.TestAccs = testAccounts
	s.TestAccPrivs = privKeys
}

func (s *AppTestSuite) AssertEventEmitted(ctx sdk.Context, eventTypeExpected string, numEventsExpected int) {
	allEvents := ctx.EventManager().Events()
	// filter out other events
	actualEvents := make([]sdk.Event, 0)
	for _, event := range allEvents {
		if event.Type == eventTypeExpected {
			actualEvents = append(actualEvents, event)
		}
	}
	s.Require().Equal(numEventsExpected, len(actualEvents))
}

// CreateRandomAccounts is a function return a list of randomly generated AccAddresses
func (s *AppTestSuite) CreateRandomAccounts(numAccts int) ([]sdk.AccAddress, []cryptotypes.PrivKey) {
	testAddrs := make([]sdk.AccAddress, numAccts)
	testPrivKeys := make([]cryptotypes.PrivKey, numAccts)
	for i := 0; i < numAccts; i++ {
		priv := secp256k1.GenPrivKey()
		pk := priv.PubKey()
		testAddrs[i] = sdk.AccAddress(pk.Address())
		testPrivKeys[i] = priv

		err := s.FundAcc(testAddrs[i], sdk.NewCoins(sdk.NewInt64Coin("uluna", 100000000)))
		s.Require().NoError(err)
	}

	return testAddrs, testPrivKeys
}

// FundAcc funds target address with specified amount.
func (s *AppTestSuite) FundAcc(acc sdk.AccAddress, amounts sdk.Coins) (err error) {
	s.Require().NoError(err)
	if err := s.App.Keepers.BankKeeper.MintCoins(s.Ctx, minttypes.ModuleName, amounts); err != nil {
		return err
	}

	return s.App.Keepers.BankKeeper.SendCoinsFromModuleToAccount(s.Ctx, minttypes.ModuleName, acc, amounts)
}

// FundAcc funds target address with specified amount.
func (s *AppTestSuite) FundModule(moduleAccount string, amounts sdk.Coins) (err error) {
	s.Require().NoError(err)
	if err := s.App.Keepers.BankKeeper.MintCoins(s.Ctx, minttypes.ModuleName, amounts); err != nil {
		return err
	}

	return s.App.Keepers.BankKeeper.SendCoinsFromModuleToModule(s.Ctx, minttypes.ModuleName, moduleAccount, amounts)
}
