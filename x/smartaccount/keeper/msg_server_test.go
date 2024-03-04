package keeper_test

import (
	smartaccountkeeper "github.com/terra-money/core/v2/x/smartaccount/keeper"

	"github.com/terra-money/core/v2/x/smartaccount/types"
)

func (s *IntegrationTestSuite) TestMsgCreateSmartAccount() {
	s.Setup()
	sender := s.TestAccs[0]

	// Ensure that the smart account was not created
	_, err := s.App.Keepers.SmartAccountKeeper.GetSetting(s.Ctx, sender.String())
	s.Require().Error(err)

	// Create a smart account
	ms := smartaccountkeeper.NewMsgServer(s.App.Keepers.SmartAccountKeeper)
	msg := types.NewMsgCreateSmartAccount(sender.String())
	_, err = ms.CreateSmartAccount(s.Ctx, msg)
	s.Require().NoError(err)

	// Ensure that the smart account was created
	_, err = s.App.Keepers.SmartAccountKeeper.GetSetting(s.Ctx, sender.String())
	s.Require().NoError(err)

	// Ensure that the smart account cannot be created again
	_, err = ms.CreateSmartAccount(s.Ctx, msg)
	s.Require().Error(err)
}

func (s *IntegrationTestSuite) TestMsgUpdateAuthorization() {
	s.Setup()
	sender := s.TestAccs[0]

	// Create a smart account
	ms := smartaccountkeeper.NewMsgServer(s.App.Keepers.SmartAccountKeeper)
	msg := types.NewMsgCreateSmartAccount(sender.String())
	_, err := ms.CreateSmartAccount(s.Ctx, msg)
	s.Require().NoError(err)

	// update authorization
	authorization := types.AuthorizationMsg{
		ContractAddress: "abc",
		InitMsg:         "abc",
	}
	msgUpdate := types.NewMsgUpdateAuthorization(sender.String(), []types.AuthorizationMsg{&authorization}, true)
	_, err = ms.UpdateAuthorization(s.Ctx, msgUpdate)
	s.Require().NoError(err)

	// Ensure that the smart account was updated
	setting, err := s.App.Keepers.SmartAccountKeeper.GetSetting(s.Ctx, sender.String())
	s.Require().NoError(err)
	s.Require().Equal(sender.String(), setting.Owner)
	s.Require().Equal([]*types.AuthorizationMsg{&authorization}, *setting.Authorization[0])
}
