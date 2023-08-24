package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	v2 "github.com/terra-money/core/v2/x/tokenfactory/migrations/v2"
)

type Migrator struct {
	keeper Keeper
}

func NewMigrator(keeper Keeper) Migrator {
	return Migrator{
		keeper: keeper,
	}
}

func (m Migrator) Migrate1to2(ctx sdk.Context) error {
	return v2.MigrateStore(ctx, m.keeper.paramSpace)
}
