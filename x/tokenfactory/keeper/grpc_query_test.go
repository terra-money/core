package keeper_test

import "github.com/terra-money/core/v2/x/tokenfactory/types"

func (s *KeeperTestSuite) TestQueryParams() {
	params, err := s.App.TokenFactoryKeeper.Params(s.Ctx, &types.QueryParamsRequest{})

	s.Require().NoError(err)

	expected := types.QueryParamsResponse{
		Params: s.App.TokenFactoryKeeper.GetParams(s.Ctx),
	}
	s.Require().Equal(&expected, params)

}

func (s *KeeperTestSuite) TestQueryBeforeSendHookEmptyAddress() {
	res, err := s.App.TokenFactoryKeeper.BeforeSendHookAddress(s.Ctx, &types.QueryBeforeSendHookAddressRequest{})

	s.Require().NoError(err)

	expected := types.QueryBeforeSendHookAddressResponse{
		CosmwasmAddress: "",
	}
	s.Require().Equal(&expected, res)

}

func (s *KeeperTestSuite) TestQueryBeforeSendHookNonRegisteredAddress() {
	res, err := s.App.TokenFactoryKeeper.BeforeSendHookAddress(s.Ctx, &types.QueryBeforeSendHookAddressRequest{
		Denom: "bitcoin",
	})
	s.Require().NoError(err)

	expected := types.QueryBeforeSendHookAddressResponse{
		CosmwasmAddress: "",
	}
	s.Require().Equal(&expected, res)

}
