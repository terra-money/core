package keeper_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/terra-money/core/v2/x/tokenfactory/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func TestKeeperMsgServer(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

// TestMintDenomMsg tests TypeMsgMint message is emitted on a successful mint
func (s *KeeperTestSuite) TestMintDenomMsg() {
	// Create a denom
	res, _ := s.msgServer.CreateDenom(sdk.WrapSDKContext(s.Ctx), types.NewMsgCreateDenom(s.TestAccs[0].String(), "bitcoin"))
	defaultDenom := res.GetNewTokenDenom()

	for _, tc := range []struct {
		desc                  string
		amount                int64
		mintDenom             string
		admin                 string
		valid                 bool
		expectedMessageEvents int
		expectedResType       interface{}
	}{
		{
			desc:      "denom does not exist",
			amount:    10,
			mintDenom: "factory/osmo1t7egva48prqmzl59x5ngv4zx0dtrwewc9m7z44/evmos",
			admin:     s.TestAccs[0].String(),
			valid:     false,
		},
		{
			desc:                  "success case",
			amount:                10,
			mintDenom:             defaultDenom,
			admin:                 s.TestAccs[0].String(),
			valid:                 true,
			expectedMessageEvents: 2,
			expectedResType:       &types.MsgMintResponse{},
		},
	} {
		s.Run(fmt.Sprintf("Case %s", tc.desc), func() {
			ctx := s.Ctx.WithEventManager(sdk.NewEventManager())
			s.Require().Equal(0, len(ctx.EventManager().Events()))
			// Test mint message
			mint := types.NewMsgMint(tc.admin, sdk.NewInt64Coin(tc.mintDenom, 10))
			err := mint.ValidateBasic()
			s.Require().NoError(err)
			res, err := s.msgServer.Mint(sdk.WrapSDKContext(ctx), mint)

			mintTo := types.NewMsgMintTo(tc.admin, sdk.NewInt64Coin(tc.mintDenom, 10), tc.admin)
			err2 := mintTo.ValidateBasic()
			s.Require().NoError(err2)
			res2, err2 := s.msgServer.Mint(sdk.WrapSDKContext(ctx), mintTo)

			if tc.valid {
				s.Require().NoError(err)
				s.Require().NoError(err2)
				s.Require().Equal(tc.expectedResType, res)
				s.Require().Equal(tc.expectedResType, res2)
			} else {
				s.Require().Error(err)
				s.Require().Error(err2)
				s.Require().Nil(res)
				s.Require().Nil(res2)
			}
			// Ensure current number and type of event is emitted
			s.AssertEventEmitted(ctx, types.TypeMsgMint, tc.expectedMessageEvents)
		})
	}
}

// TestForceTransferMsg tests MsgForceTransfer message is emitted on a successful send
func (s *KeeperTestSuite) TestForceTransferMsg() {
	// Create a denom
	admin := s.TestAccs[0].String()
	res, err := s.msgServer.CreateDenom(s.Ctx, types.NewMsgCreateDenom(admin, "bitcoin"))
	s.Require().NoError(err)

	// Mint tokens
	defaultDenom := res.GetNewTokenDenom()
	_, err = s.msgServer.Mint(s.Ctx, types.NewMsgMint(admin, sdk.NewInt64Coin(res.NewTokenDenom, 10)))
	s.Require().NoError(err)

	for _, tc := range []struct {
		desc                  string
		amount                int64
		admin                 string
		transferTo            string
		valid                 bool
		expectedMessageEvents int
		expectedResType       interface{}
	}{
		{
			desc:                  "success case",
			amount:                1,
			admin:                 admin,
			transferTo:            s.TestAccs[1].String(),
			valid:                 true,
			expectedMessageEvents: 1,
			expectedResType:       &types.MsgForceTransferResponse{},
		},
		{
			desc:                  "fail to transfer because user that force transfer is not the admin",
			amount:                1,
			admin:                 s.TestAccs[1].String(),
			transferTo:            admin,
			valid:                 false,
			expectedMessageEvents: 0,
		},
	} {
		s.Run(fmt.Sprintf("Case %s", tc.desc), func() {
			ctx := s.Ctx.WithEventManager(sdk.NewEventManager())
			s.Require().Equal(0, len(ctx.EventManager().Events()))
			// Test force transfer message
			msg := types.NewMsgForceTransfer(tc.admin, sdk.NewInt64Coin(defaultDenom, tc.amount), tc.admin, tc.transferTo)
			err := msg.ValidateBasic()
			s.Require().NoError(err)

			res, err := s.msgServer.ForceTransfer(ctx, msg)
			if tc.valid {
				s.Require().NoError(err)
				s.Require().Equal(tc.expectedResType, res)
			} else {
				s.Require().Error(err)
				s.Require().Nil(res)
			}
			// Ensure current number and type of event is emitted
			s.AssertEventEmitted(ctx, types.TypeMsgForceTransfer, tc.expectedMessageEvents)
		})
	}
}

// TestBurnDenomMsg tests TypeMsgBurn message is emitted on a successful burn
func (s *KeeperTestSuite) TestBurnDenomMsg() {
	// Create a denom
	res, _ := s.msgServer.CreateDenom(sdk.WrapSDKContext(s.Ctx), types.NewMsgCreateDenom(s.TestAccs[0].String(), "bitcoin"))
	defaultDenom := res.GetNewTokenDenom()

	// mint 10 default token for testAcc[0]
	_, err := s.msgServer.Mint(sdk.WrapSDKContext(s.Ctx), types.NewMsgMint(s.TestAccs[0].String(), sdk.NewInt64Coin(defaultDenom, 10)))
	s.Require().NoError(err)

	for _, tc := range []struct {
		desc                  string
		amount                int64
		burnDenom             string
		admin                 string
		valid                 bool
		expectedMessageEvents int
	}{
		{
			desc:      "denom does not exist",
			burnDenom: "factory/osmo1t7egva48prqmzl59x5ngv4zx0dtrwewc9m7z44/evmos",
			admin:     s.TestAccs[0].String(),
			valid:     false,
		},
		{
			desc:                  "success case",
			burnDenom:             defaultDenom,
			admin:                 s.TestAccs[0].String(),
			valid:                 true,
			expectedMessageEvents: 1,
		},
	} {
		s.Run(fmt.Sprintf("Case %s", tc.desc), func() {
			ctx := s.Ctx.WithEventManager(sdk.NewEventManager())
			s.Require().Equal(0, len(ctx.EventManager().Events()))
			// Test burn message
			_, err := s.msgServer.Burn(sdk.WrapSDKContext(ctx), types.NewMsgBurn(tc.admin, sdk.NewInt64Coin(tc.burnDenom, 10)))
			if tc.valid {
				s.Require().NoError(err)
			} else {
				s.Require().Error(err)
			}
			// Ensure current number and type of event is emitted
			s.AssertEventEmitted(ctx, types.TypeMsgBurn, tc.expectedMessageEvents)
		})
	}
}

// TestCreateDenomMsg tests TypeMsgCreateDenom message is emitted on a successful denom creation
func (s *KeeperTestSuite) TestCreateDenomMsg() {
	for _, tc := range []struct {
		desc                  string
		subdenom              string
		valid                 bool
		expectedMessageEvents int
	}{
		{
			desc:     "subdenom too long",
			subdenom: "assadsadsadasdasdsadsadsadsadsadsadsklkadaskkkdasdasedskhanhassyeunganassfnlksdflksafjlkasd",
			valid:    false,
		},
		{
			desc:                  "success case: defaultDenomCreationFee",
			subdenom:              "evmos",
			valid:                 true,
			expectedMessageEvents: 1,
		},
	} {
		s.SetupTest()
		s.Run(fmt.Sprintf("Case %s", tc.desc), func() {
			ctx := s.Ctx.WithEventManager(sdk.NewEventManager())
			s.Require().Equal(0, len(ctx.EventManager().Events()))
			// Set denom creation fee in params
			// Test create denom message
			_, err := s.msgServer.CreateDenom(sdk.WrapSDKContext(ctx), types.NewMsgCreateDenom(s.TestAccs[0].String(), tc.subdenom))
			if tc.valid {
				s.Require().NoError(err)
			} else {
				s.Require().Error(err)
			}
			// Ensure current number and type of event is emitted
			s.AssertEventEmitted(ctx, types.TypeMsgCreateDenom, tc.expectedMessageEvents)
		})
	}
}

// TestChangeAdminDenomMsg tests TypeMsgChangeAdmin message is emitted on a successful admin change
func (s *KeeperTestSuite) TestChangeAdminDenomMsg() {
	for _, tc := range []struct {
		desc                    string
		msgChangeAdmin          func(denom string) *types.MsgChangeAdmin
		expectedChangeAdminPass bool
		expectedAdminIndex      int
		msgMint                 func(denom string) *types.MsgMint
		expectedMintPass        bool
		expectedMessageEvents   int
	}{
		{
			desc: "non-admins can't change the existing admin",
			msgChangeAdmin: func(denom string) *types.MsgChangeAdmin {
				return types.NewMsgChangeAdmin(s.TestAccs[1].String(), denom, s.TestAccs[2].String())
			},
			expectedChangeAdminPass: false,
			expectedAdminIndex:      0,
		},
		{
			desc: "success change admin",
			msgChangeAdmin: func(denom string) *types.MsgChangeAdmin {
				return types.NewMsgChangeAdmin(s.TestAccs[0].String(), denom, s.TestAccs[1].String())
			},
			expectedAdminIndex:      1,
			expectedChangeAdminPass: true,
			expectedMessageEvents:   1,
			msgMint: func(denom string) *types.MsgMint {
				return types.NewMsgMint(s.TestAccs[1].String(), sdk.NewInt64Coin(denom, 5))
			},
			expectedMintPass: true,
		},
	} {
		s.Run(fmt.Sprintf("Case %s", tc.desc), func() {
			// setup test
			s.SetupTest()
			ctx := s.Ctx.WithEventManager(sdk.NewEventManager())
			s.Require().Equal(0, len(ctx.EventManager().Events()))
			// Create a denom and mint
			res, err := s.msgServer.CreateDenom(sdk.WrapSDKContext(ctx), types.NewMsgCreateDenom(s.TestAccs[0].String(), "bitcoin"))
			s.Require().NoError(err)
			testDenom := res.GetNewTokenDenom()
			_, err = s.msgServer.Mint(sdk.WrapSDKContext(ctx), types.NewMsgMint(s.TestAccs[0].String(), sdk.NewInt64Coin(testDenom, 10)))
			s.Require().NoError(err)
			// Test change admin message
			_, err = s.msgServer.ChangeAdmin(sdk.WrapSDKContext(ctx), tc.msgChangeAdmin(testDenom))
			if tc.expectedChangeAdminPass {
				s.Require().NoError(err)
			} else {
				s.Require().Error(err)
			}
			// Ensure current number and type of event is emitted
			s.AssertEventEmitted(ctx, types.TypeMsgChangeAdmin, tc.expectedMessageEvents)
		})
	}
}

// TestSetDenomMetaDataMsg tests TypeMsgSetDenomMetadata message is emitted on a successful denom metadata change
func (s *KeeperTestSuite) TestSetDenomMetaDataMsg() {
	// setup test
	s.SetupTest()

	// Create a denom
	res, _ := s.msgServer.CreateDenom(sdk.WrapSDKContext(s.Ctx), types.NewMsgCreateDenom(s.TestAccs[0].String(), "bitcoin"))
	defaultDenom := res.GetNewTokenDenom()

	for _, tc := range []struct {
		desc                  string
		msgSetDenomMetadata   types.MsgSetDenomMetadata
		expectedPass          bool
		expectedMessageEvents int
	}{
		{
			desc: "successful set denom metadata",
			msgSetDenomMetadata: *types.NewMsgSetDenomMetadata(s.TestAccs[0].String(), banktypes.Metadata{
				Description: "yeehaw",
				DenomUnits: []*banktypes.DenomUnit{
					{
						Denom:    defaultDenom,
						Exponent: 0,
					},
					{
						Denom:    "uosmo",
						Exponent: 6,
					},
				},
				Base:    defaultDenom,
				Display: "uosmo",
				Name:    "OSMO",
				Symbol:  "OSMO",
			}),
			expectedPass:          true,
			expectedMessageEvents: 1,
		},
		{
			desc: "non existent factory denom name",
			msgSetDenomMetadata: *types.NewMsgSetDenomMetadata(s.TestAccs[0].String(), banktypes.Metadata{
				Description: "yeehaw",
				DenomUnits: []*banktypes.DenomUnit{
					{
						Denom:    fmt.Sprintf("factory/%s/litecoin", s.TestAccs[0].String()),
						Exponent: 0,
					},
					{
						Denom:    "uosmo",
						Exponent: 6,
					},
				},
				Base:    fmt.Sprintf("factory/%s/litecoin", s.TestAccs[0].String()),
				Display: "uosmo",
				Name:    "OSMO",
				Symbol:  "OSMO",
			}),
			expectedPass: false,
		},
	} {
		s.Run(fmt.Sprintf("Case %s", tc.desc), func() {
			ctx := s.Ctx.WithEventManager(sdk.NewEventManager())
			s.Require().Equal(0, len(ctx.EventManager().Events()))
			// Test set denom metadata message
			_, err := s.msgServer.SetDenomMetadata(sdk.WrapSDKContext(ctx), &tc.msgSetDenomMetadata)
			if tc.expectedPass {
				s.Require().NoError(err)
			} else {
				s.Require().Error(err)
			}
			// Ensure current number and type of event is emitted
			s.AssertEventEmitted(ctx, types.TypeMsgSetDenomMetadata, tc.expectedMessageEvents)
		})
	}
}
