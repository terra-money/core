package test_helpers

import (
	"github.com/terra-money/core/v2/app/test_helpers"
	"github.com/terra-money/core/v2/x/smartaccount/keeper"
)

type SmartAccountTestSuite struct {
	test_helpers.AppTestSuite

	SmartAccountKeeper keeper.Keeper
}

func (s *SmartAccountTestSuite) SetupTests() {
	s.Setup()
	s.SmartAccountKeeper = s.App.Keepers.SmartAccountKeeper
}
