package tests

import (
	"encoding/base64"
	"encoding/json"
	"testing"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/terra-money/core/v2/x/smartaccount/ante"
	"github.com/terra-money/core/v2/x/smartaccount/test_helpers"
	smartaccounttypes "github.com/terra-money/core/v2/x/smartaccount/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
)

type AnteTestSuite struct {
	test_helpers.SmartAccountTestSuite

	AuthDecorator  ante.SmartAccountAuthDecorator
	PreTxDecorator ante.PreTransactionHookDecorator
	WasmKeeper     *wasmkeeper.PermissionedKeeper
}

func TestAnteSuite(t *testing.T) {
	suite.Run(t, new(AnteTestSuite))
}

func (s *AnteTestSuite) Setup() {
	s.SmartAccountTestSuite.SetupTests()
	s.WasmKeeper = wasmkeeper.NewDefaultPermissionKeeper(s.App.Keepers.WasmKeeper)
	s.AuthDecorator = ante.NewSmartAccountAuthDecorator(s.SmartAccountKeeper, s.WasmKeeper, s.App.Keepers.AccountKeeper, nil, s.EncodingConfig.TxConfig.SignModeHandler())
	s.PreTxDecorator = ante.NewPreTransactionHookDecorator(s.SmartAccountKeeper, s.WasmKeeper)
	s.Ctx = s.Ctx.WithChainID("test")
}

func (s *AnteTestSuite) TestAuthAnteHandler() {
	s.Setup()

	// testAcc1 using private key of testAcc0
	acc := s.TestAccs[1]
	pubKey := s.TestAccPrivs[0].PubKey()
	// endcoding this since this should be encoded in base64 when submitted by the user
	pkEncoded := []byte(base64.StdEncoding.EncodeToString(pubKey.Bytes()))

	codeId, _, err := s.WasmKeeper.Create(s.Ctx, acc, test_helpers.SmartAuthContractWasm, nil)
	require.NoError(s.T(), err)
	contractAddr, _, err := s.WasmKeeper.Instantiate(s.Ctx, codeId, acc, acc, []byte("{}"), "auth", sdk.NewCoins())
	require.NoError(s.T(), err)

	// create initMsg
	initMsg := smartaccounttypes.Initialization{
		Account: acc.String(),
		Msg:     pkEncoded,
	}
	sudoInitMsg := smartaccounttypes.SudoMsg{Initialization: &initMsg}
	sudoInitMsgBs, err := json.Marshal(sudoInitMsg)
	require.NoError(s.T(), err)

	_, err = s.WasmKeeper.Sudo(s.Ctx, contractAddr, sudoInitMsgBs)
	require.NoError(s.T(), err)

	// set settings
	authMsg := &smartaccounttypes.AuthorizationMsg{
		ContractAddress: contractAddr.String(),
		InitMsg:         sudoInitMsg.Initialization,
	}
	err = s.SmartAccountKeeper.SetSetting(s.Ctx, smartaccounttypes.Setting{
		Owner:         acc.String(),
		Authorization: []*smartaccounttypes.AuthorizationMsg{authMsg},
	})
	require.NoError(s.T(), err)

	// signing with testAcc1 pk which should error
	txBuilder := s.BuildDefaultMsgTx(1, &types.MsgSend{
		FromAddress: acc.String(),
		ToAddress:   acc.String(),
		Amount:      sdk.NewCoins(sdk.NewInt64Coin("uluna", 1)),
	})
	_, err = s.AuthDecorator.AnteHandle(s.Ctx, txBuilder.GetTx(), false, sdk.ChainAnteDecorators(sdk.Terminator{}))
	require.Error(s.T(), err)

	// signing with testAcc0 pk which should pass
	txBuilder = s.BuildDefaultMsgTx(0, &types.MsgSend{
		FromAddress: acc.String(),
		ToAddress:   acc.String(),
		Amount:      sdk.NewCoins(sdk.NewInt64Coin("uluna", 1)),
	})
	_, err = s.AuthDecorator.AnteHandle(s.Ctx, txBuilder.GetTx(), false, sdk.ChainAnteDecorators(sdk.Terminator{}))
	require.NoError(s.T(), err)
}

func (s *AnteTestSuite) TestAuthAnteHandlerWithMultisig() {
	s.Setup()

	type MultiSigSetting struct {
		Signers    []string `json:"signers"`
		PublicKeys []string `json:"public_keys"`
		Threshold  uint8    `json:"threshold"`
	}

	acc1 := s.TestAccs[1]
	acc2 := s.TestAccs[2]
	pubKey1 := s.TestAccPrivs[1].PubKey()
	pubKey2 := s.TestAccPrivs[2].PubKey()
	pkEncoded1 := []byte(base64.StdEncoding.EncodeToString(pubKey1.Bytes()))
	pkEncoded2 := []byte(base64.StdEncoding.EncodeToString(pubKey2.Bytes()))

	codeId, _, err := s.WasmKeeper.Create(s.Ctx, acc1, test_helpers.SmartMultiSigWasm, nil)
	require.NoError(s.T(), err)
	contractAddr, _, err := s.WasmKeeper.Instantiate(s.Ctx, codeId, acc1, acc1, []byte("{}"), "auth", sdk.NewCoins())
	require.NoError(s.T(), err)

	setting := MultiSigSetting{
		Signers:    []string{acc1.String(), acc2.String()},
		PublicKeys: []string{string(pkEncoded1), string(pkEncoded2)},
		Threshold:  2,
	}
	settingBs, err := json.Marshal(setting)
	require.NoError(s.T(), err)

	// create initMsg
	initMsg := smartaccounttypes.Initialization{
		Account: acc1.String(),
		Msg:     settingBs,
	}
	sudoInitMsg := smartaccounttypes.SudoMsg{Initialization: &initMsg}
	sudoInitMsgBs, err := json.Marshal(sudoInitMsg)
	require.NoError(s.T(), err)

	_, err = s.WasmKeeper.Sudo(s.Ctx, contractAddr, sudoInitMsgBs)
	require.NoError(s.T(), err)

	// set settings
	authMsg := &smartaccounttypes.AuthorizationMsg{
		ContractAddress: contractAddr.String(),
		InitMsg:         sudoInitMsg.Initialization,
	}
	err = s.SmartAccountKeeper.SetSetting(s.Ctx, smartaccounttypes.Setting{
		Owner:         acc1.String(),
		Fallback:      false,
		Authorization: []*smartaccounttypes.AuthorizationMsg{authMsg},
	})
	require.NoError(s.T(), err)

	// signing with multisig should pass
	txBuilder := s.BuildDefaultMultiSigMsgTx([]int{1, 2}, &types.MsgSend{
		FromAddress: acc1.String(),
		ToAddress:   acc1.String(),
		Amount:      sdk.NewCoins(sdk.NewInt64Coin("uluna", 1)),
	})
	_, err = s.AuthDecorator.AnteHandle(s.Ctx, txBuilder.GetTx(), false, sdk.ChainAnteDecorators(sdk.Terminator{}))
	require.NoError(s.T(), err)

	// signing with multisig with only 1 signer should fail
	txBuilder = s.BuildDefaultMultiSigMsgTx([]int{1}, &types.MsgSend{
		FromAddress: acc1.String(),
		ToAddress:   acc1.String(),
		Amount:      sdk.NewCoins(sdk.NewInt64Coin("uluna", 1)),
	})
	_, err = s.AuthDecorator.AnteHandle(s.Ctx, txBuilder.GetTx(), false, sdk.ChainAnteDecorators(sdk.Terminator{}))
	require.Error(s.T(), err)
	require.ErrorContainsf(s.T(), err, "required: 2, found: 1", "")
}

func (s *AnteTestSuite) BuildDefaultMsgTx(accountIndex int, msgs ...sdk.Msg) client.TxBuilder {
	pk := s.TestAccPrivs[accountIndex]
	sender := s.TestAccs[accountIndex]
	senderAcc := s.App.Keepers.AccountKeeper.GetAccount(s.Ctx, sender)
	senderSeq := senderAcc.GetSequence()
	txBuilder := s.EncodingConfig.TxConfig.NewTxBuilder()
	err := txBuilder.SetMsgs(
		msgs...,
	)
	require.NoError(s.T(), err)

	signer := authsigning.SignerData{
		Address:       sender.String(),
		ChainID:       "test",
		AccountNumber: senderAcc.GetAccountNumber(),
		Sequence:      senderSeq,
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
		senderSeq,
	)
	require.NoError(s.T(), err)

	err = txBuilder.SetSignatures(sigV2)
	require.NoError(s.T(), err)

	return txBuilder
}

func (s *AnteTestSuite) BuildDefaultMultiSigMsgTx(accountIndices []int, msgs ...sdk.Msg) client.TxBuilder {

	txBuilder := s.EncodingConfig.TxConfig.NewTxBuilder()
	err := txBuilder.SetMsgs(
		msgs...,
	)
	require.NoError(s.T(), err)

	var emptySigs []signing.SignatureV2

	for _, accountIndex := range accountIndices {
		sender := s.TestAccs[accountIndex]
		senderAcc := s.App.Keepers.AccountKeeper.GetAccount(s.Ctx, sender)
		senderSeq := senderAcc.GetSequence()
		pk := s.TestAccPrivs[accountIndex]

		signer := authsigning.SignerData{
			Address:       sender.String(),
			ChainID:       "test",
			AccountNumber: senderAcc.GetAccountNumber(),
			Sequence:      senderSeq,
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

		emptySigs = append(emptySigs, emptySig)
	}
	err = txBuilder.SetSignatures(emptySigs...)

	require.NoError(s.T(), err)

	var signed []signing.SignatureV2
	for _, accountIndex := range accountIndices {
		sender := s.TestAccs[accountIndex]
		senderAcc := s.App.Keepers.AccountKeeper.GetAccount(s.Ctx, sender)
		senderSeq := senderAcc.GetSequence()
		pk := s.TestAccPrivs[accountIndex]

		signer := authsigning.SignerData{
			Address:       sender.String(),
			ChainID:       "test",
			AccountNumber: senderAcc.GetAccountNumber(),
			Sequence:      senderSeq,
			PubKey:        pk.PubKey(),
		}
		sigV2, err := tx.SignWithPrivKey(
			s.EncodingConfig.TxConfig.SignModeHandler().DefaultMode(),
			signer,
			txBuilder,
			pk,
			s.EncodingConfig.TxConfig,
			senderSeq,
		)
		require.NoError(s.T(), err)
		signed = append(signed, sigV2)
	}

	err = txBuilder.SetSignatures(signed...)
	require.NoError(s.T(), err)

	return txBuilder
}
