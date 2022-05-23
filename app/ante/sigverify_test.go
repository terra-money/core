package ante_test

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	cosmosante "github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/auth/legacy/legacytx"

	"github.com/terra-money/core/v2/app/ante"
)

func (suite *AnteTestSuite) TestSigVerification() {
	suite.SetupTest(true) // setup
	suite.txBuilder = suite.clientCtx.TxConfig.NewTxBuilder()

	// make block height non-zero to ensure account numbers part of signBytes
	suite.ctx = suite.ctx.WithBlockHeight(1)

	// keys and addresses
	priv1, _, addr1 := testdata.KeyTestPubAddr()
	priv2, _, addr2 := testdata.KeyTestPubAddr()
	priv3, _, addr3 := testdata.KeyTestPubAddr()

	addrs := []sdk.AccAddress{addr1, addr2, addr3}

	msgs := make([]sdk.Msg, len(addrs))
	// set accounts and create msg for each address
	for i, addr := range addrs {
		acc := suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, addr)
		suite.Require().NoError(acc.SetAccountNumber(uint64(i)))
		suite.app.AccountKeeper.SetAccount(suite.ctx, acc)
		msgs[i] = testdata.NewTestMsg(addr)
	}

	feeAmount := testdata.NewTestFeeAmount()
	gasLimit := testdata.NewTestGasLimit()

	spkd := cosmosante.NewSetPubKeyDecorator(suite.app.AccountKeeper)
	svd := ante.NewSigVerificationDecorator(suite.app.AccountKeeper, suite.clientCtx.TxConfig.SignModeHandler())
	antehandler := sdk.ChainAnteDecorators(spkd, svd)

	type testCase struct {
		name        string
		privs       []cryptotypes.PrivKey
		accNums     []uint64
		accSeqs     []uint64
		invalidSigs bool
		recheck     bool
		shouldErr   bool
	}
	validSigs := false
	testCases := []testCase{
		{"no signers", []cryptotypes.PrivKey{}, []uint64{}, []uint64{}, validSigs, false, true},
		{"not enough signers", []cryptotypes.PrivKey{priv1, priv2}, []uint64{0, 1}, []uint64{0, 0}, validSigs, false, true},
		{"wrong order signers", []cryptotypes.PrivKey{priv3, priv2, priv1}, []uint64{2, 1, 0}, []uint64{0, 0, 0}, validSigs, false, true},
		{"wrong accnums", []cryptotypes.PrivKey{priv1, priv2, priv3}, []uint64{7, 8, 9}, []uint64{0, 0, 0}, validSigs, false, true},
		{"wrong sequences", []cryptotypes.PrivKey{priv1, priv2, priv3}, []uint64{0, 1, 2}, []uint64{3, 4, 5}, validSigs, false, true},
		{"valid tx", []cryptotypes.PrivKey{priv1, priv2, priv3}, []uint64{0, 1, 2}, []uint64{0, 0, 0}, validSigs, false, false},
		{"no err on recheck", []cryptotypes.PrivKey{priv1, priv2, priv3}, []uint64{0, 0, 0}, []uint64{0, 0, 0}, !validSigs, true, false},
	}
	for i, tc := range testCases {
		suite.ctx = suite.ctx.WithIsReCheckTx(tc.recheck)
		suite.txBuilder = suite.clientCtx.TxConfig.NewTxBuilder() // Create new txBuilder for each test

		suite.Require().NoError(suite.txBuilder.SetMsgs(msgs...))
		suite.txBuilder.SetFeeAmount(feeAmount)
		suite.txBuilder.SetGasLimit(gasLimit)

		tx, err := suite.CreateTestTx(tc.privs, tc.accNums, tc.accSeqs, suite.ctx.ChainID())
		suite.Require().NoError(err)
		if tc.invalidSigs {
			txSigs, _ := tx.GetSignaturesV2()
			badSig, _ := tc.privs[0].Sign([]byte("unrelated message"))
			txSigs[0] = signing.SignatureV2{
				PubKey: tc.privs[0].PubKey(),
				Data: &signing.SingleSignatureData{
					SignMode:  suite.clientCtx.TxConfig.SignModeHandler().DefaultMode(),
					Signature: badSig,
				},
				Sequence: tc.accSeqs[0],
			}
			suite.txBuilder.SetSignatures(txSigs...)
			tx = suite.txBuilder.GetTx()
		}

		_, err = antehandler(suite.ctx, tx, false)
		if tc.shouldErr {
			suite.Require().NotNil(err, "TestCase %d: %s did not error as expected", i, tc.name)
		} else {
			suite.Require().Nil(err, "TestCase %d: %s errored unexpectedly. Err: %v", i, tc.name, err)
		}
	}
}

// This test is exactly like the one above, but we set the codec explicitly to
// Amino.
// Once https://github.com/cosmos/cosmos-sdk/issues/6190 is in, we can remove
// this, since it'll be handled by the test matrix.
// In the meantime, we want to make double-sure amino compatibility works.
// ref: https://github.com/cosmos/cosmos-sdk/issues/7229
func (suite *AnteTestSuite) TestSigVerification_ExplicitAmino() {
	tempDir := suite.T().TempDir()
	suite.app, suite.ctx = createTestApp(true, tempDir)
	suite.ctx = suite.ctx.WithBlockHeight(1)

	// Set up TxConfig.
	aminoCdc := codec.NewLegacyAmino()
	// We're using TestMsg amino encoding in some tests, so register it here.
	txConfig := legacytx.StdTxConfig{Cdc: aminoCdc}

	suite.clientCtx = client.Context{}.
		WithTxConfig(txConfig)

	anteHandler, err := ante.NewAnteHandler(
		ante.HandlerOptions{
			HandlerOptions: cosmosante.HandlerOptions{
				AccountKeeper:   suite.app.AccountKeeper,
				BankKeeper:      suite.app.BankKeeper,
				FeegrantKeeper:  suite.app.FeeGrantKeeper,
				SignModeHandler: txConfig.SignModeHandler(),
				SigGasConsumer:  cosmosante.DefaultSigVerificationGasConsumer,
			},
		},
	)

	suite.Require().NoError(err)
	suite.anteHandler = anteHandler

	suite.txBuilder = suite.clientCtx.TxConfig.NewTxBuilder()

	// make block height non-zero to ensure account numbers part of signBytes
	suite.ctx = suite.ctx.WithBlockHeight(1)

	// keys and addresses
	priv1, _, addr1 := testdata.KeyTestPubAddr()
	priv2, _, addr2 := testdata.KeyTestPubAddr()
	priv3, _, addr3 := testdata.KeyTestPubAddr()

	addrs := []sdk.AccAddress{addr1, addr2, addr3}

	msgs := make([]sdk.Msg, len(addrs))
	// set accounts and create msg for each address
	for i, addr := range addrs {
		acc := suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, addr)
		suite.Require().NoError(acc.SetAccountNumber(uint64(i)))
		suite.app.AccountKeeper.SetAccount(suite.ctx, acc)
		msgs[i] = testdata.NewTestMsg(addr)
	}

	feeAmount := testdata.NewTestFeeAmount()
	gasLimit := testdata.NewTestGasLimit()

	spkd := cosmosante.NewSetPubKeyDecorator(suite.app.AccountKeeper)
	svd := ante.NewSigVerificationDecorator(suite.app.AccountKeeper, suite.clientCtx.TxConfig.SignModeHandler())
	antehandler := sdk.ChainAnteDecorators(spkd, svd)

	type testCase struct {
		name      string
		privs     []cryptotypes.PrivKey
		accNums   []uint64
		accSeqs   []uint64
		recheck   bool
		shouldErr bool
	}
	testCases := []testCase{
		{"no signers", []cryptotypes.PrivKey{}, []uint64{}, []uint64{}, false, true},
		{"not enough signers", []cryptotypes.PrivKey{priv1, priv2}, []uint64{0, 1}, []uint64{0, 0}, false, true},
		{"wrong order signers", []cryptotypes.PrivKey{priv3, priv2, priv1}, []uint64{2, 1, 0}, []uint64{0, 0, 0}, false, true},
		{"wrong accnums", []cryptotypes.PrivKey{priv1, priv2, priv3}, []uint64{7, 8, 9}, []uint64{0, 0, 0}, false, true},
		{"wrong sequences", []cryptotypes.PrivKey{priv1, priv2, priv3}, []uint64{0, 1, 2}, []uint64{3, 4, 5}, false, true},
		{"valid tx", []cryptotypes.PrivKey{priv1, priv2, priv3}, []uint64{0, 1, 2}, []uint64{0, 0, 0}, false, false},
		{"no err on recheck", []cryptotypes.PrivKey{priv1, priv2, priv3}, []uint64{0, 1, 2}, []uint64{0, 0, 0}, true, false},
	}
	for i, tc := range testCases {
		suite.ctx = suite.ctx.WithIsReCheckTx(tc.recheck)
		suite.txBuilder = suite.clientCtx.TxConfig.NewTxBuilder() // Create new txBuilder for each test

		suite.Require().NoError(suite.txBuilder.SetMsgs(msgs...))
		suite.txBuilder.SetFeeAmount(feeAmount)
		suite.txBuilder.SetGasLimit(gasLimit)

		tx, err := suite.CreateTestTx(tc.privs, tc.accNums, tc.accSeqs, suite.ctx.ChainID())
		suite.Require().NoError(err)

		_, err = antehandler(suite.ctx, tx, false)
		if tc.shouldErr {
			suite.Require().NotNil(err, "TestCase %d: %s did not error as expected", i, tc.name)
		} else {
			suite.Require().Nil(err, "TestCase %d: %s errored unexpectedly. Err: %v", i, tc.name, err)
		}
	}
}
