package keeper_test

import (
	"github.com/terra-money/core/v2/x/tokenfactory/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func (s *KeeperTestSuite) TestGenesis() {
	genesisState := types.GenesisState{
		FactoryDenoms: []types.GenesisDenom{
			{
				Denom: "factory/terra13s4gwzxv6dycfctvddfuy6r3zm7d6zklynzzj5/bitcoin",
				AuthorityMetadata: types.DenomAuthorityMetadata{
					Admin: "terra13s4gwzxv6dycfctvddfuy6r3zm7d6zklynzzj5",
				},
			},
			{
				Denom: "factory/terra13s4gwzxv6dycfctvddfuy6r3zm7d6zklynzzj5/diff-admin",
				AuthorityMetadata: types.DenomAuthorityMetadata{
					Admin: "terra16jpsrgl423fqg6n0e9edllew9z0gm7rhl5300u",
				},
			},
			{
				Denom: "factory/terra13s4gwzxv6dycfctvddfuy6r3zm7d6zklynzzj5/litecoin",
				AuthorityMetadata: types.DenomAuthorityMetadata{
					Admin: "terra13s4gwzxv6dycfctvddfuy6r3zm7d6zklynzzj5",
				},
			},
		},
	}

	app := s.App

	// Test both with bank denom metadata set, and not set.
	for i, denom := range genesisState.FactoryDenoms {
		// hacky, sets bank metadata to exist if i != 0, to cover both cases.
		if i != 0 {
			app.Keepers.BankKeeper.SetDenomMetaData(s.Ctx, banktypes.Metadata{Base: denom.GetDenom(), Display: "test"})
		}
	}

	// check before initGenesis that the module account is nil
	tokenfactoryModuleAccount := app.Keepers.AccountKeeper.GetAccount(s.Ctx, app.Keepers.AccountKeeper.GetModuleAddress(types.ModuleName))
	s.Require().Nil(tokenfactoryModuleAccount)

	app.Keepers.TokenFactoryKeeper.SetParams(s.Ctx, types.Params{DenomCreationFee: sdk.Coins{sdk.NewInt64Coin("uosmo", 100)}})
	app.Keepers.TokenFactoryKeeper.InitGenesis(s.Ctx, genesisState)

	// check that the module account is now initialized
	tokenfactoryModuleAccount = app.Keepers.AccountKeeper.GetAccount(s.Ctx, app.Keepers.AccountKeeper.GetModuleAddress(types.ModuleName))
	s.Require().NotNil(tokenfactoryModuleAccount)

	exportedGenesis := app.Keepers.TokenFactoryKeeper.ExportGenesis(s.Ctx)
	s.Require().NotNil(exportedGenesis)
	s.Require().Equal(genesisState, *exportedGenesis)

	app.Keepers.BankKeeper.SetParams(s.Ctx, banktypes.DefaultParams())
	app.Keepers.BankKeeper.InitGenesis(s.Ctx, app.Keepers.BankKeeper.ExportGenesis(s.Ctx))
	for i, denom := range genesisState.FactoryDenoms {
		// hacky, check whether bank metadata is not replaced if i != 0, to cover both cases.
		if i != 0 {
			metadata, found := app.Keepers.BankKeeper.GetDenomMetaData(s.Ctx, denom.GetDenom())
			s.Require().True(found)
			s.Require().Equal(metadata, banktypes.Metadata{Base: denom.GetDenom(), Display: "test"})
		}
	}
}
