package ante_test

import (
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/terra-money/core/v2/app/ante"
	"github.com/terra-money/core/v2/app/ante/types"
)

func (suite *AnteTestSuite) TestMinCommission() {
	suite.SetupTest(true) // setup
	suite.txBuilder = suite.clientCtx.TxConfig.NewTxBuilder()

	// make block height non-zero to ensure account numbers part of signBytes
	suite.ctx = suite.ctx.WithBlockHeight(1)

	// keys and addresses
	_, pub1, addr1 := testdata.KeyTestPubAddr()
	_, _, addr2 := testdata.KeyTestPubAddr()

	min := ante.NewMinCommissionDecorator(suite.app.AppCodec(), &suite.app.AnteKeeper)
	lowCommission := types.DefaultMinimumCommission.QuoInt64(2)
	highCommission := types.DefaultMinimumCommission

	// create validator
	createValidator, err := stakingtypes.NewMsgCreateValidator(
		sdk.ValAddress(addr1),
		pub1,
		sdk.NewInt64Coin("foo", 1),
		stakingtypes.Description{},
		stakingtypes.NewCommissionRates(
			lowCommission,
			sdk.NewDecWithPrec(100, 2), // 100%
			sdk.NewDecWithPrec(1, 2),   // 1%
		),
		sdk.NewInt(1),
	)
	suite.NoError(err)

	antehandler := sdk.ChainAnteDecorators(min)

	// with low commission
	suite.txBuilder.SetMsgs(createValidator)
	_, err = antehandler(suite.ctx, suite.txBuilder.GetTx(), false)
	suite.Error(err)

	execMsg := authz.NewMsgExec(addr1, []sdk.Msg{createValidator})
	suite.txBuilder.SetMsgs(&execMsg)
	_, err = antehandler(suite.ctx, suite.txBuilder.GetTx(), false)
	suite.Error(err)

	// with high commission
	createValidator.Commission.Rate = highCommission
	suite.txBuilder.SetMsgs(createValidator)
	_, err = antehandler(suite.ctx, suite.txBuilder.GetTx(), false)
	suite.NoError(err)

	execMsg = authz.NewMsgExec(addr1, []sdk.Msg{createValidator})
	suite.txBuilder.SetMsgs(&execMsg)
	_, err = antehandler(suite.ctx, suite.txBuilder.GetTx(), false)
	suite.NoError(err)

	// edit validator
	editValidator := stakingtypes.NewMsgEditValidator(
		sdk.ValAddress(addr2),
		stakingtypes.Description{},
		&lowCommission,
		nil,
	)

	// with low commission
	suite.txBuilder.SetMsgs(editValidator)
	_, err = antehandler(suite.ctx, suite.txBuilder.GetTx(), false)
	suite.Error(err)

	execMsg = authz.NewMsgExec(addr1, []sdk.Msg{editValidator})
	suite.txBuilder.SetMsgs(&execMsg)
	_, err = antehandler(suite.ctx, suite.txBuilder.GetTx(), false)
	suite.Error(err)

	// with high commission
	editValidator.CommissionRate = &highCommission
	suite.txBuilder.SetMsgs(editValidator)
	_, err = antehandler(suite.ctx, suite.txBuilder.GetTx(), false)
	suite.NoError(err)

	execMsg = authz.NewMsgExec(addr1, []sdk.Msg{editValidator})
	suite.txBuilder.SetMsgs(&execMsg)
	_, err = antehandler(suite.ctx, suite.txBuilder.GetTx(), false)
	suite.NoError(err)
}
