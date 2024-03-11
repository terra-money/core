package keeper_test

import (
	"encoding/base64"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/terra-money/core/v2/x/smartaccount/test_helpers"
	"github.com/terra-money/core/v2/x/smartaccount/types"
)

func (s *IntegrationTestSuite) TestMsgCreateAndDisableSmartAccount() {
	s.Setup()
	sender := s.TestAccs[0]

	// Ensure that the smart account was not created
	_, err := s.App.Keepers.SmartAccountKeeper.GetSetting(s.Ctx, sender.String())
	s.Require().Error(err)

	// Create a smart account
	msg := types.NewMsgCreateSmartAccount(sender.String())
	_, err = s.msgServer.CreateSmartAccount(s.Ctx, msg)
	s.Require().NoError(err)

	// Ensure that the smart account was created
	_, err = s.App.Keepers.SmartAccountKeeper.GetSetting(s.Ctx, sender.String())
	s.Require().NoError(err)

	// Ensure that the smart account cannot be created again
	_, err = s.msgServer.CreateSmartAccount(s.Ctx, msg)
	s.Require().Error(err)

	// Disable the smart account
	msgDisable := types.NewMsgDisableSmartAccount(sender.String())
	_, err = s.msgServer.DisableSmartAccount(s.Ctx, msgDisable)
	s.Require().NoError(err)

	// Ensure that the smart account was disabled
	_, err = s.App.Keepers.SmartAccountKeeper.GetSetting(s.Ctx, sender.String())
	s.Require().Error(err)
}

func (s *IntegrationTestSuite) TestMsgUpdateAuthorization() {
	s.Setup()

	// create smart account 1
	acc := s.TestAccs[1]
	msg := types.NewMsgCreateSmartAccount(acc.String())
	_, err := s.msgServer.CreateSmartAccount(s.Ctx, msg)
	s.Require().NoError(err)

	// testAcc1 using private key of testAcc0
	pubKey := s.TestAccPrivs[0].PubKey()
	// endcoding this since this should be encoded in base64 when submitted by the user
	pkEncoded := []byte(base64.StdEncoding.EncodeToString(pubKey.Bytes()))

	codeId, _, err := s.wasmKeeper.Create(s.Ctx, acc, test_helpers.SmartAuthContractWasm, nil)
	require.NoError(s.T(), err)
	contractAddr, _, err := s.wasmKeeper.Instantiate(s.Ctx, codeId, acc, acc, []byte("{}"), "auth", sdk.NewCoins())
	require.NoError(s.T(), err)

	// create updateAuth msg
	initMsg := types.Initialization{
		Account: acc.String(),
		Msg:     pkEncoded,
	}
	authMsg := &types.AuthorizationMsg{
		ContractAddress: contractAddr.String(),
		InitMsg:         &initMsg,
	}
	_, err = s.msgServer.UpdateAuthorization(s.Ctx, types.NewMsgUpdateAuthorization(
		acc.String(),
		[]*types.AuthorizationMsg{authMsg},
		false,
	))
	require.NoError(s.T(), err)

	// Ensure that the smart account was updated
	setting, err := s.App.Keepers.SmartAccountKeeper.GetSetting(s.Ctx, acc.String())
	s.Require().NoError(err)
	s.Require().Equal(acc.String(), setting.Owner)
	s.CheckAuthorizationEqual([]*types.AuthorizationMsg{authMsg}, setting.Authorization)

	// deploy another auth contract
	contractAddr2, _, err := s.wasmKeeper.Instantiate(s.Ctx, codeId, acc, acc, []byte("{}"), "auth", sdk.NewCoins())
	require.NoError(s.T(), err)

	// create updateAuth msg
	authMsg2 := &types.AuthorizationMsg{
		ContractAddress: contractAddr2.String(),
		InitMsg:         &initMsg,
	}
	_, err = s.msgServer.UpdateAuthorization(s.Ctx, types.NewMsgUpdateAuthorization(
		acc.String(),
		[]*types.AuthorizationMsg{authMsg2},
		false,
	))
	require.NoError(s.T(), err)

	// Ensure that the smart account was updated again
	setting, err = s.App.Keepers.SmartAccountKeeper.GetSetting(s.Ctx, acc.String())
	s.Require().NoError(err)
	s.Require().Equal(acc.String(), setting.Owner)
	s.CheckAuthorizationEqual([]*types.AuthorizationMsg{authMsg2}, setting.Authorization)
}

func (s *IntegrationTestSuite) TestMsgUpdateTransactionHooks() {
	s.Setup()
	sender := s.TestAccs[0]

	// Create a smart account
	msg := types.NewMsgCreateSmartAccount(sender.String())
	_, err := s.msgServer.CreateSmartAccount(s.Ctx, msg)
	s.Require().NoError(err)

	// update transaction hooks
	pretx := []string{"hook1", "hook2"}
	posttx := []string{"hook3", "hook4"}
	msgUpdate := types.NewMsgUpdateTransactionHooks(
		sender.String(),
		pretx,
		posttx,
	)
	_, err = s.msgServer.UpdateTransactionHooks(s.Ctx, msgUpdate)
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
	_, err = s.msgServer.UpdateTransactionHooks(s.Ctx, msgUpdate)
	s.Require().NoError(err)

	// Ensure that the smart account was updated again
	setting, err = s.App.Keepers.SmartAccountKeeper.GetSetting(s.Ctx, sender.String())
	s.Require().NoError(err)
	s.Require().Equal(sender.String(), setting.Owner)
	s.Require().Equal(pretx, setting.PreTransaction)
	s.Require().Equal(posttx, setting.PostTransaction)
}

func (s *IntegrationTestSuite) CheckSettingEqual(a types.Setting, b types.Setting) {
	s.Require().Equal(a.Owner, b.Owner)
	s.Require().Equal(a.Fallback, b.Fallback)
	s.CheckAuthorizationEqual(a.Authorization, b.Authorization)
	s.Require().Equal(a.PreTransaction, b.PreTransaction)
	s.Require().Equal(a.PostTransaction, b.PostTransaction)
}

func (s *IntegrationTestSuite) CheckAuthorizationEqual(a []*types.AuthorizationMsg, b []*types.AuthorizationMsg) {
	s.Require().Equal(len(a), len(b))
	for i := range a {
		s.Require().Equal(a[i].ContractAddress, b[i].ContractAddress)
		s.Require().Equal(a[i].InitMsg.Msg, b[i].InitMsg.Msg)
		s.Require().Equal(a[i].InitMsg.Account, b[i].InitMsg.Account)
	}
}
