package types

import (
	"github.com/cometbft/cometbft/libs/log"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/CosmWasm/wasmd/x/wasm/types"
)

type KeeperInterface interface {
	GetParams(ctx sdk.Context) types.Params
	SetParams(ctx sdk.Context, ps types.Params) error
	GetAuthority() string
	GetGasRegister() types.GasRegister
	Sudo(ctx sdk.Context, contractAddress sdk.AccAddress, msg []byte) ([]byte, error)
	IterateContractsByCreator(ctx sdk.Context, creator sdk.AccAddress, cb func(address sdk.AccAddress) bool)
	IterateContractsByCode(ctx sdk.Context, codeID uint64, cb func(address sdk.AccAddress) bool)
	GetContractHistory(ctx sdk.Context, contractAddr sdk.AccAddress) []types.ContractCodeHistoryEntry
	QuerySmart(ctx sdk.Context, contractAddr sdk.AccAddress, req []byte) ([]byte, error)
	QueryRaw(ctx sdk.Context, contractAddress sdk.AccAddress, key []byte) []byte
	GetContractInfo(ctx sdk.Context, contractAddress sdk.AccAddress) *types.ContractInfo
	HasContractInfo(ctx sdk.Context, contractAddress sdk.AccAddress) bool
	IterateContractInfo(ctx sdk.Context, cb func(sdk.AccAddress, types.ContractInfo) bool)
	IterateContractState(ctx sdk.Context, contractAddress sdk.AccAddress, cb func(key, value []byte) bool)
	GetCodeInfo(ctx sdk.Context, codeID uint64) *types.CodeInfo
	IterateCodeInfos(ctx sdk.Context, cb func(uint64, types.CodeInfo) bool)
	GetByteCode(ctx sdk.Context, codeID uint64) ([]byte, error)
	IsPinnedCode(ctx sdk.Context, codeID uint64) bool
	InitializePinnedCodes(ctx sdk.Context) error
	PeekAutoIncrementID(ctx sdk.Context, sequenceKey []byte) uint64
	Logger(ctx sdk.Context) log.Logger
	QueryGasLimit() sdk.Gas
}
