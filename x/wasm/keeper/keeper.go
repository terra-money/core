package keeper

import (
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/CosmWasm/wasmd/x/wasm/types"
	keepertypes "github.com/terra-money/core/v2/x/wasm/types"
)

var _ keepertypes.KeeperInterface = Keeper{}

type Keeper struct {
	*wasmkeeper.Keeper
	storeKey storetypes.StoreKey
	cdc      codec.Codec
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
		Keeper:   &keeper,
		storeKey: storeKey,
		cdc:      cdc,
	}
}

// After executing the contract, get all executed
// contract addresses from the store, if there is
// a store already then check if the contract address
// exists in the list, if not then update the store,
// If the contract does not exist in the store, add it.
func (k Keeper) AfterExecuteContract(ctx sdk.Context, msg *types.MsgExecuteContract, res *types.MsgExecuteContractResponse) error {
	contracts, found := k.GetExecutedContractAddresses(ctx)

	if found {
		for _, contract := range contracts.ContractAddresses {
			if contract == msg.Contract {
				return nil
			}
		}
	}

	contracts.ContractAddresses = append(contracts.ContractAddresses, msg.Contract)

	err := k.SetExecutedContractAddresses(ctx, contracts)
	if err != nil {
		return err
	}
	return nil
}
