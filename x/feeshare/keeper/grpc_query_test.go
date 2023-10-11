package keeper_test

import (
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	app "github.com/terra-money/core/v2/app/app_test"
	"github.com/terra-money/core/v2/x/feeshare/types"
)

func (s *IntegrationTestSuite) TestFeeShares() {
	s.Setup()
	sender := s.AppTestSuite.TestAccs[0]
	_ = s.FundAccount(s.Ctx, sender, sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(1_000_000))))

	_, _, withdrawer := testdata.KeyTestPubAddr()

	var contractAddressList []string
	var index uint64
	for index < 5 {
		contractAddress := s.InstantiateContract(sender.String(), "")
		contractAddressList = append(contractAddressList, contractAddress)
		index++
	}

	// RegsisFeeShare
	var feeShares []types.FeeShare
	for _, contractAddress := range contractAddressList {
		goCtx := sdk.WrapSDKContext(s.Ctx)
		msg := &types.MsgRegisterFeeShare{
			ContractAddress:   contractAddress,
			DeployerAddress:   sender.String(),
			WithdrawerAddress: withdrawer.String(),
		}

		feeShare := types.FeeShare{
			ContractAddress:   contractAddress,
			DeployerAddress:   sender.String(),
			WithdrawerAddress: withdrawer.String(),
		}

		feeShares = append(feeShares, feeShare)

		_, err := s.feeShareMsgServer.RegisterFeeShare(goCtx, msg)
		s.Require().NoError(err)
	}

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryFeeSharesRequest {
		return &types.QueryFeeSharesRequest{
			Pagination: &query.PageRequest{
				Key:        next,
				Offset:     offset,
				Limit:      limit,
				CountTotal: total,
			},
		}
	}
	s.Run("ByOffset", func() {
		step := 2
		goCtx := sdk.WrapSDKContext(s.Ctx)
		for i := 0; i < len(contractAddressList); i += step {
			resp, err := s.queryClient.FeeShares(goCtx, request(nil, uint64(i), uint64(step), false))
			s.Require().NoError(err)
			s.Require().LessOrEqual(len(resp.Feeshare), step)
			s.Require().Subset(app.Fill(feeShares), app.Fill(resp.Feeshare))
		}
	})
	s.Run("ByKey", func() {
		step := 2
		var next []byte
		goCtx := sdk.WrapSDKContext(s.Ctx)
		for i := 0; i < len(contractAddressList); i += step {
			resp, err := s.queryClient.FeeShares(goCtx, request(next, 0, uint64(step), false))
			s.Require().NoError(err)
			s.Require().LessOrEqual(len(resp.Feeshare), step)
			s.Require().Subset(app.Fill(feeShares), app.Fill(resp.Feeshare))
			next = resp.Pagination.NextKey
		}
	})
	s.Run("Total", func() {
		goCtx := sdk.WrapSDKContext(s.Ctx)
		resp, err := s.queryClient.FeeShares(goCtx, request(nil, 0, 0, true))
		s.Require().NoError(err)
		s.Require().Equal(len(feeShares), int(resp.Pagination.Total))
		s.Require().ElementsMatch(app.Fill(feeShares), app.Fill(resp.Feeshare))
	})
}

func (s *IntegrationTestSuite) TestFeeShare() {
	s.AppTestSuite.Setup()
	sender := s.AppTestSuite.TestAccs[0]
	err := s.FundAccount(s.Ctx, sender, sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(1_000_000))))
	s.Require().NoError(err)

	_, _, withdrawer := testdata.KeyTestPubAddr()

	contractAddress := s.InstantiateContract(sender.String(), "")
	goCtx := sdk.WrapSDKContext(s.Ctx)
	msg := &types.MsgRegisterFeeShare{
		ContractAddress:   contractAddress,
		DeployerAddress:   sender.String(),
		WithdrawerAddress: withdrawer.String(),
	}

	feeShare := types.FeeShare{
		ContractAddress:   contractAddress,
		DeployerAddress:   sender.String(),
		WithdrawerAddress: withdrawer.String(),
	}
	_, err = s.feeShareMsgServer.RegisterFeeShare(goCtx, msg)
	s.Require().NoError(err)

	req := &types.QueryFeeShareRequest{
		ContractAddress: contractAddress,
	}
	goCtx = sdk.WrapSDKContext(s.Ctx)
	resp, err := s.queryClient.FeeShare(goCtx, req)
	s.Require().NoError(err)
	s.Require().Equal(resp.Feeshare, feeShare)
}

func (s *IntegrationTestSuite) TestDeployerFeeShares() {
	s.AppTestSuite.Setup()
	sender := s.AppTestSuite.TestAccs[0]
	err := s.FundAccount(s.Ctx, sender, sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(1_000_000))))
	s.Require().NoError(err)

	_, _, withdrawer := testdata.KeyTestPubAddr()

	var contractAddressList []string
	var index uint64
	for index < 5 {
		contractAddress := s.InstantiateContract(sender.String(), "")
		contractAddressList = append(contractAddressList, contractAddress)
		index++
	}

	// RegsisFeeShare
	for _, contractAddress := range contractAddressList {
		goCtx := sdk.WrapSDKContext(s.Ctx)
		msg := &types.MsgRegisterFeeShare{
			ContractAddress:   contractAddress,
			DeployerAddress:   sender.String(),
			WithdrawerAddress: withdrawer.String(),
		}

		_, err := s.feeShareMsgServer.RegisterFeeShare(goCtx, msg)
		s.Require().NoError(err)
	}

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryDeployerFeeSharesRequest {
		return &types.QueryDeployerFeeSharesRequest{
			DeployerAddress: sender.String(),
			Pagination: &query.PageRequest{
				Key:        next,
				Offset:     offset,
				Limit:      limit,
				CountTotal: total,
			},
		}
	}
	s.Run("ByOffset", func() {
		step := 2
		goCtx := sdk.WrapSDKContext(s.Ctx)
		for i := 0; i < len(contractAddressList); i += step {
			resp, err := s.queryClient.DeployerFeeShares(goCtx, request(nil, uint64(i), uint64(step), false))
			s.Require().NoError(err)
			s.Require().LessOrEqual(len(resp.ContractAddresses), step)
			s.Require().Subset(app.Fill(contractAddressList), app.Fill(resp.ContractAddresses))
		}
	})
	s.Run("ByKey", func() {
		step := 2
		var next []byte
		goCtx := sdk.WrapSDKContext(s.Ctx)
		for i := 0; i < len(contractAddressList); i += step {
			resp, err := s.queryClient.DeployerFeeShares(goCtx, request(next, 0, uint64(step), false))
			s.Require().NoError(err)
			s.Require().LessOrEqual(len(resp.ContractAddresses), step)
			s.Require().Subset(app.Fill(contractAddressList), app.Fill(resp.ContractAddresses))
			next = resp.Pagination.NextKey
		}
	})
	s.Run("Total", func() {
		goCtx := sdk.WrapSDKContext(s.Ctx)
		resp, err := s.queryClient.DeployerFeeShares(goCtx, request(nil, 0, 0, true))
		s.Require().NoError(err)
		s.Require().Equal(len(contractAddressList), int(resp.Pagination.Total))
		s.Require().ElementsMatch(app.Fill(contractAddressList), app.Fill(resp.ContractAddresses))
	})
}

func (s *IntegrationTestSuite) TestWithdrawerFeeShares() {
	s.AppTestSuite.Setup()

	sender := s.AppTestSuite.TestAccs[0]
	err := s.FundAccount(s.Ctx, sender, sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(1_000_000))))
	s.Require().NoError(err)

	_, _, withdrawer := testdata.KeyTestPubAddr()

	var contractAddressList []string
	var index uint64
	for index < 5 {
		contractAddress := s.InstantiateContract(sender.String(), "")
		contractAddressList = append(contractAddressList, contractAddress)
		index++
	}

	// RegsisFeeShare
	for _, contractAddress := range contractAddressList {
		goCtx := sdk.WrapSDKContext(s.Ctx)
		msg := &types.MsgRegisterFeeShare{
			ContractAddress:   contractAddress,
			DeployerAddress:   sender.String(),
			WithdrawerAddress: withdrawer.String(),
		}

		_, err := s.feeShareMsgServer.RegisterFeeShare(goCtx, msg)
		s.Require().NoError(err)
	}

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryWithdrawerFeeSharesRequest {
		return &types.QueryWithdrawerFeeSharesRequest{
			WithdrawerAddress: withdrawer.String(),
			Pagination: &query.PageRequest{
				Key:        next,
				Offset:     offset,
				Limit:      limit,
				CountTotal: total,
			},
		}
	}
	s.Run("ByOffset", func() {
		step := 2
		goCtx := sdk.WrapSDKContext(s.Ctx)
		for i := 0; i < len(contractAddressList); i += step {
			resp, err := s.queryClient.WithdrawerFeeShares(goCtx, request(nil, uint64(i), uint64(step), false))
			s.Require().NoError(err)
			s.Require().LessOrEqual(len(resp.ContractAddresses), step)
			s.Require().Subset(app.Fill(contractAddressList), app.Fill(resp.ContractAddresses))
		}
	})
	s.Run("ByKey", func() {
		step := 2
		var next []byte
		goCtx := sdk.WrapSDKContext(s.Ctx)
		for i := 0; i < len(contractAddressList); i += step {
			resp, err := s.queryClient.WithdrawerFeeShares(goCtx, request(next, 0, uint64(step), false))
			s.Require().NoError(err)
			s.Require().LessOrEqual(len(resp.ContractAddresses), step)
			s.Require().Subset(app.Fill(contractAddressList), app.Fill(resp.ContractAddresses))
			next = resp.Pagination.NextKey
		}
	})
	s.Run("Total", func() {
		goCtx := sdk.WrapSDKContext(s.Ctx)
		resp, err := s.queryClient.WithdrawerFeeShares(goCtx, request(nil, 0, 0, true))
		s.Require().NoError(err)
		s.Require().Equal(len(contractAddressList), int(resp.Pagination.Total))
		s.Require().ElementsMatch(app.Fill(contractAddressList), app.Fill(resp.ContractAddresses))
	})
}
