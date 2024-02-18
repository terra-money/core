package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/terra-money/core/v2/x/smartaccount/types"
)

type MsgServer struct {
	k Keeper
}

// NewMsgServer returns the MsgServer implementation.
func NewMsgServer(k Keeper) types.MsgServer {
	return &MsgServer{k}
}

func (ms MsgServer) CreateSmartAccount(goCtx context.Context, msg *types.MsgCreateSmartAccount) (*types.MsgCreateSmartAccountResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := ms.k.SetSetting(ctx, msg.Account, types.Setting{
		Owner: msg.Account,
	}); err != nil {
		return nil, err
	}
	return &types.MsgCreateSmartAccountResponse{}, nil
}

func (ms MsgServer) Authorization(goCtx context.Context, msg *types.MsgAuthorization) (*types.MsgAuthorizationResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	_ = ctx
	return &types.MsgAuthorizationResponse{}, nil
}

func (ms MsgServer) UpdateAuthorization(goCtx context.Context, msg *types.MsgUpdateAuthorization) (*types.MsgUpdateAuthorizationResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	_ = ctx
	return &types.MsgUpdateAuthorizationResponse{}, nil
}

func (ms MsgServer) UpdateTransactionHooks(goCtx context.Context, msg *types.MsgUpdateTransactionHooks) (*types.MsgUpdateTransactionHooksResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	_ = ctx
	return &types.MsgUpdateTransactionHooksResponse{}, nil
}

func (ms MsgServer) DisableSmartAccount(goCtx context.Context, msg *types.MsgDisableSmartAccount) (*types.MsgDisableSmartAccountResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	_ = ctx
	return &types.MsgDisableSmartAccountResponse{}, nil
}
