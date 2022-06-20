package keeper_test

import (
	gocontext "context"
	"fmt"

	"github.com/terra-money/core/v2/app/ante/types"
)

func (suite *KeeperTestSuite) TestGRPCQueryParams() {
	suite.SetupTest(true) // setup
	queryClient := suite.queryClient

	var req *types.QueryParamsRequest
	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"valid request",
			func() {
				req = &types.QueryParamsRequest{}
			},
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()
			res, err := queryClient.Params(gocontext.Background(), req)
			if tc.expPass {
				suite.NoError(err)
				suite.True(res.Params.Equal(types.DefaultParams()))
			} else {
				suite.Error(err)
				suite.Nil(res)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCQueryMinimumCommission() {
	suite.SetupTest(true) // setup
	queryClient := suite.queryClient

	var req *types.QueryMinimumCommissionRequest
	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"valid request",
			func() {
				req = &types.QueryMinimumCommissionRequest{}
			},
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()
			res, err := queryClient.MinimumCommission(gocontext.Background(), req)
			if tc.expPass {
				suite.NoError(err)
				suite.True(res.MinimumCommission.Equal(types.DefaultMinimumCommission))
			} else {
				suite.Error(err)
				suite.Nil(res)
			}
		})
	}
}
