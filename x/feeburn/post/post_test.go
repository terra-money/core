package post_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	app "github.com/terra-money/core/v2/app/app_test"
	"github.com/terra-money/core/v2/app/post/mocks"
	post "github.com/terra-money/core/v2/x/feeburn/post"
	"github.com/terra-money/core/v2/x/feeburn/types"
)

type PostTestSuite struct {
	app.AppTestSuite
}

func TestAnteSuite(t *testing.T) {
	suite.Run(t, new(PostTestSuite))
}

func (suite *PostTestSuite) TestCalculateFee() {
	suite.Setup()

	// Create a mocked next post hanlder to assert the function being called.
	ctrl := gomock.NewController(suite.T())
	mockedPostDecorator := mocks.NewMockPostDecorator(ctrl)

	feeBurnPost := post.NewFeeBurnDecorator(
		suite.App.Keepers.FeeBurnKeeper,
		suite.App.Keepers.BankKeeper,
	)
	testCases := []struct {
		name string

		// Input values
		gasLimit uint64
		gasUsed  uint64
		feeTx    func() sdk.FeeTx

		// Expected values
		expectedEvents sdk.Events
	}{
		{
			"Must burn half the transaction fee",
			100_000,
			48_991, // gas used is not exaclty 50k because it uses some gas to read from store and do some maths before burning .
			func() sdk.FeeTx {
				txBuilder := suite.EncodingConfig.TxConfig.NewTxBuilder()
				txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewCoin("uluna", sdk.NewInt(500))))
				return txBuilder.GetTx()
			},
			sdk.Events(
				sdk.Events{
					{
						Type: "coin_spent",
						Attributes: []abci.EventAttribute{
							{
								Key:   "spender",
								Value: "terra17xpfvakm2amg962yls6f84z3kell8c5lkaeqfa",
								Index: false,
							},
							{
								Key:   "amount",
								Value: "250uluna",
								Index: false,
							},
						},
					},
					{
						Type: "burn",
						Attributes: []abci.EventAttribute{
							{
								Key:   "burner",
								Value: "terra17xpfvakm2amg962yls6f84z3kell8c5lkaeqfa",
								Index: false,
							},
							{
								Key:   "amount",
								Value: "250uluna",
								Index: false,
							},
						},
					},
					{
						Type: "terra.feeburn.v1.FeeBurnEvent",
						Attributes: []abci.EventAttribute{
							{
								Key:   "burn_rate",
								Value: "\"0.500000000000000000\"",
								Index: false,
							},
							{
								Key:   "fees_burn",
								Value: "[{\"denom\":\"uluna\",\"amount\":\"250\"}]",
								Index: false,
							},
						},
					},
				}),
		},
		{
			"Must burn half the transaction fee containing two tokens",
			100_000,
			48_991, // gas used is not exaclty 50k because it uses some gas to read from store and do some maths before burning .
			func() sdk.FeeTx {
				txBuilder := suite.EncodingConfig.TxConfig.NewTxBuilder()
				txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewCoin("uluna", sdk.NewInt(500)), sdk.NewCoin("utoken", sdk.NewInt(250))))
				return txBuilder.GetTx()
			},
			sdk.Events(
				sdk.Events{
					{
						Type: "coin_spent",
						Attributes: []abci.EventAttribute{
							{
								Key:   "spender",
								Value: "terra17xpfvakm2amg962yls6f84z3kell8c5lkaeqfa",
								Index: false,
							},
							{
								Key:   "amount",
								Value: "250uluna,125utoken",
								Index: false,
							},
						},
					},
					{
						Type: "burn",
						Attributes: []abci.EventAttribute{
							{
								Key:   "burner",
								Value: "terra17xpfvakm2amg962yls6f84z3kell8c5lkaeqfa",
								Index: false,
							},
							{
								Key:   "amount",
								Value: "250uluna,125utoken",
								Index: false,
							},
						},
					},
					{
						Type: "terra.feeburn.v1.FeeBurnEvent",
						Attributes: []abci.EventAttribute{
							{
								Key:   "burn_rate",
								Value: "\"0.500000000000000000\"",
								Index: false,
							},
							{
								Key:   "fees_burn",
								Value: "[{\"denom\":\"uluna\",\"amount\":\"250\"},{\"denom\":\"utoken\",\"amount\":\"125\"}]",
								Index: false,
							},
						},
					},
				}),
		},
		{
			"Must burn a quarter of the transaction fee containing two tokens rounding down to the nearest integer",
			100_000,
			23_991, // gas used is not exaclty 25k because it uses some gas to read from store and do some maths before burning .
			func() sdk.FeeTx {
				txBuilder := suite.EncodingConfig.TxConfig.NewTxBuilder()
				txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewCoin("uluna", sdk.NewInt(250)), sdk.NewCoin("utoken", sdk.NewInt(125))))
				return txBuilder.GetTx()
			},
			sdk.Events(
				sdk.Events{
					{
						Type: "coin_spent",
						Attributes: []abci.EventAttribute{
							{
								Key:   "spender",
								Value: "terra17xpfvakm2amg962yls6f84z3kell8c5lkaeqfa",
								Index: false,
							},
							{
								Key:   "amount",
								Value: "187uluna,93utoken",
								Index: false,
							},
						},
					},
					{
						Type: "burn",
						Attributes: []abci.EventAttribute{
							{
								Key:   "burner",
								Value: "terra17xpfvakm2amg962yls6f84z3kell8c5lkaeqfa",
								Index: false,
							},
							{
								Key:   "amount",
								Value: "187uluna,93utoken",
								Index: false,
							},
						},
					},
					{
						Type: "terra.feeburn.v1.FeeBurnEvent",
						Attributes: []abci.EventAttribute{
							{
								Key:   "burn_rate",
								Value: "\"0.750000000000000000\"",
								Index: false,
							},
							{
								Key:   "fees_burn",
								Value: "[{\"denom\":\"uluna\",\"amount\":\"187\"},{\"denom\":\"utoken\",\"amount\":\"93\"}]",
								Index: false,
							},
						},
					},
				}),
		},
		{
			"Must define zero fees to skip the burning process",
			5_000,
			0,
			func() sdk.FeeTx {
				txBuilder := suite.EncodingConfig.TxConfig.NewTxBuilder()
				txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewCoin("uluna", sdk.NewInt(0))))
				return txBuilder.GetTx()
			},
			sdk.EmptyEvents(),
		},
		{
			"Must disable the module to skip the burning process",
			5_000,
			0,
			func() sdk.FeeTx {
				suite.App.Keepers.FeeBurnKeeper.SetParams(suite.Ctx, types.Params{EnableFeeBurn: false})
				return suite.EncodingConfig.TxConfig.NewTxBuilder().GetTx()
			},
			sdk.EmptyEvents(),
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			// Setup the test case, removing the events from the previous
			// gas meters set the gas limit and gas used.
			feeTx := tc.feeTx()
			suite.Ctx = suite.Ctx.WithEventManager(sdk.NewEventManager())
			suite.Ctx = suite.Ctx.WithGasMeter(sdk.NewGasMeter(tc.gasLimit))
			suite.Ctx.GasMeter().ConsumeGas(tc.gasUsed, "test")

			// assert the next hanlder is called once
			mockedPostDecorator.
				EXPECT().
				PostHandle(gomock.Any(), gomock.Any(), false, true, gomock.Any()).
				Times(1)

			_, err := feeBurnPost.PostHandle(suite.Ctx,
				feeTx,
				false,
				true,
				func(ctx sdk.Context, tx sdk.Tx, simulate bool, success bool) (sdk.Context, error) {
					// Overwrite the context with the context returned from the handler.
					suite.Ctx = ctx
					return mockedPostDecorator.PostHandle(ctx, tx, simulate, success, nil)
				},
			)

			suite.Require().NoError(err)
			suite.Require().NotNil(suite.Ctx)
			suite.Require().Equal(
				tc.expectedEvents,
				suite.Ctx.EventManager().Events(),
			)
		})
	}
}
