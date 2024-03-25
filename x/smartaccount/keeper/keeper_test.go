package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/terra-money/core/v2/x/smartaccount/keeper"
	"github.com/terra-money/core/v2/x/smartaccount/test_helpers"
	"github.com/terra-money/core/v2/x/smartaccount/types"
)

type IntegrationTestSuite struct {
	test_helpers.SmartAccountTestSuite
	msgServer  types.MsgServer
	wasmKeeper *wasmkeeper.PermissionedKeeper
}

func (s *IntegrationTestSuite) Setup() {
	s.SmartAccountTestSuite.SetupTests()
	s.msgServer = keeper.NewMsgServer(s.SmartAccountKeeper)
	s.wasmKeeper = wasmkeeper.NewDefaultPermissionKeeper(s.App.Keepers.WasmKeeper)
	s.Ctx = s.Ctx.WithChainID("test")
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
