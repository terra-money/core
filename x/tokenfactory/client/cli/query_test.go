package cli_test

import (
	gocontext "context"
	"testing"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/suite"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/terra-money/core/v2/app/app_testing"
	"github.com/terra-money/core/v2/app/config"
	"github.com/terra-money/core/v2/x/tokenfactory/types"
)

type QueryTestSuite struct {
	app_testing.AppTestSuite
}

func (s *QueryTestSuite) TestQueriesNeverAlterState() {
	s.Setup()

	// fund acc
	fundAccsAmount := sdk.NewCoins(sdk.NewCoin(config.BondDenom, math.NewInt(1_000_000_000)))
	s.FundAcc(s.TestAccs[0], fundAccsAmount)
	// create new token
	_, err := s.App.TokenFactoryKeeper.CreateDenom(s.Ctx, s.TestAccs[0].String(), "tokenfactory")
	s.Require().NoError(err)

	testCases := []struct {
		name   string
		query  string
		input  interface{}
		output interface{}
	}{
		{
			"Query denom authority metadata",
			"/osmosis.tokenfactory.v1beta1.Query/DenomAuthorityMetadata",
			&types.QueryDenomAuthorityMetadataRequest{Denom: "tokenfactory"},
			&types.QueryDenomAuthorityMetadataResponse{},
		},
		{
			"Query denoms by creator",
			"/osmosis.tokenfactory.v1beta1.Query/DenomsFromCreator",
			&types.QueryDenomsFromCreatorRequest{Creator: s.TestAccs[0].String()},
			&types.QueryDenomsFromCreatorResponse{},
		},
		{
			"Query params",
			"/osmosis.tokenfactory.v1beta1.Query/Params",
			&types.QueryParamsRequest{},
			&types.QueryParamsResponse{},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			err := s.QueryHelper.Invoke(gocontext.Background(), tc.query, tc.input, tc.output)
			s.Require().NoError(err)
		})
	}
}

func TestQueryTestSuite(t *testing.T) {
	suite.Run(t, new(QueryTestSuite))
}
