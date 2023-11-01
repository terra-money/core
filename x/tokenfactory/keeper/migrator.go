package keeper

import (
	"github.com/terra-money/core/v2/x/tokenfactory/exported"
	v2 "github.com/terra-money/core/v2/x/tokenfactory/migrations/v2"
	v3 "github.com/terra-money/core/v2/x/tokenfactory/migrations/v3"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Migrator struct {
	keeper         Keeper
	legacySubspace exported.Subspace
}

func NewMigrator(keeper Keeper, legacySubspace exported.Subspace) Migrator {
	return Migrator{
		keeper:         keeper,
		legacySubspace: legacySubspace,
	}
}

func (m Migrator) Migrate1to2(ctx sdk.Context) error {
	return v2.MigrateStore(ctx, m.legacySubspace, m.keeper.cdc)
}

func (m Migrator) Migrate2to3(ctx sdk.Context) error {
	return v3.MigrateStore(ctx, m.keeper.storeKey, m.legacySubspace, m.keeper.cdc)
}
