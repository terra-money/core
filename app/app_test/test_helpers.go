package app

// DONTCOVER

import (
	"os"
	"time"

	"reflect"
	"unsafe"

	"github.com/CosmWasm/wasmd/x/wasm"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	dbm "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/crypto/ed25519"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmtypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/suite"
	"github.com/terra-money/core/v2/app"
	terra_app "github.com/terra-money/core/v2/app"
	appparams "github.com/terra-money/core/v2/app/params"
	terrraParams "github.com/terra-money/core/v2/app/params"
	"github.com/terra-money/core/v2/app/wasmconfig"
	tokenfactorytypes "github.com/terra-money/core/v2/x/tokenfactory/types"
)

type AppTestSuite struct {
	suite.Suite

	App         *app.TerraApp
	Ctx         sdk.Context
	QueryHelper *baseapp.QueryServiceTestHelper
	TestAccs    []sdk.AccAddress
}

// Setup sets up basic environment for suite (App, Ctx, and test accounts)
func (s *AppTestSuite) Setup() {
	appparams.RegisterAddressesConfig()
	baseTestAccts := CreateRandomAccounts(3)
	encCfg := terra_app.MakeEncodingConfig()

	db := dbm.NewMemDB()
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
		wasmconfig.DefaultConfig(),
	)

	s.TestAccs = []sdk.AccAddress{}
	s.TestAccs = append(s.TestAccs, baseTestAccts...)
	s.Ctx = s.App.NewContext(true, tmproto.Header{Height: 1, Time: time.Now()})
	s.QueryHelper = &baseapp.QueryServiceTestHelper{
		GRPCQueryRouter: s.App.GRPCQueryRouter(),
		Ctx:             s.Ctx,
	}
	err := s.App.BankKeeper.SetParams(s.Ctx, banktypes.NewParams(true))
	s.Require().NoError(err)
	err = s.App.WasmKeeper.SetParams(s.Ctx, wasmtypes.DefaultParams())
	s.Require().NoError(err)

	err = s.App.TokenFactoryKeeper.SetParams(s.Ctx, tokenfactorytypes.DefaultParams())
	s.Require().NoError(err)
	s.App.DistrKeeper.SetFeePool(s.Ctx, distrtypes.InitialFeePool())
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
func CreateRandomAccounts(numAccts int) []sdk.AccAddress {
	testAddrs := make([]sdk.AccAddress, numAccts)
	for i := 0; i < numAccts; i++ {
		pk := ed25519.GenPrivKey().PubKey()
		testAddrs[i] = sdk.AccAddress(pk.Address())
	}

	return testAddrs
}

// FundAcc funds target address with specified amount.
func (s *AppTestSuite) FundAcc(acc sdk.AccAddress, amounts sdk.Coins) (err error) {
	s.Require().NoError(err)
	if err := s.App.BankKeeper.MintCoins(s.Ctx, minttypes.ModuleName, amounts); err != nil {
		return err
	}

	return s.App.BankKeeper.SendCoinsFromModuleToAccount(s.Ctx, minttypes.ModuleName, acc, amounts)
}

func SetupGenesisValSet(
	valSet *tmtypes.ValidatorSet,
	genAccs []authtypes.GenesisAccount,
	opts []wasm.Option,
	app *terra_app.TerraApp,
	encCfg terrraParams.EncodingConfig,
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

var (
	coinType  = reflect.TypeOf(sdk.Coin{})
	coinsType = reflect.TypeOf(sdk.Coins{})
)

// Fill analyze all struct fields and slices with
// reflection and initialize the nil and empty slices,
// structs, and pointers.
func Fill(x interface{}) interface{} {
	v := reflect.Indirect(reflect.ValueOf(x))
	switch v.Kind() {
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			obj := v.Index(i)
			objPt := reflect.NewAt(obj.Type(), unsafe.Pointer(obj.UnsafeAddr())).Interface()
			objPt = Fill(objPt)
			obj.Set(reflect.ValueOf(objPt))
		}
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			f := reflect.Indirect(v.Field(i))
			if !f.CanSet() {
				continue
			}
			switch f.Kind() {
			case reflect.Slice:
				f.Set(reflect.MakeSlice(f.Type(), 0, 0))
			case reflect.Struct:
				switch f.Type() {
				case coinType:
					coin := reflect.New(coinType).Interface()
					s := reflect.ValueOf(coin).Elem()
					f.Set(s)
				case coinsType:
					coins := reflect.New(coinsType).Interface()
					s := reflect.ValueOf(coins).Elem()
					f.Set(s)
				default:
					objPt := reflect.NewAt(f.Type(), unsafe.Pointer(f.UnsafeAddr())).Interface()
					s := Fill(objPt)
					f.Set(reflect.ValueOf(s))
				}
			}
		}
	}
	return reflect.Indirect(v).Interface()
}
