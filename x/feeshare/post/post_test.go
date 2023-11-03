package ante_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	sdk "github.com/cosmos/cosmos-sdk/types"

	ante "github.com/terra-money/core/v2/x/feeshare/ante"
)

type AnteTestSuite struct {
	suite.Suite
}

func TestAnteSuite(t *testing.T) {
	suite.Run(t, new(AnteTestSuite))
}

func (suite *AnteTestSuite) TestFeeLogic() {
	// We expect all to pass
	feeCoins := sdk.NewCoins(sdk.NewCoin("ujuno", sdk.NewInt(500)), sdk.NewCoin("utoken", sdk.NewInt(250)))

	testCases := []struct {
		name               string
		incomingFee        sdk.Coins
		govPercent         sdk.Dec
		numContracts       int
		expectedFeePayment sdk.Coins
	}{
		{
			"100% fee / 1 contract",
			feeCoins,
			sdk.NewDecWithPrec(100, 2),
			1,
			sdk.NewCoins(sdk.NewCoin("ujuno", sdk.NewInt(500)), sdk.NewCoin("utoken", sdk.NewInt(250))),
		},
		{
			"100% fee / 2 contracts",
			feeCoins,
			sdk.NewDecWithPrec(100, 2),
			2,
			sdk.NewCoins(sdk.NewCoin("ujuno", sdk.NewInt(250)), sdk.NewCoin("utoken", sdk.NewInt(125))),
		},
		{
			"100% fee / 10 contracts",
			feeCoins,
			sdk.NewDecWithPrec(100, 2),
			10,
			sdk.NewCoins(sdk.NewCoin("ujuno", sdk.NewInt(50)), sdk.NewCoin("utoken", sdk.NewInt(25))),
		},
		{
			"67% fee / 7 contracts",
			feeCoins,
			sdk.NewDecWithPrec(67, 2),
			7,
			sdk.NewCoins(sdk.NewCoin("ujuno", sdk.NewInt(48)), sdk.NewCoin("utoken", sdk.NewInt(24))),
		},
		{
			"50% fee / 1 contracts",
			feeCoins,
			sdk.NewDecWithPrec(50, 2),
			1,
			sdk.NewCoins(sdk.NewCoin("ujuno", sdk.NewInt(250)), sdk.NewCoin("utoken", sdk.NewInt(125))),
		},
		{
			"50% fee / 2 contracts",
			feeCoins,
			sdk.NewDecWithPrec(50, 2),
			2,
			sdk.NewCoins(sdk.NewCoin("ujuno", sdk.NewInt(125)), sdk.NewCoin("utoken", sdk.NewInt(62))),
		},
		{
			"50% fee / 3 contracts",
			feeCoins,
			sdk.NewDecWithPrec(50, 2),
			3,
			sdk.NewCoins(sdk.NewCoin("ujuno", sdk.NewInt(83)), sdk.NewCoin("utoken", sdk.NewInt(42))),
		},
		{
			"25% fee / 2 contracts",
			feeCoins,
			sdk.NewDecWithPrec(25, 2),
			2,
			sdk.NewCoins(sdk.NewCoin("ujuno", sdk.NewInt(62)), sdk.NewCoin("utoken", sdk.NewInt(31))),
		},
		{
			"15% fee / 3 contracts",
			feeCoins,
			sdk.NewDecWithPrec(15, 2),
			3,
			sdk.NewCoins(sdk.NewCoin("ujuno", sdk.NewInt(25)), sdk.NewCoin("utoken", sdk.NewInt(12))),
		},
		{
			"1% fee / 2 contracts",
			feeCoins,
			sdk.NewDecWithPrec(1, 2),
			2,
			sdk.NewCoins(sdk.NewCoin("ujuno", sdk.NewInt(2)), sdk.NewCoin("utoken", sdk.NewInt(1))),
		},
	}

	for _, tc := range testCases {
		coins := ante.FeePayLogic(tc.incomingFee, tc.govPercent, tc.numContracts)

		for _, coin := range coins {
			for _, expectedCoin := range tc.expectedFeePayment {
				if coin.Denom == expectedCoin.Denom {
					suite.Require().Equal(expectedCoin.Amount.Int64(), coin.Amount.Int64(), tc.name)
				}
			}
		}
	}
}
