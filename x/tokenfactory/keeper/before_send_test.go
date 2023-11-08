package keeper_test

import (
	"fmt"
	"os"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/terra-money/core/v2/x/tokenfactory/keeper"
	"github.com/terra-money/core/v2/x/tokenfactory/types"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
)

type SendMsgTestCase struct {
	desc       string
	msg        func(denom string) *banktypes.MsgSend
	expectPass bool
}

func (s *KeeperTestSuite) TestBeforeSendHook() {
	for _, tc := range []struct {
		desc     string
		wasmFile string
		sendMsgs []SendMsgTestCase
	}{
		{
			desc:     "should not allow sending 100 amount of *any* denom",
			wasmFile: "./testdata/no100.wasm",
			sendMsgs: []SendMsgTestCase{
				{
					desc: "sending 1 of factorydenom should not error",
					msg: func(factorydenom string) *banktypes.MsgSend {
						return banktypes.NewMsgSend(
							s.TestAccs[0],
							s.TestAccs[1],
							sdk.NewCoins(sdk.NewInt64Coin(factorydenom, 1)),
						)
					},
					expectPass: true,
				},
				{
					desc: "sending 100 of non-factorydenom should not error",
					msg: func(factorydenom string) *banktypes.MsgSend {
						return banktypes.NewMsgSend(
							s.TestAccs[0],
							s.TestAccs[1],
							sdk.NewCoins(sdk.NewInt64Coin(factorydenom, 1)),
						)
					},
					expectPass: true,
				},
				{
					desc: "sending 100 of factorydenom should not work",
					msg: func(factorydenom string) *banktypes.MsgSend {
						return banktypes.NewMsgSend(
							s.TestAccs[0],
							s.TestAccs[1],
							sdk.NewCoins(sdk.NewInt64Coin("foo", 100)),
						)
					},
					expectPass: false,
				},
				{
					desc: "having 100 coin within coins should not work",
					msg: func(factorydenom string) *banktypes.MsgSend {
						return banktypes.NewMsgSend(
							s.TestAccs[0],
							s.TestAccs[1],
							sdk.NewCoins(sdk.NewInt64Coin(factorydenom, 100), sdk.NewInt64Coin("foo", 1)),
						)
					},
					expectPass: false,
				},
			},
		},
	} {
		s.Run(fmt.Sprintf("Case %s", tc.desc), func() {
			// setup test
			s.SetupTest()

			// upload and instantiate wasm code
			wasmCode, err := os.ReadFile(tc.wasmFile)
			s.Require().NoError(err, "test: %v", tc.desc)
			codeID, _, err := s.contractKeeper.Create(s.Ctx, s.TestAccs[0], wasmCode, &wasmtypes.AccessConfig{Permission: wasmtypes.AccessTypeEverybody})
			s.Require().NoError(err, "test: %v", tc.desc)
			cosmwasmAddress, _, err := s.contractKeeper.Instantiate(s.Ctx, codeID, s.TestAccs[0], s.TestAccs[0], []byte("{}"), "", sdk.NewCoins())
			s.Require().NoError(err, "test: %v", tc.desc)

			// create new denom
			res, err := s.msgServer.CreateDenom(sdk.WrapSDKContext(s.Ctx), types.NewMsgCreateDenom(s.TestAccs[0].String(), "bitcoin"))
			s.Require().NoError(err, "test: %v", tc.desc)
			denom := res.GetNewTokenDenom()

			// mint enough coins to the creator
			_, err = s.msgServer.Mint(sdk.WrapSDKContext(s.Ctx), types.NewMsgMint(s.TestAccs[0].String(), sdk.NewInt64Coin(denom, 1000000000)))
			s.Require().NoError(err)

			// set beforesend hook to the new denom
			msg := types.NewMsgSetBeforeSendHook(s.TestAccs[0].String(), denom, cosmwasmAddress.String())
			err = msg.ValidateBasic()
			s.Require().NoError(err, "test: %v", tc.desc)
			_, err = s.msgServer.SetBeforeSendHook(sdk.WrapSDKContext(s.Ctx), msg)
			s.Require().NoError(err, "test: %v", tc.desc)

			for _, sendTc := range tc.sendMsgs {
				_, err := s.bankMsgServer.Send(sdk.WrapSDKContext(s.Ctx), sendTc.msg(denom))
				if sendTc.expectPass {
					s.Require().NoError(err, "test: %v", sendTc.desc)
				} else {
					s.Require().Error(err, "test: %v", sendTc.desc)
				}
			}
		})
	}
}

// TestInfiniteTrackBeforeSend tests gas metering with infinite loop contract
// to properly test if we are gas metering trackBeforeSend properly.
func (s *KeeperTestSuite) TestInfiniteTrackBeforeSend() {
	for _, tc := range []struct {
		name            string
		wasmFile        string
		tokenToSend     sdk.Coins
		useFactoryDenom bool
		expectedError   bool
	}{
		{
			name:            "sending tokenfactory denom from module to module with infinite contract should panic",
			wasmFile:        "./testdata/infinite_track_beforesend.wasm",
			useFactoryDenom: true,
		},
		{
			name:            "sending non-tokenfactory denom from module to module with infinite contract should not panic",
			wasmFile:        "./testdata/infinite_track_beforesend.wasm",
			tokenToSend:     sdk.NewCoins(sdk.NewInt64Coin("foo", 1000000)),
			useFactoryDenom: false,
		},
		{
			name:            "Try using no 100 ",
			wasmFile:        "./testdata/no100.wasm",
			useFactoryDenom: true,
		},
	} {
		s.Run(fmt.Sprintf("Case %s", tc.name), func() {
			// setup test
			s.SetupTest()

			// load wasm file
			wasmCode, err := os.ReadFile(tc.wasmFile)
			s.Require().NoError(err)

			// instantiate wasm code
			codeID, _, err := s.contractKeeper.Create(s.Ctx, s.TestAccs[0], wasmCode, &wasmtypes.AccessConfig{Permission: wasmtypes.AccessTypeEverybody})
			s.Require().NoError(err, "test: %v", tc.name)
			cosmwasmAddress, _, err := s.contractKeeper.Instantiate(s.Ctx, codeID, s.TestAccs[0], s.TestAccs[0], []byte("{}"), "", sdk.NewCoins())
			s.Require().NoError(err, "test: %v", tc.name)

			// create new denom
			res, err := s.msgServer.CreateDenom(sdk.WrapSDKContext(s.Ctx), types.NewMsgCreateDenom(s.TestAccs[0].String(), "bitcoin"))
			s.Require().NoError(err, "test: %v", tc.name)
			factoryDenom := res.GetNewTokenDenom()

			var tokenToSend sdk.Coins
			if tc.useFactoryDenom {
				tokenToSend = sdk.NewCoins(sdk.NewInt64Coin(factoryDenom, 100))
			} else {
				tokenToSend = tc.tokenToSend
			}

			// send the mint module tokenToSend
			if err := s.App.Keepers.BankKeeper.MintCoins(s.Ctx, minttypes.ModuleName, tokenToSend); err != nil {
				s.Require().NoError(err)
			}

			// set beforesend hook to the new denom
			// we register infinite loop contract here to test if we are gas metering properly
			msg := types.NewMsgSetBeforeSendHook(s.TestAccs[0].String(), factoryDenom, cosmwasmAddress.String())
			err = msg.ValidateBasic()
			s.Require().NoError(err, "test: %v", tc.name)
			_, err = s.msgServer.SetBeforeSendHook(sdk.WrapSDKContext(s.Ctx), msg)
			s.Require().NoError(err, "test: %v", tc.name)

			// track before send suppresses in any case, thus we expect no error
			err = s.App.Keepers.BankKeeper.SendCoinsFromModuleToModule(s.Ctx, "mint", "distribution", tokenToSend)
			s.Require().NoError(err)

			// send should happen regardless of trackBeforeSend results
			distributionModuleAddress := s.App.Keepers.AccountKeeper.GetModuleAddress("distribution")
			distributionModuleBalances := s.App.Keepers.BankKeeper.GetAllBalances(s.Ctx, distributionModuleAddress)
			s.Require().True(distributionModuleBalances[0].IsEqual(tokenToSend[0]))
		})
	}
}

// Test paarse coins from sdk to cosmwasm types
func (s *KeeperTestSuite) TestCoinsFromSDK() {
	parsed_coin := keeper.CWCoinFromSDKCoin(sdk.NewInt64Coin("foo", 100))
	parsed_coins := keeper.CWCoinsFromSDKCoins(sdk.NewCoins(sdk.NewInt64Coin("foo", 100), sdk.NewInt64Coin("bar", 100)))

	expected_coin := wasmvmtypes.Coin{
		Denom:  "foo",
		Amount: "100",
	}

	expected_coins := wasmvmtypes.Coins{
		wasmvmtypes.Coin{
			Denom:  "bar",
			Amount: "100",
		},
		wasmvmtypes.Coin{
			Denom:  "foo",
			Amount: "100",
		},
	}

	s.Require().Equal(expected_coin, parsed_coin)
	s.Require().Equal(expected_coins, parsed_coins)

}

// Test get before send hook with no hook set
func (s *KeeperTestSuite) TestGetBeforeSendHook() {
	res := s.App.Keepers.TokenFactoryKeeper.GetBeforeSendHook(s.Ctx, "foo")

	expected := ""

	s.Require().Equal(expected, res)
}
