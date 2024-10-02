package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/terra-money/core/v2/x/feeburn/types"
)

var _ types.QueryServer = Querier{}

// Querier defines a wrapper around the x/feeburn
// keeper providing gRPC method handlers.
type Querier struct {
	Keeper
}

func NewQuerier(k Keeper) Querier {
	return Querier{Keeper: k}
}

// Params returns the fees module params
func (q Querier) Params(c context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	params := q.GetParams(ctx)
	return &types.QueryParamsResponse{Params: params}, nil
}
