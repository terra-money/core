package keeper

import (
	"context"
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/terra-money/core/v2/x/smartaccount/types"
)

type MsgServer struct {
	k Keeper
}

// NewMsgServer returns the MsgServer implementation.
func NewMsgServer(k Keeper) types.MsgServer {
	return &MsgServer{k}
}

func (ms MsgServer) CreateSmartAccount(
	goCtx context.Context, msg *types.MsgCreateSmartAccount,
) (*types.MsgCreateSmartAccountResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	setting, _ := ms.k.GetSetting(ctx, msg.Account)
	if setting != nil {
		return nil, sdkerrors.ErrInvalidRequest.Wrapf("smart account already exists for %s", msg.Account)
	}

	if err := ms.k.SetSetting(ctx, types.Setting{
		Owner:    msg.Account,
		Fallback: true,
	}); err != nil {
		return nil, err
	}
	return &types.MsgCreateSmartAccountResponse{}, nil
}

func (ms MsgServer) UpdateAuthorization(
	goCtx context.Context, msg *types.MsgUpdateAuthorization,
) (*types.MsgUpdateAuthorizationResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	setting, err := ms.k.GetSetting(ctx, msg.Account)
	if sdkerrors.ErrKeyNotFound.Is(err) {
		setting = &types.Setting{
			Owner: msg.Account,
		}
	} else if err != nil {
		return nil, err
	}
	setting.Authorization = msg.AuthorizationMsgs
	// TODO: check if this is right
	for _, auth := range msg.AuthorizationMsgs {
		sudoInitMsg := types.SudoMsg{Initialization: auth.InitMsg}
		sudoInitMsgBs, err := json.Marshal(sudoInitMsg)
		if err != nil {
			return nil, err
		}

		contractAddress := sdk.AccAddress(auth.ContractAddress)
		ms.k.wasmKeeper.Sudo(ctx, contractAddress, sudoInitMsgBs)
	}
	setting.Fallback = msg.Fallback
	if err := ms.k.SetSetting(ctx, *setting); err != nil {
		return nil, err
	}
	return &types.MsgUpdateAuthorizationResponse{}, nil
}

func (ms MsgServer) UpdateTransactionHooks(
	goCtx context.Context, msg *types.MsgUpdateTransactionHooks,
) (*types.MsgUpdateTransactionHooksResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	setting, err := ms.k.GetSetting(ctx, msg.Account)
	if sdkerrors.ErrKeyNotFound.Is(err) {
		setting = &types.Setting{
			Owner: msg.Account,
		}
	} else if err != nil {
		return nil, err
	}
	setting.PostTransaction = msg.PostTransactionHooks
	setting.PreTransaction = msg.PreTransactionHooks
	if err := ms.k.SetSetting(ctx, *setting); err != nil {
		return nil, err
	}
	return &types.MsgUpdateTransactionHooksResponse{}, nil
}

// DisableSmartAccount converts smart acc back to a basic acc
func (ms MsgServer) DisableSmartAccount(
	goCtx context.Context, msg *types.MsgDisableSmartAccount,
) (*types.MsgDisableSmartAccountResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := ms.k.DeleteSetting(ctx, msg.Account); err != nil {
		return nil, err
	}
	return &types.MsgDisableSmartAccountResponse{}, nil
}
