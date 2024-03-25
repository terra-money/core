package ante_test

import (
	"testing"

	errorsmod "cosmossdk.io/errors"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/terra-money/core/v2/app/post/mocks"
	"github.com/terra-money/core/v2/app/test_helpers"
	post "github.com/terra-money/core/v2/x/feeshare/post"
	"github.com/terra-money/core/v2/x/feeshare/types"
	customwasmtypes "github.com/terra-money/core/v2/x/wasm/types"
)

type AnteTestSuite struct {
	test_helpers.AppTestSuite
}

func TestAnteSuite(t *testing.T) {
	suite.Run(t, new(AnteTestSuite))
}

func (suite *AnteTestSuite) TestGetWithdrawalAddressFromContract() {
	suite.Setup()

	feeshareKeeper := suite.AppTestSuite.App.Keepers.FeeShareKeeper
	feeshareKeeper.SetFeeShare(suite.Ctx, types.FeeShare{
		ContractAddress:   "terra1jwyzzsaag4t0evnuukc35ysyrx9arzdde2kg9cld28alhjurtthq0prs2s",
		DeployerAddress:   "",
		WithdrawerAddress: "terra1zdpgj8am5nqqvht927k3etljyl6a52kwqup0je",
	})
	feeshareKeeper.SetFeeShare(suite.Ctx, types.FeeShare{
		ContractAddress:   "terra1mdpvgjc8jmv60a4x68nggsh9w8uyv69sqls04a76m9med5hsqmwsse8sxa",
		DeployerAddress:   "",
		WithdrawerAddress: "",
	})

	testCases := []struct {
		name                    string
		contractAddresses       []string
		expectedWithdrawerAddrs []sdk.AccAddress
		expectErr               bool
	}{
		{
			"valid contract addresses",
			[]string{"terra1jwyzzsaag4t0evnuukc35ysyrx9arzdde2kg9cld28alhjurtthq0prs2s"},
			[]sdk.AccAddress{
				sdk.MustAccAddressFromBech32("terra1zdpgj8am5nqqvht927k3etljyl6a52kwqup0je"),
			},
			false,
		},
		{
			"two valid contract addresses with one not registered",
			[]string{
				"terra1u3z42fpctuhh8mranz4tatacqhty6a8yk7l5wvj7dshsuytcms2qda4f5x", // not registered address
				"terra1jwyzzsaag4t0evnuukc35ysyrx9arzdde2kg9cld28alhjurtthq0prs2s",
			},
			[]sdk.AccAddress{
				sdk.MustAccAddressFromBech32("terra1zdpgj8am5nqqvht927k3etljyl6a52kwqup0je"),
			},
			false,
		},
		{
			"without withdrawer contract addresses",
			[]string{"terra1mdpvgjc8jmv60a4x68nggsh9w8uyv69sqls04a76m9med5hsqmwsse8sxa"},
			[]sdk.AccAddress(nil),
			false,
		},
		{
			"invalid contract address",
			[]string{"invalidAddress"},
			nil,
			true,
		},
	}

	for _, tc := range testCases {
		withdrawerAddrs, err := post.GetWithdrawalAddressFromContract(
			suite.Ctx,
			tc.contractAddresses,
			feeshareKeeper,
		)

		if tc.expectErr {
			suite.Require().Error(err, tc.name)
		} else {
			suite.Require().NoError(err, tc.name)
			suite.Require().Equal(tc.expectedWithdrawerAddrs, withdrawerAddrs, tc.name)
		}
	}
}

func (suite *AnteTestSuite) TestCalculateFee() {
	feeCoins := sdk.NewCoins(sdk.NewCoin("uluna", sdk.NewInt(500)), sdk.NewCoin("utoken", sdk.NewInt(250)))

	testCases := []struct {
		name               string
		incomingFee        sdk.Coins
		devShares          sdk.Dec
		numOfdevs          int
		allowdDenoms       []string
		expectedFeePayment sdk.Coins
	}{
		{
			"100% fee / 1 contract",
			feeCoins,
			sdk.NewDecWithPrec(100, 2),
			1,
			[]string{},
			sdk.NewCoins(sdk.NewCoin("uluna", sdk.NewInt(500)), sdk.NewCoin("utoken", sdk.NewInt(250))),
		},
		{
			"100% fee / 2 contracts",
			feeCoins,
			sdk.NewDecWithPrec(100, 2),
			2,
			[]string{},
			sdk.NewCoins(sdk.NewCoin("uluna", sdk.NewInt(250)), sdk.NewCoin("utoken", sdk.NewInt(125))),
		},
		{
			"100% fee / 10 contracts / 1 allowed denom",
			feeCoins,
			sdk.NewDecWithPrec(100, 2),
			10,
			[]string{"uluna"},
			sdk.NewCoins(sdk.NewCoin("uluna", sdk.NewInt(50))),
		},
		{
			"67% fee / 7 contracts",
			feeCoins,
			sdk.NewDecWithPrec(67, 2),
			7,
			[]string{},
			sdk.NewCoins(sdk.NewCoin("uluna", sdk.NewInt(48)), sdk.NewCoin("utoken", sdk.NewInt(24))),
		},
		{
			"50% fee / 1 contracts / 1 allowed denom",
			feeCoins,
			sdk.NewDecWithPrec(50, 2),
			1,
			[]string{"utoken"},
			sdk.NewCoins(sdk.NewCoin("utoken", sdk.NewInt(125))),
		},
		{
			"50% fee / 2 contracts / 2 allowed denoms",
			feeCoins,
			sdk.NewDecWithPrec(50, 2),
			2,
			[]string{"uluna", "utoken"},
			sdk.NewCoins(sdk.NewCoin("uluna", sdk.NewInt(125)), sdk.NewCoin("utoken", sdk.NewInt(62))),
		},
		{
			"50% fee / 3 contracts",
			feeCoins,
			sdk.NewDecWithPrec(50, 2),
			3,
			[]string{},
			sdk.NewCoins(sdk.NewCoin("uluna", sdk.NewInt(83)), sdk.NewCoin("utoken", sdk.NewInt(42))),
		},
		{
			"25% fee / 2 contracts",
			feeCoins,
			sdk.NewDecWithPrec(25, 2),
			2,
			[]string{},
			sdk.NewCoins(sdk.NewCoin("uluna", sdk.NewInt(62)), sdk.NewCoin("utoken", sdk.NewInt(31))),
		},
		{
			"15% fee / 3 contracts / inexistent denom",
			feeCoins,
			sdk.NewDecWithPrec(15, 2),
			3,
			[]string{"ubtc"},
			sdk.Coins(nil),
		},
		{
			"1% fee / 2 contracts",
			feeCoins,
			sdk.NewDecWithPrec(1, 2),
			2,
			[]string{},
			sdk.NewCoins(sdk.NewCoin("uluna", sdk.NewInt(2)), sdk.NewCoin("utoken", sdk.NewInt(1))),
		},
	}

	for _, tc := range testCases {
		feeToBePaid := post.CalculateFee(tc.incomingFee, tc.devShares, tc.numOfdevs, tc.allowdDenoms)

		suite.Require().Equal(tc.expectedFeePayment, feeToBePaid, tc.name)
	}
}

func (suite *AnteTestSuite) TestPostHandler() {
	suite.Setup()

	// Create a mocked next post handler to assert the function being called.
	ctrl := gomock.NewController(suite.T())
	mockedPostDecorator := mocks.NewMockPostDecorator(ctrl)

	// Register the feeshare contract...
	suite.App.Keepers.FeeShareKeeper.SetFeeShare(suite.Ctx, types.FeeShare{
		ContractAddress:   "terra1mdpvgjc8jmv60a4x68nggsh9w8uyv69sqls04a76m9med5hsqmwsse8sxa",
		DeployerAddress:   "",
		WithdrawerAddress: "terra1zdpgj8am5nqqvht927k3etljyl6a52kwqup0je",
	})
	// ... append the executed contract addresses in the wasm keeper ...
	suite.App.Keepers.WasmKeeper.SetExecutedContractAddresses(suite.Ctx, customwasmtypes.ExecutedContracts{
		ContractAddresses: []string{"terra1mdpvgjc8jmv60a4x68nggsh9w8uyv69sqls04a76m9med5hsqmwsse8sxa"},
	})

	// build a tx with a fee amount ...
	txFee := sdk.NewCoins(sdk.NewCoin("uluna", sdk.NewInt(500)), sdk.NewCoin("utoken", sdk.NewInt(250)))
	txBuilder := suite.EncodingConfig.TxConfig.NewTxBuilder()
	txBuilder.SetFeeAmount(txFee)
	txBuilder.SetMsgs(&wasmtypes.MsgExecuteContract{
		Sender:   "terra1zdpgj8am5nqqvht927k3etljyl6a52kwqup0je",
		Contract: "terra1mdpvgjc8jmv60a4x68nggsh9w8uyv69sqls04a76m9med5hsqmwsse8sxa",
		Msg:      nil,
		Funds:    nil,
	})
	// ... create the feeshare post handler ...
	handler := post.NewFeeSharePayoutDecorator(
		suite.App.Keepers.FeeShareKeeper,
		suite.App.Keepers.BankKeeper,
		suite.App.Keepers.WasmKeeper,
	)
	// Remove all events from the context to assert the events being added correctly.
	suite.Ctx = suite.Ctx.WithEventManager(sdk.NewEventManager())

	// Assert the next handler is called once
	mockedPostDecorator.
		EXPECT().
		PostHandle(gomock.Any(), gomock.Any(), false, true, gomock.Any()).
		Times(1)

	// Execute the PostHandle function
	_, err := handler.PostHandle(
		suite.Ctx,
		txBuilder.GetTx(),
		false,
		true,
		func(ctx sdk.Context, tx sdk.Tx, simulate bool, success bool) (sdk.Context, error) {
			return mockedPostDecorator.PostHandle(ctx, tx, simulate, success, nil)
		},
	)
	suite.Require().NoError(err)
	suite.Require().Equal(suite.Ctx.EventManager().ABCIEvents(),
		[]abci.Event{
			{
				Type: "coin_spent",
				Attributes: []abci.EventAttribute{
					{Key: "spender", Value: "terra17xpfvakm2amg962yls6f84z3kell8c5lkaeqfa", Index: false},
					{Key: "amount", Value: "250uluna,125utoken", Index: false},
				},
			},
			{
				Type: "coin_received",
				Attributes: []abci.EventAttribute{
					{Key: "receiver", Value: "terra1zdpgj8am5nqqvht927k3etljyl6a52kwqup0je", Index: false},
					{Key: "amount", Value: "250uluna,125utoken", Index: false},
				},
			},
			{
				Type: "transfer",
				Attributes: []abci.EventAttribute{
					{Key: "recipient", Value: "terra1zdpgj8am5nqqvht927k3etljyl6a52kwqup0je", Index: false},
					{Key: "sender", Value: "terra17xpfvakm2amg962yls6f84z3kell8c5lkaeqfa", Index: false},
					{Key: "amount", Value: "250uluna,125utoken", Index: false},
				},
			},
			{
				Type: "message",
				Attributes: []abci.EventAttribute{
					{Key: "sender", Value: "terra17xpfvakm2amg962yls6f84z3kell8c5lkaeqfa", Index: false},
				},
			},
			{
				Type: "juno.feeshare.v1.FeePayoutEvent",
				Attributes: []abci.EventAttribute{
					{Key: "fees_paid", Value: "[{\"denom\":\"uluna\",\"amount\":\"250\"},{\"denom\":\"utoken\",\"amount\":\"125\"}]", Index: false},
					{Key: "withdraw_address", Value: "\"terra1zdpgj8am5nqqvht927k3etljyl6a52kwqup0je\"", Index: false},
				},
			},
		})
}

func (suite *AnteTestSuite) TestDisabledPostHandle() {
	suite.Setup()

	// Create a mocked next post handler to assert the function being called.
	ctrl := gomock.NewController(suite.T())
	mockedPostDecorator := mocks.NewMockPostDecorator(ctrl)

	// Disable the feeshare module...
	err := suite.App.Keepers.FeeShareKeeper.SetParams(suite.Ctx, types.Params{
		EnableFeeShare:  false,
		DeveloperShares: sdk.MustNewDecFromStr("0.5"),
		AllowedDenoms:   []string{},
	})
	suite.Require().NoError(err)

	// build a tx with a fee amount ...
	txFee := sdk.NewCoins(sdk.NewCoin("uluna", sdk.NewInt(500)), sdk.NewCoin("utoken", sdk.NewInt(250)))
	txBuilder := suite.EncodingConfig.TxConfig.NewTxBuilder()
	txBuilder.SetFeeAmount(txFee)
	txBuilder.SetMsgs(&wasmtypes.MsgExecuteContract{})

	// ... create the feeshare post handler ...
	handler := post.NewFeeSharePayoutDecorator(
		suite.App.Keepers.FeeShareKeeper,
		suite.App.Keepers.BankKeeper,
		suite.App.Keepers.WasmKeeper,
	)

	// Assert the next handler is called once
	mockedPostDecorator.
		EXPECT().
		PostHandle(gomock.Any(), gomock.Any(), false, true, gomock.Any()).
		Times(1)

	// Execute the PostHandle function
	_, err = handler.PostHandle(
		suite.Ctx,
		txBuilder.GetTx(),
		false,
		true,
		func(ctx sdk.Context, tx sdk.Tx, simulate bool, success bool) (sdk.Context, error) {
			return mockedPostDecorator.PostHandle(ctx, tx, simulate, success, nil)
		},
	)
	suite.Require().NoError(err)
}

func (suite *AnteTestSuite) TestWithZeroFeesPostHandle() {
	suite.Setup()

	// Create a mocked next post handler to assert the function being called.
	ctrl := gomock.NewController(suite.T())
	mockedPostDecorator := mocks.NewMockPostDecorator(ctrl)

	// Build a tx with a fee amount ...
	txBuilder := suite.EncodingConfig.TxConfig.NewTxBuilder()

	// ... create the feeshare post handler ...
	handler := post.NewFeeSharePayoutDecorator(
		suite.App.Keepers.FeeShareKeeper,
		suite.App.Keepers.BankKeeper,
		suite.App.Keepers.WasmKeeper,
	)

	// Assert the next handler is called once
	mockedPostDecorator.
		EXPECT().
		PostHandle(gomock.Any(), gomock.Any(), false, true, gomock.Any()).
		Times(1)

	// Execute the PostHandle function
	_, err := handler.PostHandle(
		suite.Ctx,
		txBuilder.GetTx(),
		false,
		true,
		func(ctx sdk.Context, tx sdk.Tx, simulate bool, success bool) (sdk.Context, error) {
			return mockedPostDecorator.PostHandle(ctx, tx, simulate, success, nil)
		},
	)
	suite.Require().NoError(err)
}

func (suite *AnteTestSuite) TestPostHandlerWithEmptySmartContractStore() {
	suite.Setup()

	// Create a mocked next post handler to assert the function being called.
	ctrl := gomock.NewController(suite.T())
	mockedPostDecorator := mocks.NewMockPostDecorator(ctrl)

	// Register the feeshare contract...
	suite.App.Keepers.FeeShareKeeper.SetFeeShare(suite.Ctx, types.FeeShare{
		ContractAddress:   "terra1mdpvgjc8jmv60a4x68nggsh9w8uyv69sqls04a76m9med5hsqmwsse8sxa",
		DeployerAddress:   "",
		WithdrawerAddress: "terra1zdpgj8am5nqqvht927k3etljyl6a52kwqup0je",
	})

	// build a tx with a fee amount ...
	txFee := sdk.NewCoins(sdk.NewCoin("uluna", sdk.NewInt(500)), sdk.NewCoin("utoken", sdk.NewInt(250)))
	txBuilder := suite.EncodingConfig.TxConfig.NewTxBuilder()
	txBuilder.SetFeeAmount(txFee)
	txBuilder.SetMsgs(&wasmtypes.MsgExecuteContract{
		Sender:   "terra1zdpgj8am5nqqvht927k3etljyl6a52kwqup0je",
		Contract: "terra1mdpvgjc8jmv60a4x68nggsh9w8uyv69sqls04a76m9med5hsqmwsse8sxa",
		Msg:      nil,
		Funds:    nil,
	})
	// ... create the feeshare post handler ...
	handler := post.NewFeeSharePayoutDecorator(
		suite.App.Keepers.FeeShareKeeper,
		suite.App.Keepers.BankKeeper,
		suite.App.Keepers.WasmKeeper,
	)

	// Assert the next handler is called once
	mockedPostDecorator.
		EXPECT().
		PostHandle(gomock.Any(), gomock.Any(), false, true, gomock.Any()).
		Times(1)

	// Execute the PostHandle function
	_, err := handler.PostHandle(
		suite.Ctx,
		txBuilder.GetTx(),
		false,
		true,
		func(ctx sdk.Context, tx sdk.Tx, simulate bool, success bool) (sdk.Context, error) {
			return mockedPostDecorator.PostHandle(ctx, tx, simulate, success, nil)
		},
	)
	suite.Require().NoError(err)
}

func (suite *AnteTestSuite) TestPostHandlerNoSmartContractExecuted() {
	suite.Setup()

	// Create a mocked next post handler to assert the function being called.
	ctrl := gomock.NewController(suite.T())
	mockedPostDecorator := mocks.NewMockPostDecorator(ctrl)

	// Register the feeshare contract...
	suite.App.Keepers.FeeShareKeeper.SetFeeShare(suite.Ctx, types.FeeShare{
		ContractAddress:   "terra1mdpvgjc8jmv60a4x68nggsh9w8uyv69sqls04a76m9med5hsqmwsse8sxa",
		DeployerAddress:   "",
		WithdrawerAddress: "terra1zdpgj8am5nqqvht927k3etljyl6a52kwqup0je",
	})
	// ... create the store key ...
	suite.App.Keepers.WasmKeeper.SetExecutedContractAddresses(suite.Ctx, customwasmtypes.ExecutedContracts{
		ContractAddresses: []string{},
	})

	// build a tx with a fee amount ...
	txFee := sdk.NewCoins(sdk.NewCoin("uluna", sdk.NewInt(500)), sdk.NewCoin("utoken", sdk.NewInt(250)))
	txBuilder := suite.EncodingConfig.TxConfig.NewTxBuilder()
	txBuilder.SetFeeAmount(txFee)
	txBuilder.SetMsgs(&wasmtypes.MsgExecuteContract{
		Sender:   "terra1zdpgj8am5nqqvht927k3etljyl6a52kwqup0je",
		Contract: "terra1mdpvgjc8jmv60a4x68nggsh9w8uyv69sqls04a76m9med5hsqmwsse8sxa",
		Msg:      nil,
		Funds:    nil,
	})
	// ... create the feeshare post handler ...
	handler := post.NewFeeSharePayoutDecorator(
		suite.App.Keepers.FeeShareKeeper,
		suite.App.Keepers.BankKeeper,
		suite.App.Keepers.WasmKeeper,
	)

	// Assert the next handler is called once
	mockedPostDecorator.
		EXPECT().
		PostHandle(gomock.Any(), gomock.Any(), false, true, gomock.Any()).
		Times(1)

	// Execute the PostHandle function
	_, err := handler.PostHandle(
		suite.Ctx,
		txBuilder.GetTx(),
		false,
		true,
		func(ctx sdk.Context, tx sdk.Tx, simulate bool, success bool) (sdk.Context, error) {
			return mockedPostDecorator.PostHandle(ctx, tx, simulate, success, nil)
		},
	)
	suite.Require().NoError(err)
}

func (suite *AnteTestSuite) TestPostHandlerWithInvalidContractAddrOnExecution() {
	suite.Setup()

	// Create a mocked next post handler to assert the function being called.
	ctrl := gomock.NewController(suite.T())
	mockedPostDecorator := mocks.NewMockPostDecorator(ctrl)

	// Register the feeshare contract...
	suite.App.Keepers.FeeShareKeeper.SetFeeShare(suite.Ctx, types.FeeShare{
		ContractAddress:   "terra1mdpvgjc8jmv60a4x68nggsh9w8uyv69sqls04a76m9med5hsqmwsse8sxa",
		DeployerAddress:   "",
		WithdrawerAddress: "terra1zdpgj8am5nqqvht927k3etljyl6a52kwqup0je",
	})
	// ... create the store key ...
	suite.App.Keepers.WasmKeeper.SetExecutedContractAddresses(suite.Ctx, customwasmtypes.ExecutedContracts{
		ContractAddresses: []string{"invalid_contract_addr"},
	})

	// build a tx with a fee amount ...
	txFee := sdk.NewCoins(sdk.NewCoin("uluna", sdk.NewInt(500)), sdk.NewCoin("utoken", sdk.NewInt(250)))
	txBuilder := suite.EncodingConfig.TxConfig.NewTxBuilder()
	txBuilder.SetFeeAmount(txFee)
	txBuilder.SetMsgs(&wasmtypes.MsgExecuteContract{})

	// ... create the feeshare post handler ...
	handler := post.NewFeeSharePayoutDecorator(
		suite.App.Keepers.FeeShareKeeper,
		suite.App.Keepers.BankKeeper,
		suite.App.Keepers.WasmKeeper,
	)

	// Execute the PostHandle function
	_, err := handler.PostHandle(
		suite.Ctx,
		txBuilder.GetTx(),
		false,
		true,
		func(ctx sdk.Context, tx sdk.Tx, simulate bool, success bool) (sdk.Context, error) {
			return mockedPostDecorator.PostHandle(ctx, tx, simulate, success, nil)
		},
	)
	suite.
		Require().
		ErrorIs(err, errorsmod.Wrapf(sdkerrors.ErrLogic, err.Error()))
}
