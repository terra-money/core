package app_test

import (
	"os"
	"testing"

	dbm "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"

	"github.com/stretchr/testify/require"
	"github.com/terra-money/core/v2/app"
	"github.com/terra-money/core/v2/app/keepers"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"

	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	simulationtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	simcli "github.com/cosmos/cosmos-sdk/x/simulation/client/cli"
)

func init() {
	simcli.GetSimulatorFlags()
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
	enabled := simcli.FlagEnabledValue

	db, dir, logger, _, err := simtestutil.SetupSimulation(config, "goleveldb-app-sim", "Simulation", true, enabled)
	require.NoError(b, err, "simulation setup failed")

	b.Cleanup(func() {
		db.Close()
		err = os.RemoveAll(dir)
		require.NoError(b, err)
	})

	encoding := app.MakeEncodingConfig()

	terraApp := app.NewTerraApp(
		logger,
		db,
		nil,
		true,
		map[int64]bool{},
		app.DefaultNodeHome,
		0,
		encoding,
		simtestutil.EmptyAppOptions{},
		wasmtypes.DefaultWasmConfig(),
	)

	// Run randomized simulations
	_, simParams, simErr := simulation.SimulateFromSeed(
		b,
		os.Stdout,
		terraApp.BaseApp,
		simtestutil.AppStateFn(terraApp.AppCodec(), terraApp.SimulationManager(), terraApp.DefaultGenesis()),
		simulationtypes.RandomAccounts,
		simtestutil.SimulationOperations(terraApp, terraApp.AppCodec(), config),
		keepers.ModuleAccountAddrs(),
		config,
		terraApp.AppCodec(),
	)

	// export state and simParams before the simulation error is checked
	err = simtestutil.CheckExportSimulation(terraApp, config, simParams)
	require.NoError(b, err)
	require.NoError(b, simErr)

	if config.Commit {
		simtestutil.PrintStats(db)
	}
}

func TestSimulationManager(t *testing.T) {
	db := dbm.NewMemDB()
	encoding := app.MakeEncodingConfig()

	terraApp := app.NewTerraApp(
		log.NewTMLogger(log.NewSyncWriter(os.Stdout)),
		db,
		nil,
		true,
		map[int64]bool{},
		app.DefaultNodeHome,
		0,
		encoding,
		simtestutil.EmptyAppOptions{},
		wasmtypes.DefaultWasmConfig(),
	)
	sm := terraApp.SimulationManager()
	require.NotNil(t, sm)
}
