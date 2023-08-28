package keeper_test

import (
	"testing"
	"time"

	"cosmossdk.io/math"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"

	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/terra-money/core/v2/app/app_testing"
	"github.com/terra-money/core/v2/app/config"
	"github.com/terra-money/core/v2/x/tokenfactory/keeper"
	"github.com/terra-money/core/v2/x/tokenfactory/types"
)

type KeeperTestSuite struct {
	app_testing.AppTestSuite

	queryClient    types.QueryClient
	msgServer      types.MsgServer
	contractKeeper wasmtypes.ContractOpsKeeper
	bankMsgServer  banktypes.MsgServer
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (s *KeeperTestSuite) SetupTest() {
	s.Setup()
	// Fund every TestAcc with two denoms, one of which is the denom creation fee
	fundAccsAmount := sdk.NewCoins(sdk.NewCoin(config.BondDenom, math.NewInt(1_000_000_000)))
	for _, acc := range s.TestAccs {
		s.FundAcc(acc, fundAccsAmount)
	}
	s.contractKeeper = wasmkeeper.NewGovPermissionKeeper(s.App.WasmKeeper)
	s.queryClient = types.NewQueryClient(s.QueryHelper)
	s.msgServer = keeper.NewMsgServerImpl(s.App.TokenFactoryKeeper)
	s.bankMsgServer = bankkeeper.NewMsgServerImpl(s.App.BankKeeper)
}

func (s *KeeperTestSuite) TestCreateModuleAccount() {
	s.Setup()
	app := s.App

	// setup new next account number
	nextAccountNumber := app.AccountKeeper.NextAccountNumber(s.Ctx)

	// ensure module account was removed
	s.Ctx = app.NewContext(true, tmproto.Header{Time: time.Now()})
	tokenfactoryModuleAccount := app.AccountKeeper.GetAccount(s.Ctx, app.AccountKeeper.GetModuleAddress(types.ModuleName))
	s.Require().Nil(tokenfactoryModuleAccount)

	// create module account
	app.TokenFactoryKeeper.CreateModuleAccount(s.Ctx)

	// check that the module account is now initialized
	tokenfactoryModuleAccount = app.AccountKeeper.GetAccount(s.Ctx, app.AccountKeeper.GetModuleAddress(types.ModuleName))
	s.Require().NotNil(tokenfactoryModuleAccount)

	// check that the account number of the module account is now initialized correctly
	s.Require().Equal(nextAccountNumber+1, tokenfactoryModuleAccount.GetAccountNumber())
}
