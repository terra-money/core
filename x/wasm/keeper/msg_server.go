package keeper

import (
	"context"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/CosmWasm/wasmd/x/wasm/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ types.MsgServer = customMsgServer{}

// grpc message server implementation
type customMsgServer struct {
	msgServer types.MsgServer
	keeper    *Keeper
}

// NewCustomMsgServerImpl default constructor
func NewCustomMsgServerImpl(k *Keeper) types.MsgServer {
	msgServer := wasmkeeper.NewMsgServerImpl(k.Keeper)

	return customMsgServer{
		msgServer: msgServer,
		keeper:    k,
	}
}

// EexcuteContract wraps the original call but it collects
// the addresses of contracts involved in the transaction
func (m customMsgServer) ExecuteContract(goCtx context.Context, msg *types.MsgExecuteContract) (*types.MsgExecuteContractResponse, error) {
	res, err := m.msgServer.ExecuteContract(goCtx, msg)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	err = m.keeper.AfterExecuteContract(ctx, msg, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (m customMsgServer) StoreCode(goCtx context.Context, msg *types.MsgStoreCode) (*types.MsgStoreCodeResponse, error) {
	return m.msgServer.StoreCode(goCtx, msg)
}
func (m customMsgServer) InstantiateContract(goCtx context.Context, msg *types.MsgInstantiateContract) (*types.MsgInstantiateContractResponse, error) {
	return m.msgServer.InstantiateContract(goCtx, msg)
}
func (m customMsgServer) InstantiateContract2(goCtx context.Context, msg *types.MsgInstantiateContract2) (*types.MsgInstantiateContract2Response, error) {
	return m.msgServer.InstantiateContract2(goCtx, msg)
}
func (m customMsgServer) MigrateContract(goCtx context.Context, msg *types.MsgMigrateContract) (*types.MsgMigrateContractResponse, error) {
	return m.msgServer.MigrateContract(goCtx, msg)
}
func (m customMsgServer) UpdateAdmin(goCtx context.Context, msg *types.MsgUpdateAdmin) (*types.MsgUpdateAdminResponse, error) {
	return m.msgServer.UpdateAdmin(goCtx, msg)
}
func (m customMsgServer) ClearAdmin(goCtx context.Context, msg *types.MsgClearAdmin) (*types.MsgClearAdminResponse, error) {
	return m.msgServer.ClearAdmin(goCtx, msg)
}
func (m customMsgServer) UpdateInstantiateConfig(goCtx context.Context, msg *types.MsgUpdateInstantiateConfig) (*types.MsgUpdateInstantiateConfigResponse, error) {
	return m.msgServer.UpdateInstantiateConfig(goCtx, msg)
}
func (m customMsgServer) UpdateParams(goCtx context.Context, req *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	return m.msgServer.UpdateParams(goCtx, req)
}
func (m customMsgServer) PinCodes(goCtx context.Context, req *types.MsgPinCodes) (*types.MsgPinCodesResponse, error) {
	return m.msgServer.PinCodes(goCtx, req)
}
func (m customMsgServer) UnpinCodes(goCtx context.Context, req *types.MsgUnpinCodes) (*types.MsgUnpinCodesResponse, error) {
	return m.msgServer.UnpinCodes(goCtx, req)
}
func (m customMsgServer) SudoContract(goCtx context.Context, req *types.MsgSudoContract) (*types.MsgSudoContractResponse, error) {
	return m.msgServer.SudoContract(goCtx, req)
}
func (m customMsgServer) StoreAndInstantiateContract(goCtx context.Context, req *types.MsgStoreAndInstantiateContract) (*types.MsgStoreAndInstantiateContractResponse, error) {
	return m.msgServer.StoreAndInstantiateContract(goCtx, req)
}
func (m customMsgServer) AddCodeUploadParamsAddresses(goCtx context.Context, req *types.MsgAddCodeUploadParamsAddresses) (*types.MsgAddCodeUploadParamsAddressesResponse, error) {
	return m.msgServer.AddCodeUploadParamsAddresses(goCtx, req)
}
func (m customMsgServer) RemoveCodeUploadParamsAddresses(goCtx context.Context, req *types.MsgRemoveCodeUploadParamsAddresses) (*types.MsgRemoveCodeUploadParamsAddressesResponse, error) {
	return m.msgServer.RemoveCodeUploadParamsAddresses(goCtx, req)
}
func (m customMsgServer) StoreAndMigrateContract(goCtx context.Context, req *types.MsgStoreAndMigrateContract) (*types.MsgStoreAndMigrateContractResponse, error) {
	return m.msgServer.StoreAndMigrateContract(goCtx, req)
}
func (m customMsgServer) UpdateContractLabel(goCtx context.Context, msg *types.MsgUpdateContractLabel) (*types.MsgUpdateContractLabelResponse, error) {
	return m.msgServer.UpdateContractLabel(goCtx, msg)
}
