package v2_11

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	cfg module.Configurator,
	cdc codec.Codec,
	sk *stakingkeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		// Set commission based on this text proposal
		// https://station.terra.money/proposal/phoenix-1/4803
		minCommission := sdk.MustNewDecFromStr("0.05")
		stakingParams := sk.GetParams(ctx)
		stakingParams.MinCommissionRate = minCommission
		err := sk.SetParams(ctx, stakingParams)
		if err != nil {
			return nil, err
		}

		sk.IterateValidators(ctx, func(index int64, validator stakingtypes.ValidatorI) (stop bool) {
			if validator.GetCommission().LT(minCommission) {
				val := validator.(stakingtypes.Validator)
				_, err = sk.UpdateValidatorCommission(ctx, val, minCommission)
				if err != nil {
					return true
				}
			}
			return false
		})
		if err != nil {
			return nil, err
		}

		return mm.RunMigrations(ctx, cfg, fromVM)
	}
}
