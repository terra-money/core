package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/terra-money/core/v2/x/smartaccount/types"
)

var _ types.QueryServer = Querier{}

// Querier defines a wrapper around the x/SmartAccounts keeper providing gRPC method
// handlers.
type Querier struct {
	Keeper
}

func NewQuerier(k Keeper) Querier {
	return Querier{Keeper: k}
}

// Params returns the fees module params
func (q Querier) Params(
	c context.Context,
	_ *types.QueryParamsRequest,
) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	params := q.GetParams(ctx)
	return &types.QueryParamsResponse{Params: params}, nil
}

// Setting returns the fees module setting
func (q Querier) Setting(
	c context.Context,
	req *types.QuerySettingRequest,
) (*types.QuerySettingResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	setting, err := q.GetSetting(ctx, req.Address)
	if err != nil {
		return nil, err
	}
	return &types.QuerySettingResponse{Setting: *setting}, nil
}
