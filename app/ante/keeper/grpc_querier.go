package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/terra-money/core/v2/app/ante/types"
)

// Querier is used as Keeper will have duplicate methods if used directly, and gRPC names take precedence over keeper
type Querier struct {
	Keeper
}

var _ types.QueryServer = Querier{}

// Params queries the ante parameters
func (k Querier) Params(c context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	params := k.GetParams(ctx)

	return &types.QueryParamsResponse{Params: params}, nil
}

// MinimumCommission queries the minimum commission
func (k Querier) MinimumCommission(c context.Context, _ *types.QueryMinimumCommissionRequest) (*types.QueryMinimumCommissionResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	minimumCommission := k.GetMinimumCommission(ctx)

	return &types.QueryMinimumCommissionResponse{MinimumCommission: minimumCommission}, nil
}
