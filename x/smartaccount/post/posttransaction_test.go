package post_test

import (
	"testing"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/terra-money/core/v2/x/smartaccount/post"
	"github.com/terra-money/core/v2/x/smartaccount/test_helpers"
	smartaccounttypes "github.com/terra-money/core/v2/x/smartaccount/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type PostTxTestSuite struct {
	test_helpers.SmartAccountTestSuite

	PostTxDecorator post.PostTransactionHookDecorator
	WasmKeeper      *wasmkeeper.PermissionedKeeper
}

func TestAnteSuite(t *testing.T) {
	suite.Run(t, new(PostTxTestSuite))
}

func (s *PostTxTestSuite) Setup() {
	s.SmartAccountTestSuite.SetupTests()
	s.WasmKeeper = wasmkeeper.NewDefaultPermissionKeeper(s.App.Keepers.WasmKeeper)
	s.PostTxDecorator = post.NewPostTransactionHookDecorator(s.SmartAccountKeeper, s.WasmKeeper)
	s.Ctx = s.Ctx.WithChainID("test")
}

func (s *PostTxTestSuite) TestPostTransactionHookWithoutSmartAccount() {
	s.Setup()
	txBuilder := s.BuildDefaultMsgTx(0, &types.MsgSend{
		FromAddress: s.TestAccs[0].String(),
		ToAddress:   s.TestAccs[1].String(),
		Amount:      sdk.NewCoins(sdk.NewInt64Coin("uluna", 100000000)),
	})
	_, err := s.PostTxDecorator.PostHandle(s.Ctx, txBuilder.GetTx(), false, true, sdk.ChainPostDecorators(sdk.Terminator{}))
	require.NoError(s.T(), err)
}

func (s *PostTxTestSuite) TestPostTransactionHookWithEmptySmartAccount() {
	s.Setup()
	s.Ctx = s.Ctx.WithValue(smartaccounttypes.ModuleName, smartaccounttypes.Setting{})
	txBuilder := s.BuildDefaultMsgTx(0, &types.MsgSend{
		FromAddress: s.TestAccs[0].String(),
		ToAddress:   s.TestAccs[1].String(),
		Amount:      sdk.NewCoins(sdk.NewInt64Coin("uluna", 100000000)),
	})
	_, err := s.PostTxDecorator.PostHandle(s.Ctx, txBuilder.GetTx(), false, true, sdk.ChainPostDecorators(sdk.Terminator{}))
	require.NoError(s.T(), err)
}

func (s *PostTxTestSuite) TestInvalidContractAddress() {
	s.Setup()
	s.Ctx = s.Ctx.WithValue(smartaccounttypes.ModuleName, &smartaccounttypes.Setting{
		PostTransaction: []string{s.TestAccs[0].String()},
	})
	txBuilder := s.BuildDefaultMsgTx(0, &types.MsgSend{
		FromAddress: s.TestAccs[0].String(),
		ToAddress:   s.TestAccs[1].String(),
		Amount:      sdk.NewCoins(sdk.NewInt64Coin("uluna", 100000000)),
	})
	_, err := s.PostTxDecorator.PostHandle(s.Ctx, txBuilder.GetTx(), false, true, sdk.ChainPostDecorators(sdk.Terminator{}))
	require.ErrorContainsf(s.T(), err, "no such contract", "error message: %s", err)
}

func (s *PostTxTestSuite) TestSendWithinLimitWithLimitCoinsSendHook() {
	s.Setup()

	acc := s.TestAccs[0]
	codeId, _, err := s.WasmKeeper.Create(s.Ctx, acc, test_helpers.LimitMinCoinsHookWasm, nil)
	require.NoError(s.T(), err)
	contractAddr, _, err := s.WasmKeeper.Instantiate(s.Ctx, codeId, acc, acc, []byte("{}"), "limit send", sdk.NewCoins())
	require.NoError(s.T(), err)

	s.Ctx = s.Ctx.WithValue(smartaccounttypes.ModuleName, &smartaccounttypes.Setting{
		PostTransaction: []string{contractAddr.String()},
	})

	err = s.App.Keepers.BankKeeper.SendCoinsFromAccountToModule(s.Ctx, acc, smartaccounttypes.ModuleName, sdk.NewCoins(sdk.NewInt64Coin("uluna", 10000000)))
	require.NoError(s.T(), err)

	txBuilder := s.BuildDefaultMsgTx(0, &types.MsgSend{
		FromAddress: acc.String(),
		ToAddress:   acc.String(),
		Amount:      sdk.NewCoins(sdk.NewInt64Coin("uluna", 10000000)),
	})
	_, err = s.PostTxDecorator.PostHandle(s.Ctx, txBuilder.GetTx(), false, true, sdk.ChainPostDecorators(sdk.Terminator{}))
	require.NoError(s.T(), err)
}

func (s *PostTxTestSuite) TestSendOverLimitWithLimitCoinsSendHook() {
	s.Setup()

	acc := s.TestAccs[0]
	codeId, _, err := s.WasmKeeper.Create(s.Ctx, acc, test_helpers.LimitMinCoinsHookWasm, nil)
	require.NoError(s.T(), err)
	contractAddr, _, err := s.WasmKeeper.Instantiate(s.Ctx, codeId, acc, acc, []byte("{}"), "limit send", sdk.NewCoins())
	require.NoError(s.T(), err)

	s.Ctx = s.Ctx.WithValue(smartaccounttypes.ModuleName, &smartaccounttypes.Setting{
		PostTransaction: []string{contractAddr.String()},
	})

	err = s.App.Keepers.BankKeeper.SendCoinsFromAccountToModule(s.Ctx, acc, smartaccounttypes.ModuleName, sdk.NewCoins(sdk.NewInt64Coin("uluna", 100000000)))
	require.NoError(s.T(), err)

	txBuilder := s.BuildDefaultMsgTx(0, &stakingtypes.MsgDelegate{
		DelegatorAddress: acc.String(),
		ValidatorAddress: acc.String(),
		Amount:           sdk.NewInt64Coin("uluna", 100000000),
	})
	_, err = s.PostTxDecorator.PostHandle(s.Ctx, txBuilder.GetTx(), false, true, sdk.ChainPostDecorators(sdk.Terminator{}))
	require.ErrorContainsf(s.T(), err, "Failed post transaction process", "error message: %s", err)
}

func (s *PostTxTestSuite) BuildDefaultMsgTx(accountIndex int, msgs ...sdk.Msg) client.TxBuilder {
	pk := s.TestAccPrivs[accountIndex]
	sender := s.TestAccs[accountIndex]
	acc := s.App.Keepers.AccountKeeper.GetAccount(s.Ctx, msgs[0].GetSigners()[0])
	txBuilder := s.EncodingConfig.TxConfig.NewTxBuilder()
	err := txBuilder.SetMsgs(
		msgs...,
	)
	require.NoError(s.T(), err)

	signer := authsigning.SignerData{
		Address:       sender.String(),
		ChainID:       "test",
		AccountNumber: acc.GetAccountNumber(),
		Sequence:      acc.GetSequence(),
		PubKey:        pk.PubKey(),
	}

	emptySig := signing.SignatureV2{
		PubKey: signer.PubKey,
		Data: &signing.SingleSignatureData{
			SignMode:  s.EncodingConfig.TxConfig.SignModeHandler().DefaultMode(),
			Signature: nil,
		},
		Sequence: signer.Sequence,
	}

	err = txBuilder.SetSignatures(emptySig)
	require.NoError(s.T(), err)

	sigV2, err := tx.SignWithPrivKey(
		s.EncodingConfig.TxConfig.SignModeHandler().DefaultMode(),
		signer,
		txBuilder,
		pk,
		s.EncodingConfig.TxConfig,
		acc.GetSequence(),
	)
	require.NoError(s.T(), err)

	err = txBuilder.SetSignatures(sigV2)
	require.NoError(s.T(), err)

	return txBuilder
}
