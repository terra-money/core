package keeper_test

import (
	"crypto/sha256"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"

	_ "embed"

	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/terra-money/core/v2/app/config"
	"github.com/terra-money/core/v2/x/feeshare/types"
)

//go:embed testdata/reflect.wasm
var wasmContract []byte

func (s *IntegrationTestSuite) StoreCode() {
	_, _, sender := testdata.KeyTestPubAddr()
	msg := wasmtypes.MsgStoreCodeFixture(func(m *wasmtypes.MsgStoreCode) {
		m.WASMByteCode = wasmContract
		m.Sender = sender.String()
	})
	rsp, err := s.App.MsgServiceRouter().Handler(msg)(s.Ctx, msg)
	s.Require().NoError(err)
	var result wasmtypes.MsgStoreCodeResponse
	s.Require().NoError(s.App.AppCodec().Unmarshal(rsp.Data, &result))
	s.Require().Equal(uint64(1), result.CodeID)
	expHash := sha256.Sum256(wasmContract)
	s.Require().Equal(expHash[:], result.Checksum)
	// and
	info := s.App.WasmKeeper.GetCodeInfo(s.Ctx, 1)
	s.Require().NotNil(info)
	s.Require().Equal(expHash[:], info.CodeHash)
	s.Require().Equal(sender.String(), info.Creator)
	s.Require().Equal(wasmtypes.DefaultParams().InstantiateDefaultPermission.With(sender), info.InstantiateConfig)
}

func (s *IntegrationTestSuite) InstantiateContract(sender string, admin string) string {
	msgStoreCode := wasmtypes.MsgStoreCodeFixture(func(m *wasmtypes.MsgStoreCode) {
		m.WASMByteCode = wasmContract
		m.Sender = sender
	})
	_, err := s.App.MsgServiceRouter().Handler(msgStoreCode)(s.Ctx, msgStoreCode)
	s.Require().NoError(err)

	msgInstantiate := wasmtypes.MsgInstantiateContractFixture(func(m *wasmtypes.MsgInstantiateContract) {
		m.Sender = sender
		m.Admin = admin
		m.Funds = sdk.NewCoins(sdk.NewCoin(config.MicroLuna, sdk.NewInt(1)))
		m.Msg = []byte(`{}`)
	})
	resp, err := s.App.MsgServiceRouter().Handler(msgInstantiate)(s.Ctx, msgInstantiate)
	s.Require().NoError(err)
	var result wasmtypes.MsgInstantiateContractResponse
	s.Require().NoError(s.App.AppCodec().Unmarshal(resp.Data, &result))
	contractInfo := s.App.WasmKeeper.GetContractInfo(s.Ctx, sdk.MustAccAddressFromBech32(result.Address))
	s.Require().Equal(contractInfo.CodeID, uint64(1))
	s.Require().Equal(contractInfo.Admin, admin)
	s.Require().Equal(contractInfo.Creator, sender)

	return result.Address
}

func (s *IntegrationTestSuite) TestGetContractAdminOrCreatorAddress() {
	s.Setup()
	sender := s.TestAccs[0]
	admin := s.TestAccs[1]

	noAdminContractAddress := s.InstantiateContract(sender.String(), "")
	withAdminContractAddress := s.InstantiateContract(sender.String(), admin.String())

	for _, tc := range []struct {
		desc            string
		contractAddress string
		deployerAddress string
		shouldErr       bool
	}{
		{
			desc:            "Success - Creator",
			contractAddress: noAdminContractAddress,
			deployerAddress: sender.String(),
			shouldErr:       false,
		},
		{
			desc:            "Success - Admin",
			contractAddress: withAdminContractAddress,
			deployerAddress: admin.String(),
			shouldErr:       false,
		},
		{
			desc:            "Error - Invalid deployer",
			contractAddress: noAdminContractAddress,
			deployerAddress: "Invalid",
			shouldErr:       true,
		},
	} {
		tc := tc
		s.Run(tc.desc, func() {
			if !tc.shouldErr {
				_, err := s.App.FeeShareKeeper.GetContractAdminOrCreatorAddress(s.Ctx, sdk.MustAccAddressFromBech32(tc.contractAddress), tc.deployerAddress)
				s.Require().NoError(err)
			} else {
				_, err := s.App.FeeShareKeeper.GetContractAdminOrCreatorAddress(s.Ctx, sdk.MustAccAddressFromBech32(tc.contractAddress), tc.deployerAddress)
				s.Require().Error(err)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestRegisterFeeShare() {
	s.Setup()
	sender := s.TestAccs[0]

	gov := s.App.AccountKeeper.GetModuleAddress(govtypes.ModuleName).String()
	govContract := s.InstantiateContract(sender.String(), gov)

	contractAddress := s.InstantiateContract(sender.String(), "")
	contractAddress2 := s.InstantiateContract(contractAddress, contractAddress)

	contractAddress3 := s.InstantiateContract(sender.String(), "")
	subContract := s.InstantiateContract(contractAddress3, contractAddress3)

	_, _, withdrawer := testdata.KeyTestPubAddr()

	for _, tc := range []struct {
		desc      string
		msg       *types.MsgRegisterFeeShare
		resp      *types.MsgRegisterFeeShareResponse
		shouldErr bool
	}{
		{
			desc: "Invalid contract address",
			msg: &types.MsgRegisterFeeShare{
				ContractAddress:   "Invalid",
				DeployerAddress:   sender.String(),
				WithdrawerAddress: withdrawer.String(),
			},
			resp:      &types.MsgRegisterFeeShareResponse{},
			shouldErr: true,
		},
		{
			desc: "Invalid deployer address",
			msg: &types.MsgRegisterFeeShare{
				ContractAddress:   contractAddress,
				DeployerAddress:   "Invalid",
				WithdrawerAddress: withdrawer.String(),
			},
			resp:      &types.MsgRegisterFeeShareResponse{},
			shouldErr: true,
		},
		{
			desc: "Invalid withdrawer address",
			msg: &types.MsgRegisterFeeShare{
				ContractAddress:   contractAddress,
				DeployerAddress:   sender.String(),
				WithdrawerAddress: "Invalid",
			},
			resp:      &types.MsgRegisterFeeShareResponse{},
			shouldErr: true,
		},
		{
			desc: "Success",
			msg: &types.MsgRegisterFeeShare{
				ContractAddress:   contractAddress,
				DeployerAddress:   sender.String(),
				WithdrawerAddress: withdrawer.String(),
			},
			resp:      &types.MsgRegisterFeeShareResponse{},
			shouldErr: false,
		},
		{
			desc: "Invalid withdraw address for factory contract",
			msg: &types.MsgRegisterFeeShare{
				ContractAddress:   contractAddress2,
				DeployerAddress:   sender.String(),
				WithdrawerAddress: sender.String(),
			},
			resp:      &types.MsgRegisterFeeShareResponse{},
			shouldErr: true,
		},
		{
			desc: "Success register factory contract to itself",
			msg: &types.MsgRegisterFeeShare{
				ContractAddress:   contractAddress2,
				DeployerAddress:   sender.String(),
				WithdrawerAddress: contractAddress2,
			},
			resp:      &types.MsgRegisterFeeShareResponse{},
			shouldErr: false,
		},
		{
			desc: "Invalid register gov contract withdraw address",
			msg: &types.MsgRegisterFeeShare{
				ContractAddress:   govContract,
				DeployerAddress:   sender.String(),
				WithdrawerAddress: sender.String(),
			},
			resp:      &types.MsgRegisterFeeShareResponse{},
			shouldErr: true,
		},
		{
			desc: "Success register gov contract withdraw address to self",
			msg: &types.MsgRegisterFeeShare{
				ContractAddress:   govContract,
				DeployerAddress:   sender.String(),
				WithdrawerAddress: govContract,
			},
			resp:      &types.MsgRegisterFeeShareResponse{},
			shouldErr: false,
		},
		{
			desc: "Success register contract from contractAddress3 contract as admin",
			msg: &types.MsgRegisterFeeShare{
				ContractAddress:   subContract,
				DeployerAddress:   contractAddress3,
				WithdrawerAddress: contractAddress3,
			},
			resp:      &types.MsgRegisterFeeShareResponse{},
			shouldErr: false,
		},
	} {
		s.Run(tc.desc, func() {
			if !tc.shouldErr {
				resp, err := s.App.FeeShareKeeper.RegisterFeeShare(s.Ctx, tc.msg)
				s.Require().NoError(err)
				s.Require().Equal(resp, tc.resp)
			} else {
				resp, err := s.App.FeeShareKeeper.RegisterFeeShare(s.Ctx, tc.msg)
				s.Require().Error(err)
				s.Require().Nil(resp)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestUpdateFeeShare() {
	s.Setup()
	sender := s.TestAccs[0]

	contractAddress := s.InstantiateContract(sender.String(), "")
	_, _, withdrawer := testdata.KeyTestPubAddr()

	contractAddressNoRegisFeeShare := s.InstantiateContract(sender.String(), "")
	s.Require().NotEqual(contractAddress, contractAddressNoRegisFeeShare)

	// RegsisFeeShare
	goCtx := sdk.WrapSDKContext(s.Ctx)
	msg := &types.MsgRegisterFeeShare{
		ContractAddress:   contractAddress,
		DeployerAddress:   sender.String(),
		WithdrawerAddress: withdrawer.String(),
	}
	_, err := s.App.FeeShareKeeper.RegisterFeeShare(goCtx, msg)
	s.Require().NoError(err)
	_, _, newWithdrawer := testdata.KeyTestPubAddr()
	s.Require().NotEqual(withdrawer, newWithdrawer)

	for _, tc := range []struct {
		desc      string
		msg       *types.MsgUpdateFeeShare
		resp      *types.MsgCancelFeeShareResponse
		shouldErr bool
	}{
		{
			desc: "Invalid - contract not regis",
			msg: &types.MsgUpdateFeeShare{
				ContractAddress:   contractAddressNoRegisFeeShare,
				DeployerAddress:   sender.String(),
				WithdrawerAddress: newWithdrawer.String(),
			},
			resp:      nil,
			shouldErr: true,
		},
		{
			desc: "Invalid - Invalid DeployerAddress",
			msg: &types.MsgUpdateFeeShare{
				ContractAddress:   contractAddress,
				DeployerAddress:   "Invalid",
				WithdrawerAddress: newWithdrawer.String(),
			},
			resp:      nil,
			shouldErr: true,
		},
		{
			desc: "Invalid - Invalid WithdrawerAddress",
			msg: &types.MsgUpdateFeeShare{
				ContractAddress:   contractAddress,
				DeployerAddress:   sender.String(),
				WithdrawerAddress: "Invalid",
			},
			resp:      nil,
			shouldErr: true,
		},
		{
			desc: "Invalid - Invalid WithdrawerAddress not change",
			msg: &types.MsgUpdateFeeShare{
				ContractAddress:   contractAddress,
				DeployerAddress:   sender.String(),
				WithdrawerAddress: withdrawer.String(),
			},
			resp:      nil,
			shouldErr: true,
		},
		{
			desc: "Success",
			msg: &types.MsgUpdateFeeShare{
				ContractAddress:   contractAddress,
				DeployerAddress:   sender.String(),
				WithdrawerAddress: newWithdrawer.String(),
			},
			resp:      &types.MsgCancelFeeShareResponse{},
			shouldErr: false,
		},
	} {
		tc := tc
		s.Run(tc.desc, func() {
			goCtx := sdk.WrapSDKContext(s.Ctx)
			if !tc.shouldErr {
				_, err := s.App.FeeShareKeeper.UpdateFeeShare(goCtx, tc.msg)
				s.Require().NoError(err)
			} else {
				resp, err := s.App.FeeShareKeeper.UpdateFeeShare(goCtx, tc.msg)
				s.Require().Error(err)
				s.Require().Nil(resp)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestCancelFeeShare() {
	s.AppTestSuite.Setup()
	sender := s.AppTestSuite.TestAccs[0]

	contractAddress := s.InstantiateContract(sender.String(), "")
	_, _, withdrawer := testdata.KeyTestPubAddr()

	// RegsisFeeShare
	goCtx := sdk.WrapSDKContext(s.Ctx)
	msg := &types.MsgRegisterFeeShare{
		ContractAddress:   contractAddress,
		DeployerAddress:   sender.String(),
		WithdrawerAddress: withdrawer.String(),
	}
	_, err := s.App.FeeShareKeeper.RegisterFeeShare(goCtx, msg)
	s.Require().NoError(err)

	for _, tc := range []struct {
		desc      string
		msg       *types.MsgCancelFeeShare
		resp      *types.MsgCancelFeeShareResponse
		shouldErr bool
	}{
		{
			desc: "Invalid - contract Address",
			msg: &types.MsgCancelFeeShare{
				ContractAddress: "Invalid",
				DeployerAddress: sender.String(),
			},
			resp:      nil,
			shouldErr: true,
		},
		{
			desc: "Invalid - deployer Address",
			msg: &types.MsgCancelFeeShare{
				ContractAddress: contractAddress,
				DeployerAddress: "Invalid",
			},
			resp:      nil,
			shouldErr: true,
		},
		{
			desc: "Success",
			msg: &types.MsgCancelFeeShare{
				ContractAddress: contractAddress,
				DeployerAddress: sender.String(),
			},
			resp:      &types.MsgCancelFeeShareResponse{},
			shouldErr: false,
		},
	} {
		tc := tc
		s.Run(tc.desc, func() {
			goCtx := sdk.WrapSDKContext(s.Ctx)
			if !tc.shouldErr {
				resp, err := s.App.FeeShareKeeper.CancelFeeShare(goCtx, tc.msg)
				s.Require().NoError(err)
				s.Require().Equal(resp, tc.resp)
			} else {
				resp, err := s.App.FeeShareKeeper.CancelFeeShare(goCtx, tc.msg)
				s.Require().Error(err)
				s.Require().Equal(resp, tc.resp)
			}
		})
	}
}
