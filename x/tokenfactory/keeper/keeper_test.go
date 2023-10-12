package keeper_test

import (
	"testing"
	"time"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/stretchr/testify/suite"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"

	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	app "github.com/terra-money/core/v2/app/app_test"
	"github.com/terra-money/core/v2/x/tokenfactory/keeper"
	"github.com/terra-money/core/v2/x/tokenfactory/types"
)

type KeeperTestSuite struct {
	app.AppTestSuite

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
	s.contractKeeper = wasmkeeper.NewGovPermissionKeeper(s.App.WasmKeeper)
	s.queryClient = types.NewQueryClient(s.QueryHelper)
	s.msgServer = keeper.NewMsgServerImpl(s.App.TokenFactoryKeeper)
	s.bankMsgServer = bankkeeper.NewMsgServerImpl(s.App.BankKeeper)
}

func (s *KeeperTestSuite) TestCreateModuleAccount() {
	// setup new next account number
	nextAccountNumber := s.App.AccountKeeper.NextAccountNumber(s.Ctx)

	// ensure module account was removed
	s.Ctx = s.App.NewContext(true, tmproto.Header{Time: time.Now()})
	tokenfactoryModuleAccount := s.App.AccountKeeper.GetAccount(s.Ctx, s.App.AccountKeeper.GetModuleAddress(types.ModuleName))
	s.Require().Nil(tokenfactoryModuleAccount)

	// create module account
	s.App.TokenFactoryKeeper.CreateModuleAccount(s.Ctx)

	// check that the module account is now initialized
	tokenfactoryModuleAccount = s.App.AccountKeeper.GetAccount(s.Ctx, s.App.AccountKeeper.GetModuleAddress(types.ModuleName))
	s.Require().NotNil(tokenfactoryModuleAccount)

	// check that the account number of the module account is now initialized correctly
	s.Require().Equal(nextAccountNumber+1, tokenfactoryModuleAccount.GetAccountNumber())
}
