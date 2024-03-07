package keeper_test

import (
	smartaccountkeeper "github.com/terra-money/core/v2/x/smartaccount/keeper"

	"github.com/terra-money/core/v2/x/smartaccount/types"
)

func (s *IntegrationTestSuite) TestMsgCreateAndDisableSmartAccount() {
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

	// Disable the smart account
	msgDisable := types.NewMsgDisableSmartAccount(sender.String())
	_, err = ms.DisableSmartAccount(s.Ctx, msgDisable)
	s.Require().NoError(err)

	// Ensure that the smart account was disabled
	_, err = s.App.Keepers.SmartAccountKeeper.GetSetting(s.Ctx, sender.String())
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
		InitMsg:         &types.Initialization{},
	}
	msgUpdate := types.NewMsgUpdateAuthorization(sender.String(), []*types.AuthorizationMsg{&authorization}, true)
	_, err = ms.UpdateAuthorization(s.Ctx, msgUpdate)
	s.Require().NoError(err)

	// Ensure that the smart account was updated
	setting, err := s.App.Keepers.SmartAccountKeeper.GetSetting(s.Ctx, sender.String())
	s.Require().NoError(err)
	s.Require().Equal(sender.String(), setting.Owner)
	s.Require().Equal([]*types.AuthorizationMsg{&authorization}, setting.Authorization)

	// update authorization again
	authorization2 := types.AuthorizationMsg{
		ContractAddress: "bbc",
		InitMsg:         &types.Initialization{},
	}
	msgUpdate2 := types.NewMsgUpdateAuthorization(sender.String(), []*types.AuthorizationMsg{&authorization2}, true)
	_, err = ms.UpdateAuthorization(s.Ctx, msgUpdate2)
	s.Require().NoError(err)

	// Ensure that the smart account was updated again
	setting, err = s.App.Keepers.SmartAccountKeeper.GetSetting(s.Ctx, sender.String())
	s.Require().NoError(err)
	s.Require().Equal(sender.String(), setting.Owner)
	s.Require().Equal([]*types.AuthorizationMsg{&authorization2}, setting.Authorization)
}

func (s *IntegrationTestSuite) TestMsgUpdateTransactionHooks() {
	s.Setup()
	sender := s.TestAccs[0]

	// Create a smart account
	ms := smartaccountkeeper.NewMsgServer(s.App.Keepers.SmartAccountKeeper)
	msg := types.NewMsgCreateSmartAccount(sender.String())
	_, err := ms.CreateSmartAccount(s.Ctx, msg)
	s.Require().NoError(err)

	// update transaction hooks
	pretx := []string{"hook1", "hook2"}
	posttx := []string{"hook3", "hook4"}
	msgUpdate := types.NewMsgUpdateTransactionHooks(
		sender.String(),
		pretx,
		posttx,
	)
	_, err = ms.UpdateTransactionHooks(s.Ctx, msgUpdate)
	s.Require().NoError(err)

	// Ensure that the smart account was updated
	setting, err := s.App.Keepers.SmartAccountKeeper.GetSetting(s.Ctx, sender.String())
	s.Require().NoError(err)
	s.Require().Equal(sender.String(), setting.Owner)
	s.Require().Equal(pretx, setting.PreTransaction)
	s.Require().Equal(posttx, setting.PostTransaction)

	// update authorization again
	pretx = []string{"hook5", "hook6"}
	posttx = []string{"hook7", "hook8"}
	msgUpdate = types.NewMsgUpdateTransactionHooks(
		sender.String(),
		pretx,
		posttx,
	)
	_, err = ms.UpdateTransactionHooks(s.Ctx, msgUpdate)
	s.Require().NoError(err)

	// Ensure that the smart account was updated again
	setting, err = s.App.Keepers.SmartAccountKeeper.GetSetting(s.Ctx, sender.String())
	s.Require().NoError(err)
	s.Require().Equal(sender.String(), setting.Owner)
	s.Require().Equal(pretx, setting.PreTransaction)
	s.Require().Equal(posttx, setting.PostTransaction)
}
