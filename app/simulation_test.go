package app_test

import (
	"os"
	"testing"

	dbm "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/stretchr/testify/require"
	"github.com/terra-money/core/v2/app"
	"github.com/terra-money/core/v2/app/wasmconfig"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simulationtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	simcli "github.com/cosmos/cosmos-sdk/x/simulation/client/cli"
)

func init() {
	simcli.GetSimulatorFlags()
}

type AppTest interface {
	app.TerraApp
	GetBaseApp() *baseapp.BaseApp
	AppCodec() codec.Codec
	SimulationManager() *module.SimulationManager
	ModuleAccountAddrs() map[string]bool
	Name() string
	LegacyAmino() *codec.LegacyAmino
	BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock
	EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock
	InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain
}

// BenchmarkSimulation run the chain simulation
// Running using starport command:
// `starport chain simulate -v --numBlocks 200 --blockSize 50`
// Running as go benchmark test:
// `go test -benchmem -run=^$ -bench ^BenchmarkSimulation ./app -NumBlocks=200 -BlockSize 50 -Commit=true -Verbose=true -Enabled=true`
func BenchmarkSimulation(b *testing.B) {
	config := simcli.NewConfigFromFlags()
	simcli.FlagEnabledValue = true
	simcli.FlagCommitValue = true

	db, dir, logger, _, err := simtestutil.SetupSimulation(config, "goleveldb-app-sim", "Simulation", true, false)
	require.NoError(b, err, "simulation setup failed")

	b.Cleanup(func() {
		db.Close()
		err = os.RemoveAll(dir)
		require.NoError(b, err)
	})

	encoding := app.MakeEncodingConfig()

	AppTest := app.NewTerraApp(
		logger,
		db,
		nil,
		true,
		map[int64]bool{},
		app.DefaultNodeHome,
		0,
		encoding,
		simtestutil.EmptyAppOptions{},
		wasmconfig.DefaultConfig(),
	)

	// Run randomized simulations
	_, simParams, simErr := simulation.SimulateFromSeed(
		b,
		os.Stdout,
		AppTest.BaseApp,
		simtestutil.AppStateFn(AppTest.AppCodec(), AppTest.SimulationManager(), AppTest.DefaultGenesis()),
		simulationtypes.RandomAccounts,
		simtestutil.SimulationOperations(AppTest, AppTest.AppCodec(), config),
		AppTest.ModuleAccountAddrs(),
		config,
		AppTest.AppCodec(),
	)

	// export state and simParams before the simulation error is checked
	err = simtestutil.CheckExportSimulation(AppTest, config, simParams)
	require.NoError(b, err)
	require.NoError(b, simErr)

	if config.Commit {
		simtestutil.PrintStats(db)
	}
}

func TestSimulationManager(t *testing.T) {
	db := dbm.NewMemDB()
	encoding := app.MakeEncodingConfig()

	AppTest := app.NewTerraApp(
		log.NewTMLogger(log.NewSyncWriter(os.Stdout)),
		db,
		nil,
		true,
		map[int64]bool{},
		app.DefaultNodeHome,
		0,
		encoding,
		simtestutil.EmptyAppOptions{},
		wasmconfig.DefaultConfig(),
	)
	sm := AppTest.SimulationManager()
	require.NotNil(t, sm)
}
