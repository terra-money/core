package keeper

import (
	"slices"

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

func (k Keeper) AfterExecuteContract(ctx sdk.Context, msg *types.MsgExecuteContract, res *types.MsgExecuteContractResponse) error {
	events := ctx.EventManager().Events()
	contractAddresses := []string{}

	for _, ev := range events {
		if ev.Type != "execute" {
			continue
		}

		for _, attr := range ev.Attributes {
			if attr.Key != "_contract_address" {
				continue
			}
			// if the contract address has already been
			// added just skip it to avoid duplicates
			if slices.Contains(contractAddresses, attr.Value) {
				continue
			}

			contractAddresses = append(contractAddresses, attr.Value)
		}
	}

	if len(contractAddresses) != 0 {
		executedContracts := keepertypes.ExecutedContracts{
			ContractAddresses: contractAddresses,
		}

		err := k.SetExecutedContractAddresses(ctx, executedContracts)
		if err != nil {
			return err
		}
	}

	return nil
}
