package bindings_test

import (
	"os"
	"testing"
	"time"

	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/stretchr/testify/require"

	"github.com/cometbft/cometbft/crypto"
	"github.com/cometbft/cometbft/crypto/ed25519"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"

	dbm "github.com/cometbft/cometbft-db"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/terra-money/core/v2/app"
	tokenfactorytypes "github.com/terra-money/core/v2/x/tokenfactory/types"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
)

func CreateTestInput() (*app.TerraApp, sdk.Context) {
	encCfg := app.MakeEncodingConfig()
	genesisState := app.NewDefaultGenesisState(encCfg.Marshaler)
	genesisState.SetDefaultTerraConfig(encCfg.Marshaler)
	db := dbm.NewMemDB()
	terraApp := app.NewTerraApp(
		log.NewTMLogger(log.NewSyncWriter(os.Stdout)),
		db,
		nil,
		true,
		map[int64]bool{},
		app.DefaultNodeHome,
		0,
		encCfg,
		simtestutil.EmptyAppOptions{},
		wasmtypes.DefaultWasmConfig(),
	)
	ctx := terraApp.BaseApp.NewContext(true, tmproto.Header{Height: 1, ChainID: "phoenix-1", Time: time.Now()})
	err := terraApp.Keepers.WasmKeeper.SetParams(ctx, wasmtypes.DefaultParams())
	if err != nil {
		panic(err)
	}
	terraApp.Keepers.BankKeeper.SetParams(ctx, banktypes.NewParams(true))
	if err != nil {
		panic(err)
	}
	terraApp.Keepers.TokenFactoryKeeper.SetParams(ctx, tokenfactorytypes.DefaultParams())
	if err != nil {
		panic(err)
	}
	terraApp.Keepers.DistrKeeper.SetFeePool(ctx, distrtypes.InitialFeePool())
	if err != nil {
		panic(err)
	}
	return terraApp, ctx
}

func FundAccount(t *testing.T, ctx sdk.Context, terra *app.TerraApp, acct sdk.AccAddress) {
	t.Helper()
}

// we need to make this deterministic (same every test run), as content might affect gas costs
func keyPubAddr() (crypto.PrivKey, crypto.PubKey, sdk.AccAddress) {
	key := ed25519.GenPrivKey()
	pub := key.PubKey()
	addr := sdk.AccAddress(pub.Address())
	return key, pub, addr
}

func RandomAccountAddress() sdk.AccAddress {
	_, _, addr := keyPubAddr()
	return addr
}

func RandomBech32AccountAddress() string {
	return RandomAccountAddress().String()
}

func storeReflectCode(t *testing.T, ctx sdk.Context, app *app.TerraApp, addr sdk.AccAddress) uint64 {
	t.Helper()
	wasmCode, err := os.ReadFile("./testdata/token_reflect.wasm")
	require.NoError(t, err)

	contractKeeper := keeper.NewDefaultPermissionKeeper(app.Keepers.WasmKeeper)
	codeID, _, err := contractKeeper.Create(ctx, addr, wasmCode, &wasmtypes.AccessConfig{Permission: wasmtypes.AccessTypeEverybody})
	require.NoError(t, err)

	return codeID
}

func instantiateReflectContract(t *testing.T, ctx sdk.Context, app *app.TerraApp, funder sdk.AccAddress) sdk.AccAddress {
	t.Helper()
	initMsgBz := []byte("{}")
	contractKeeper := keeper.NewDefaultPermissionKeeper(app.Keepers.WasmKeeper)
	codeID := uint64(1)
	addr, _, err := contractKeeper.Instantiate(ctx, codeID, funder, funder, initMsgBz, "demo contract", nil)
	require.NoError(t, err)

	return addr
}

func fundAccount(t *testing.T, ctx sdk.Context, app *app.TerraApp, addr sdk.AccAddress, coins sdk.Coins) {
	t.Helper()
	err := app.Keepers.BankKeeper.MintCoins(ctx, minttypes.ModuleName, coins)
	require.NoError(t, err)
	err = app.Keepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, addr, coins)
	require.NoError(t, err)
}

func SetupCustomApp(t *testing.T, addr sdk.AccAddress) (*app.TerraApp, sdk.Context) {
	t.Helper()
	app, ctx := CreateTestInput()

	storeReflectCode(t, ctx, app, addr)

	cInfo := app.Keepers.WasmKeeper.GetCodeInfo(ctx, 1)
	require.NotNil(t, cInfo)

	return app, ctx
}
