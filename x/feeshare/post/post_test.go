package ante_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	sdk "github.com/cosmos/cosmos-sdk/types"

	app "github.com/terra-money/core/v2/app/app_test"
	post "github.com/terra-money/core/v2/x/feeshare/post"
	"github.com/terra-money/core/v2/x/feeshare/types"
)

type AnteTestSuite struct {
	app.AppTestSuite
}

func TestAnteSuite(t *testing.T) {
	suite.Run(t, new(AnteTestSuite))
}

func (suite *AnteTestSuite) TestGetWithdrawalAddressFromContract() {
	suite.Setup()

	feeshareKeeper := suite.AppTestSuite.App.Keepers.FeeShareKeeper
	feeshareKeeper.SetFeeShare(suite.Ctx, types.FeeShare{
		ContractAddress:   "terra1jwyzzsaag4t0evnuukc35ysyrx9arzdde2kg9cld28alhjurtthq0prs2s",
		DeployerAddress:   "",
		WithdrawerAddress: "terra1zdpgj8am5nqqvht927k3etljyl6a52kwqup0je",
	})
	feeshareKeeper.SetFeeShare(suite.Ctx, types.FeeShare{
		ContractAddress:   "terra1mdpvgjc8jmv60a4x68nggsh9w8uyv69sqls04a76m9med5hsqmwsse8sxa",
		DeployerAddress:   "",
		WithdrawerAddress: "",
	})

	testCases := []struct {
		name                    string
		contractAddresses       []string
		expectedWithdrawerAddrs []sdk.AccAddress
		expectErr               bool
	}{
		{
			"valid contract addresses",
			[]string{"terra1jwyzzsaag4t0evnuukc35ysyrx9arzdde2kg9cld28alhjurtthq0prs2s"},
			[]sdk.AccAddress{
				sdk.MustAccAddressFromBech32("terra1zdpgj8am5nqqvht927k3etljyl6a52kwqup0je"),
			},
			false,
		},
		{
			"without withdrawer contract addresses",
			[]string{"terra1mdpvgjc8jmv60a4x68nggsh9w8uyv69sqls04a76m9med5hsqmwsse8sxa"},
			[]sdk.AccAddress(nil),
			false,
		},
		{
			"invalid contract address",
			[]string{"invalidAddress"},
			nil,
			true,
		},
	}

	for _, tc := range testCases {
		withdrawerAddrs, err := post.GetWithdrawalAddressFromContract(
			suite.Ctx,
			tc.contractAddresses,
			feeshareKeeper,
		)

		if tc.expectErr {
			suite.Require().Error(err, tc.name)
		} else {
			suite.Require().NoError(err, tc.name)
			suite.Require().Equal(tc.expectedWithdrawerAddrs, withdrawerAddrs, tc.name)
		}
	}
}

func (suite *AnteTestSuite) TestCalculateFee() {
	feeCoins := sdk.NewCoins(sdk.NewCoin("uluna", sdk.NewInt(500)), sdk.NewCoin("utoken", sdk.NewInt(250)))

	testCases := []struct {
		name               string
		incomingFee        sdk.Coins
		govPercent         sdk.Dec
		numContracts       int
		allowdDenoms       []string
		expectedFeePayment sdk.Coins
	}{
		{
			"100% fee / 1 contract",
			feeCoins,
			sdk.NewDecWithPrec(100, 2),
			1,
			[]string{},
			sdk.NewCoins(sdk.NewCoin("uluna", sdk.NewInt(500)), sdk.NewCoin("utoken", sdk.NewInt(250))),
		},
		{
			"100% fee / 2 contracts",
			feeCoins,
			sdk.NewDecWithPrec(100, 2),
			2,
			[]string{},
			sdk.NewCoins(sdk.NewCoin("uluna", sdk.NewInt(250)), sdk.NewCoin("utoken", sdk.NewInt(125))),
		},
		{
			"100% fee / 10 contracts / 1 allowed denom",
			feeCoins,
			sdk.NewDecWithPrec(100, 2),
			10,
			[]string{"uluna"},
			sdk.NewCoins(sdk.NewCoin("uluna", sdk.NewInt(50))),
		},
		{
			"67% fee / 7 contracts",
			feeCoins,
			sdk.NewDecWithPrec(67, 2),
			7,
			[]string{},
			sdk.NewCoins(sdk.NewCoin("uluna", sdk.NewInt(48)), sdk.NewCoin("utoken", sdk.NewInt(24))),
		},
		{
			"50% fee / 1 contracts / 1 allowed denom",
			feeCoins,
			sdk.NewDecWithPrec(50, 2),
			1,
			[]string{"utoken"},
			sdk.NewCoins(sdk.NewCoin("utoken", sdk.NewInt(125))),
		},
		{
			"50% fee / 2 contracts / 2 allowed denoms",
			feeCoins,
			sdk.NewDecWithPrec(50, 2),
			2,
			[]string{"uluna", "utoken"},
			sdk.NewCoins(sdk.NewCoin("uluna", sdk.NewInt(125)), sdk.NewCoin("utoken", sdk.NewInt(62))),
		},
		{
			"50% fee / 3 contracts",
			feeCoins,
			sdk.NewDecWithPrec(50, 2),
			3,
			[]string{},
			sdk.NewCoins(sdk.NewCoin("uluna", sdk.NewInt(83)), sdk.NewCoin("utoken", sdk.NewInt(42))),
		},
		{
			"25% fee / 2 contracts",
			feeCoins,
			sdk.NewDecWithPrec(25, 2),
			2,
			[]string{},
			sdk.NewCoins(sdk.NewCoin("uluna", sdk.NewInt(62)), sdk.NewCoin("utoken", sdk.NewInt(31))),
		},
		{
			"15% fee / 3 contracts / inexistent denom",
			feeCoins,
			sdk.NewDecWithPrec(15, 2),
			3,
			[]string{"ubtc"},
			sdk.Coins(nil),
		},
		{
			"1% fee / 2 contracts",
			feeCoins,
			sdk.NewDecWithPrec(1, 2),
			2,
			[]string{},
			sdk.NewCoins(sdk.NewCoin("uluna", sdk.NewInt(2)), sdk.NewCoin("utoken", sdk.NewInt(1))),
		},
	}

	for _, tc := range testCases {
		feeToBePaid := post.CalculateFee(tc.incomingFee, tc.govPercent, tc.numContracts, tc.allowdDenoms)

		suite.Require().Equal(tc.expectedFeePayment, feeToBePaid, tc.name)
	}
}
