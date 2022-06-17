package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	antetypes "github.com/terra-money/core/v2/app/ante/types"
)

// CreateUpgradeHandler make upgrade handler
func (app TerraApp) CreateUpgradeHandler() upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		minimumCommission := antetypes.DefaultMinimumCommission
		allValidators := app.StakingKeeper.GetAllValidators(ctx)
		for _, validator := range allValidators {
			// increase commission rate
			if validator.Commission.CommissionRates.Rate.LT(minimumCommission) {

				// call the before-modification hook since we're about to update the commission
				app.StakingKeeper.BeforeValidatorModified(ctx, validator.GetOperator())

				validator.Commission.Rate = minimumCommission
				validator.Commission.UpdateTime = ctx.BlockHeader().Time
			}

			// increase max commission rate
			if validator.Commission.CommissionRates.MaxRate.LT(minimumCommission) {
				validator.Commission.MaxRate = minimumCommission
			}

			app.StakingKeeper.SetValidator(ctx, validator)
		}

		return app.mm.RunMigrations(ctx, app.configurator, vm)
	}
}
