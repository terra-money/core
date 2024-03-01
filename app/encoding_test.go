package app_test

import (
	"github.com/terra-money/core/v2/app"
	"github.com/terra-money/core/v2/app/test_helpers"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

type AppCodecsTestSuite struct {
	test_helpers.AppTestSuite
}

func (acts *AppCodecsTestSuite) TestCodecs() {
	// Setting up the app
	acts.Setup()

	// generating the encoding config to assert against
	encCfg := app.MakeEncodingConfig()

	// Validate the encoding config have been configured as expected for the App
	acts.Require().Equal(encCfg.Amino, acts.App.LegacyAmino())
	acts.Require().Equal(encCfg.Marshaler, acts.App.AppCodec())
	acts.Require().Equal(encCfg.InterfaceRegistry, acts.App.InterfaceRegistry())
	acts.Require().NotEmpty(acts.App.GetKey(banktypes.StoreKey))
	acts.Require().NotEmpty(acts.App.GetTKey(paramstypes.TStoreKey))
	acts.Require().NotEmpty(acts.App.GetMemKey(capabilitytypes.MemStoreKey))
}
