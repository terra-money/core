package upgrade

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	antetypes "github.com/terra-money/core/v2/app/ante/types"
)

// CreateUpgradeHandler make upgrade handler
func CreateUpgradeHandler(stakingKeeper *stakingkeeper.Keeper) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		allValidators := stakingKeeper.GetAllValidators(ctx)
		for _, validator := range allValidators {
			// increase commission rate
			if validator.Commission.CommissionRates.Rate.LT(antetypes.DefaultMinimumCommission) {
				commission, err := stakingKeeper.UpdateValidatorCommission(ctx, validator, antetypes.DefaultMinimumCommission)
				if err != nil {
					return nil, err
				}

				// call the before-modification hook since we're about to update the commission
				stakingKeeper.BeforeValidatorModified(ctx, validator.GetOperator())

				validator.Commission = commission
			}

			// increase max commission rate
			if validator.Commission.CommissionRates.MaxRate.LT(antetypes.DefaultMinimumCommission) {
				validator.Commission.MaxRate = antetypes.DefaultMinimumCommission
			}

			stakingKeeper.SetValidator(ctx, validator)
		}

		return vm, nil
	}
}
