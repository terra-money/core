package keeper

import (
	"fmt"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/CosmWasm/wasmd/x/wasm/types"
	keepertypes "github.com/terra-money/core/v2/custom/wasmd/types"
)

var _ keepertypes.KeeperInterface = Keeper{}

type Keeper struct {
	*wasmkeeper.Keeper
}

func NewKeeper(
	cdc codec.Codec,
	storeKey storetypes.StoreKey,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	stakingKeeper types.StakingKeeper,
	distrKeeper types.DistributionKeeper,
	ics4Wrapper types.ICS4Wrapper,
	channelKeeper types.ChannelKeeper,
	portKeeper types.PortKeeper,
	capabilityKeeper types.CapabilityKeeper,
	portSource types.ICS20TransferPortSource,
	router wasmkeeper.MessageRouter,
	grpcQueryRouter wasmkeeper.GRPCQueryRouter,
	homeDir string,
	wasmConfig types.WasmConfig,
	availableCapabilities string,
	authority string,
	opts ...wasmkeeper.Option,
) Keeper {
	keeper := wasmkeeper.NewKeeper(
		cdc,
		storeKey,
		accountKeeper,
		bankKeeper,
		stakingKeeper,
		distrKeeper,
		ics4Wrapper,
		channelKeeper,
		portKeeper,
		capabilityKeeper,
		portSource,
		router,
		grpcQueryRouter,
		homeDir,
		wasmConfig,
		availableCapabilities,
		authority,
		opts...,
	)

	return Keeper{
		&keeper,
	}
}

func (k Keeper) AfterExecuteContract(ctx sdk.Context, msg *types.MsgExecuteContract, res *types.MsgExecuteContractResponse) error {
	fmt.Print("\n\n\n\n\n\n\n\nAfterExecuteContract\n")
	fmt.Print(ctx.EventManager().Events())
	fmt.Print("\n\n\n\n\n\n\n")
	return nil
}
