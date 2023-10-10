package keeper_test

import (
	"testing"
	"time"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/stretchr/testify/suite"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"

	app_test "github.com/terra-money/core/v2/app/app_test"
	"github.com/terra-money/core/v2/x/feeshare/keeper"
	"github.com/terra-money/core/v2/x/feeshare/types"
)

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	SendCoins(ctx sdk.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromModuleToModule(ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coins) error
}

type IntegrationTestSuite struct {
	*app_test.AppTestSuite

	ctx               sdk.Context
	bankKeeper        BankKeeper
	accountKeeper     types.AccountKeeper
	queryClient       types.QueryClient
	feeShareMsgServer types.MsgServer
	wasmMsgServer     wasmtypes.MsgServer
}

func (s *IntegrationTestSuite) SetupTest() {
	isCheckTx := false
	s.AppTestSuite.Setup()

	s.ctx = s.App.BaseApp.NewContext(isCheckTx, tmproto.Header{
		ChainID: "testing",
		Height:  9,
		Time:    time.Now().UTC(),
	})

	queryHelper := baseapp.NewQueryServerTestHelper(s.ctx, s.App.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, keeper.NewQuerier(s.App.FeeShareKeeper))

	s.queryClient = types.NewQueryClient(queryHelper)
	s.bankKeeper = s.App.BankKeeper
	s.accountKeeper = s.App.AccountKeeper
	s.feeShareMsgServer = s.App.FeeShareKeeper
	s.wasmMsgServer = wasmkeeper.NewMsgServerImpl(&s.App.WasmKeeper)
}

func (s *IntegrationTestSuite) FundAccount(ctx sdk.Context, addr sdk.AccAddress, amounts sdk.Coins) error {
	if err := s.bankKeeper.MintCoins(ctx, minttypes.ModuleName, amounts); err != nil {
		return err
	}

	return s.bankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, addr, amounts)
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
