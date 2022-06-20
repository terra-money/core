package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	terraapp "github.com/terra-money/core/v2/app"
	"github.com/terra-money/core/v2/app/ante/keeper"
	"github.com/terra-money/core/v2/app/ante/types"
	"github.com/terra-money/core/v2/app/wasmconfig"
)

// KeeperTestSuite is a test suite to be used with ante handler tests.
type KeeperTestSuite struct {
	suite.Suite

	app         *terraapp.TerraApp
	ctx         sdk.Context
	queryClient types.QueryClient
}

// returns context and app with params set on account keeper
func createTestApp(isCheckTx bool, tempDir string) (*terraapp.TerraApp, sdk.Context) {
	app := terraapp.NewTerraApp(
		log.NewNopLogger(), dbm.NewMemDB(), nil, true, map[int64]bool{},
		tempDir, simapp.FlagPeriodValue, terraapp.MakeEncodingConfig(),
		simapp.EmptyAppOptions{}, wasmconfig.DefaultConfig(),
	)
	ctx := app.BaseApp.NewContext(isCheckTx, tmproto.Header{})
	app.AccountKeeper.SetParams(ctx, authtypes.DefaultParams())
	app.AnteKeeper.SetParams(ctx, types.DefaultParams())
	app.AnteKeeper.SetMinimumCommission(ctx, types.DefaultMinimumCommission)

	return app, ctx
}

// SetupTest setups a new test, with new app, context, and anteHandler.
func (suite *KeeperTestSuite) SetupTest(isCheckTx bool) {
	tempDir := suite.T().TempDir()
	suite.app, suite.ctx = createTestApp(isCheckTx, tempDir)
	suite.ctx = suite.ctx.WithBlockHeight(1)

	// Set up TxConfig.
	encodingConfig := simapp.MakeTestEncodingConfig()
	// We're using TestMsg encoding in some tests, so register it here.
	encodingConfig.Amino.RegisterConcrete(&testdata.TestMsg{}, "testdata.TestMsg", nil)
	testdata.RegisterInterfaces(encodingConfig.InterfaceRegistry)

	querier := keeper.Querier{Keeper: suite.app.AnteKeeper}

	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, querier)
	suite.queryClient = types.NewQueryClient(queryHelper)
}

func (suite *KeeperTestSuite) TestParams() {
	suite.SetupTest(true) // setup
	app, ctx := suite.app, suite.ctx

	expParams := types.DefaultParams()

	//check that the empty keeper loads the default
	resParams := app.AnteKeeper.GetParams(ctx)
	suite.True(expParams.Equal(resParams))

	//modify a params, save, and retrieve
	expParams.MinimumCommissionEnforced = false
	app.AnteKeeper.SetParams(ctx, expParams)
	resParams = app.AnteKeeper.GetParams(ctx)
	suite.True(expParams.Equal(resParams))
}

func (suite *KeeperTestSuite) TestInitGenesis() {
	suite.SetupTest(true) // setup
	app, ctx := suite.app, suite.ctx

	expParams := types.DefaultParams()

	//check that the empty keeper loads the default
	app.AnteKeeper.InitGenesis(ctx, types.DefaultGenesisState())
	resParams := app.AnteKeeper.GetParams(ctx)

	suite.True(expParams.Equal(resParams))

	exportedGenesis := app.AnteKeeper.ExportGenesis(ctx)
	exportedParams := exportedGenesis.Params
	suite.True(expParams.Equal(exportedParams))
}

func (suite *KeeperTestSuite) TestGetSetMinimumCommission() {
	suite.SetupTest(true) // setup
	app, ctx := suite.app, suite.ctx

	expMinimumCommission := sdk.NewDecWithPrec(5, 2)

	//check that the empty keeper loads the default
	app.AnteKeeper.SetMinimumCommission(ctx, expMinimumCommission)
	resMinimumCommission := app.AnteKeeper.GetMinimumCommission(ctx)

	suite.True(expMinimumCommission.Equal(resMinimumCommission))
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
