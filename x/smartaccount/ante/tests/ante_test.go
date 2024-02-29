package tests

import (
	"encoding/base64"
	"encoding/json"
	"testing"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/terra-money/core/v2/x/smartaccount/ante"
	"github.com/terra-money/core/v2/x/smartaccount/test_helpers"
	smartaccounttypes "github.com/terra-money/core/v2/x/smartaccount/types"
)

type AnteTestSuite struct {
	test_helpers.SmartAccountTestSuite

	Decorator  ante.SmartAccountAuthDecorator
	WasmKeeper *wasmkeeper.PermissionedKeeper
}

func TestAnteSuite(t *testing.T) {
	suite.Run(t, new(AnteTestSuite))
}

func (s *AnteTestSuite) Setup() {
	s.SmartAccountTestSuite.SetupTests()
	s.WasmKeeper = wasmkeeper.NewDefaultPermissionKeeper(s.App.Keepers.WasmKeeper)
	s.Decorator = ante.NewSmartAccountAuthDecorator(s.SmartAccountKeeper, s.WasmKeeper, s.App.Keepers.AccountKeeper, nil, s.EncodingConfig.TxConfig.SignModeHandler())
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
		Sender:  acc.String(),
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
		InitMsg:         string(sudoInitMsgBs),
	}
	err = s.SmartAccountKeeper.SetSetting(s.Ctx, smartaccounttypes.Setting{
		Owner:         acc.String(),
		Authorization: []*smartaccounttypes.AuthorizationMsg{authMsg},
	})
	require.NoError(s.T(), err)

	txBuilder := s.BuildDefaultMsgTx(0, &types.MsgSend{
		FromAddress: acc.String(),
		ToAddress:   acc.String(),
		Amount:      sdk.NewCoins(sdk.NewInt64Coin("uluna", 1)),
	})
	_, err = s.Decorator.AnteHandle(s.Ctx, txBuilder.GetTx(), false, sdk.ChainAnteDecorators(sdk.Terminator{}))
	require.NoError(s.T(), err)
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
	require.NoError(s.T(), err)
	txBuilder.SetSignatures(sig)
	return txBuilder
}
