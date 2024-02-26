package tests

import (
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/terra-money/core/v2/x/smartaccount/ante"
	"github.com/terra-money/core/v2/x/smartaccount/test_helpers"
	"testing"
)

type AnteTestSuite struct {
	test_helpers.SmartAccountTestSuite

	Decorator  ante.PreTransactionHookDecorator
	WasmKeeper *wasmkeeper.PermissionedKeeper
}

func TestAnteSuite(t *testing.T) {
	suite.Run(t, new(AnteTestSuite))
}

func (s *AnteTestSuite) Setup() {
	s.SmartAccountTestSuite.Setup()
	s.WasmKeeper = wasmkeeper.NewDefaultPermissionKeeper(s.App.Keepers.WasmKeeper)
	s.Decorator = ante.NewPreTransactionHookDecorator(s.SmartAccountKeeper, s.WasmKeeper)
}

func (s *AnteTestSuite) BuildDefaultMsgTx(accountIndex int, msgs ...sdk.Msg) client.TxBuilder {
	pk := s.TestAccPrivs[accountIndex]
	txBuilder := s.EncodingConfig.TxConfig.NewTxBuilder()
	err := txBuilder.SetMsgs(
		msgs...,
	)
	require.NoError(s.T(), err)
	signer := authsigning.SignerData{
		ChainID:       "test",
		AccountNumber: 0,
		Sequence:      0,
	}
	sig, err := tx.SignWithPrivKey(
		s.EncodingConfig.TxConfig.SignModeHandler().DefaultMode(),
		signer,
		txBuilder,
		pk,
		s.EncodingConfig.TxConfig,
		0,
	)
	txBuilder.SetSignatures(sig)
	return txBuilder
}
