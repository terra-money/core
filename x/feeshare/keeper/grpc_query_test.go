package keeper_test

import (
	"github.com/terra-money/core/v2/x/feeshare/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
)

func (s *IntegrationTestSuite) TestFeeShares() {
	s.SetupTest()
	sender := s.TestAccs[0]
	withdrawer := s.TestAccs[1]

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

		_, err := s.App.Keepers.FeeShareKeeper.RegisterFeeShare(s.Ctx, msg)
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
		goCtx := sdk.WrapSDKContext(s.Ctx)
		step := 2
		for i := 0; i < len(contractAddressList); i += step {
			resp, err := s.queryClient.FeeShares(goCtx, request(nil, uint64(i), uint64(step), false))
			s.Require().NoError(err)
			s.Require().LessOrEqual(len(resp.Feeshare), step)
			s.Require().Subset(feeShares, resp.Feeshare)
		}
	})
	s.Run("ByKey", func() {
		goCtx := sdk.WrapSDKContext(s.Ctx)
		step := 2
		var next []byte
		for i := 0; i < len(contractAddressList); i += step {
			resp, err := s.queryClient.FeeShares(goCtx, request(next, 0, uint64(step), false))
			s.Require().NoError(err)
			s.Require().LessOrEqual(len(resp.Feeshare), step)
			s.Require().Subset(feeShares, resp.Feeshare)
			next = resp.Pagination.NextKey
		}
	})
	s.Run("Total", func() {
		goCtx := sdk.WrapSDKContext(s.Ctx)
		resp, err := s.queryClient.FeeShares(goCtx, request(nil, 0, 0, true))
		s.Require().NoError(err)
		s.Require().Equal(len(feeShares), int(resp.Pagination.Total))
		s.Require().ElementsMatch(feeShares, resp.Feeshare)
	})
}

func (s *IntegrationTestSuite) TestFeeShare() {
	s.SetupTest()
	sender := s.TestAccs[0]
	withdrawer := s.TestAccs[1]

	contractAddress := s.InstantiateContract(sender.String(), "")
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
	_, err := s.App.Keepers.FeeShareKeeper.RegisterFeeShare(s.Ctx, msg)
	s.Require().NoError(err)

	req := &types.QueryFeeShareRequest{
		ContractAddress: contractAddress,
	}
	resp, err := s.queryClient.FeeShare(s.Ctx, req)
	s.Require().NoError(err)
	s.Require().Equal(resp.Feeshare, feeShare)
}

func (s *IntegrationTestSuite) TestDeployerFeeShares() {
	s.SetupTest()
	sender := s.TestAccs[0]
	withdrawer := s.TestAccs[1]

	var contractAddressList []string
	var index uint64
	for index < 5 {
		contractAddress := s.InstantiateContract(sender.String(), "")
		contractAddressList = append(contractAddressList, contractAddress)
		index++
	}

	// RegsisFeeShare
	for _, contractAddress := range contractAddressList {
		msg := &types.MsgRegisterFeeShare{
			ContractAddress:   contractAddress,
			DeployerAddress:   sender.String(),
			WithdrawerAddress: withdrawer.String(),
		}

		_, err := s.App.Keepers.FeeShareKeeper.RegisterFeeShare(s.Ctx, msg)
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
		for i := 0; i < len(contractAddressList); i += step {
			resp, err := s.queryClient.DeployerFeeShares(s.Ctx, request(nil, uint64(i), uint64(step), false))
			s.Require().NoError(err)
			s.Require().LessOrEqual(len(resp.ContractAddresses), step)
			s.Require().Subset(contractAddressList, resp.ContractAddresses)
		}
	})
	s.Run("ByKey", func() {
		step := 2
		var next []byte
		for i := 0; i < len(contractAddressList); i += step {
			resp, err := s.queryClient.DeployerFeeShares(s.Ctx, request(next, 0, uint64(step), false))
			s.Require().NoError(err)
			s.Require().LessOrEqual(len(resp.ContractAddresses), step)
			s.Require().Subset(contractAddressList, resp.ContractAddresses)
			next = resp.Pagination.NextKey
		}
	})
	s.Run("Total", func() {
		resp, err := s.queryClient.DeployerFeeShares(s.Ctx, request(nil, 0, 0, true))
		s.Require().NoError(err)
		s.Require().Equal(len(contractAddressList), int(resp.Pagination.Total))
		s.Require().ElementsMatch(contractAddressList, resp.ContractAddresses)
	})
}

func (s *IntegrationTestSuite) TestWithdrawerFeeShares() {
	s.SetupTest()
	sender := s.TestAccs[0]
	withdrawer := s.TestAccs[1]

	var contractAddressList []string
	var index uint64
	for index < 5 {
		contractAddress := s.InstantiateContract(sender.String(), "")
		contractAddressList = append(contractAddressList, contractAddress)
		index++
	}

	// RegsisFeeShare
	for _, contractAddress := range contractAddressList {
		msg := &types.MsgRegisterFeeShare{
			ContractAddress:   contractAddress,
			DeployerAddress:   sender.String(),
			WithdrawerAddress: withdrawer.String(),
		}

		_, err := s.App.Keepers.FeeShareKeeper.RegisterFeeShare(s.Ctx, msg)
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
		for i := 0; i < len(contractAddressList); i += step {
			resp, err := s.queryClient.WithdrawerFeeShares(s.Ctx, request(nil, uint64(i), uint64(step), false))
			s.Require().NoError(err)
			s.Require().LessOrEqual(len(resp.ContractAddresses), step)
			s.Require().Subset(contractAddressList, resp.ContractAddresses)
		}
	})
	s.Run("ByKey", func() {
		step := 2
		var next []byte
		for i := 0; i < len(contractAddressList); i += step {
			resp, err := s.queryClient.WithdrawerFeeShares(s.Ctx, request(next, 0, uint64(step), false))
			s.Require().NoError(err)
			s.Require().LessOrEqual(len(resp.ContractAddresses), step)
			s.Require().Subset(contractAddressList, resp.ContractAddresses)
			next = resp.Pagination.NextKey
		}
	})
	s.Run("Total", func() {
		resp, err := s.queryClient.WithdrawerFeeShares(s.Ctx, request(nil, 0, 0, true))
		s.Require().NoError(err)
		s.Require().Equal(len(contractAddressList), int(resp.Pagination.Total))
		s.Require().ElementsMatch(contractAddressList, resp.ContractAddresses)
	})
}
