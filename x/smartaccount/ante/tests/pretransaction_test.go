package tests

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"
	"github.com/terra-money/core/v2/x/smartaccount/test_helpers"
	smartaccounttypes "github.com/terra-money/core/v2/x/smartaccount/types"
)

func (s *AnteTestSuite) TestPreTransactionHookWithoutSmartAccount() {
	s.Setup()
	txBuilder := s.BuildDefaultMsgTx(0, &types.MsgSend{
		FromAddress: s.TestAccs[0].String(),
		ToAddress:   s.TestAccs[1].String(),
		Amount:      sdk.NewCoins(sdk.NewInt64Coin("uluna", 100000000)),
	})
	_, err := s.PreTxDecorator.AnteHandle(s.Ctx, txBuilder.GetTx(), false, sdk.ChainAnteDecorators(sdk.Terminator{}))
	require.NoError(s.T(), err)
}

func (s *AnteTestSuite) TestPreTransactionHookWithEmptySmartAccount() {
	s.Setup()
	s.Ctx = s.Ctx.WithValue(smartaccounttypes.ModuleName, smartaccounttypes.Setting{})
	txBuilder := s.BuildDefaultMsgTx(0, &types.MsgSend{
		FromAddress: s.TestAccs[0].String(),
		ToAddress:   s.TestAccs[1].String(),
		Amount:      sdk.NewCoins(sdk.NewInt64Coin("uluna", 100000000)),
	})
	_, err := s.PreTxDecorator.AnteHandle(s.Ctx, txBuilder.GetTx(), false, sdk.ChainAnteDecorators(sdk.Terminator{}))
	require.NoError(s.T(), err)
}

func (s *AnteTestSuite) TestInvalidContractAddress() {
	s.Setup()
	s.Ctx = s.Ctx.WithValue(smartaccounttypes.ModuleName, smartaccounttypes.Setting{
		PreTransaction: []string{s.TestAccs[0].String()},
	})
	txBuilder := s.BuildDefaultMsgTx(0, &types.MsgSend{
		FromAddress: s.TestAccs[0].String(),
		ToAddress:   s.TestAccs[1].String(),
		Amount:      sdk.NewCoins(sdk.NewInt64Coin("uluna", 100000000)),
	})
	_, err := s.PreTxDecorator.AnteHandle(s.Ctx, txBuilder.GetTx(), false, sdk.ChainAnteDecorators(sdk.Terminator{}))
	require.ErrorContainsf(s.T(), err, "no such contract", "error message: %s", err)
}

func (s *AnteTestSuite) TestSendCoinsWithLimitSendHook() {
	s.Setup()

	acc := s.TestAccs[0]
	codeId, _, err := s.WasmKeeper.Create(s.Ctx, acc, test_helpers.LimitSendOnlyHookWasm, nil)
	require.NoError(s.T(), err)
	contractAddr, _, err := s.WasmKeeper.Instantiate(s.Ctx, codeId, acc, acc, []byte("{}"), "limit send", sdk.NewCoins())
	require.NoError(s.T(), err)

	s.Ctx = s.Ctx.WithValue(smartaccounttypes.ModuleName, smartaccounttypes.Setting{
		PreTransaction: []string{contractAddr.String()},
	})
	txBuilder := s.BuildDefaultMsgTx(0, &types.MsgSend{
		FromAddress: acc.String(),
		ToAddress:   acc.String(),
		Amount:      sdk.NewCoins(sdk.NewInt64Coin("uluna", 100000000)),
	})
	_, err = s.PreTxDecorator.AnteHandle(s.Ctx, txBuilder.GetTx(), false, sdk.ChainAnteDecorators(sdk.Terminator{}))
	require.NoError(s.T(), err)
}

func (s *AnteTestSuite) TestStakingWithLimitSendHook() {
	s.Setup()

	acc := s.TestAccs[0]
	codeId, _, err := s.WasmKeeper.Create(s.Ctx, acc, test_helpers.LimitSendOnlyHookWasm, nil)
	require.NoError(s.T(), err)
	contractAddr, _, err := s.WasmKeeper.Instantiate(s.Ctx, codeId, acc, acc, []byte("{}"), "limit send", sdk.NewCoins())
	require.NoError(s.T(), err)

	s.Ctx = s.Ctx.WithValue(smartaccounttypes.ModuleName, smartaccounttypes.Setting{
		PreTransaction: []string{contractAddr.String()},
	})
	txBuilder := s.BuildDefaultMsgTx(0, &stakingtypes.MsgDelegate{
		DelegatorAddress: acc.String(),
		ValidatorAddress: acc.String(),
		Amount:           sdk.NewInt64Coin("uluna", 100000000),
	})
	_, err = s.PreTxDecorator.AnteHandle(s.Ctx, txBuilder.GetTx(), false, sdk.ChainAnteDecorators(sdk.Terminator{}))
	require.ErrorContainsf(s.T(), err, "Unauthorized message type", "error message: %s", err)
}
