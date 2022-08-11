package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

// UpgradeHandler h for software upgrade proposal
type UpgradeHandler struct {
	*TerraApp
}

// NewUpgradeHandler return new instance of UpgradeHandler
func NewUpgradeHandler(app *TerraApp) UpgradeHandler {
	return UpgradeHandler{app}
}

func (h UpgradeHandler) CreateUpgradeHandler() upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		return h.mm.RunMigrations(ctx, h.configurator, vm)
	}
}
