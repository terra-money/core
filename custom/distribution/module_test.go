package distribution_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	sdk "github.com/cosmos/cosmos-sdk/types"

	sdkerrors "cosmossdk.io/errors"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	app "github.com/terra-money/core/v2/app/app_test"
)

type DistributionModuleTest struct {
	app.AppTestSuite
}

func TestDistributionModuleTest(t *testing.T) {
	suite.Run(t, new(DistributionModuleTest))
}

func (s *DistributionModuleTest) TestDistributonModule() {
	s.Setup()
	receiver := s.TestAccs[0]

	testCases := []struct {
		name       string
		expectPass bool
		err        error
		blockTime  time.Time
	}{
		{
			"Send transasction from community pool",
			true,
			nil,
			time.Date(2025, 1, 1, 1, 1, 1, 1, time.UTC),
		},
		{
			"Block transasction from community pool",
			false,
			sdkerrors.New(distrtypes.ModuleName, 999, "CommunityPool is blocked until 2025-01-01 00:00:00 +0000 UTC"),
			time.Date(2024, 1, 1, 1, 1, 1, 1, time.UTC),
		},
	}

	for _, tc := range testCases {
		s.Ctx = s.Ctx.WithBlockTime(tc.blockTime)

		err := s.App.DistrKeeper.DistributeFromFeePool(
			s.Ctx,
			sdk.NewCoins(sdk.NewCoin("uluna", sdk.NewInt(0))),
			receiver,
		)
		if tc.expectPass {
			s.Require().NoError(err, tc.name)
		} else {
			s.Require().EqualError(err, tc.err.Error(), tc.name)
		}
	}
}
