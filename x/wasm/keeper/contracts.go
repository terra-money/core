package keeper

import (
	"github.com/terra-money/core/v2/x/wasm/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetExecutedContractAddresses(ctx sdk.Context) (contracts types.ExecutedContracts, found bool) {
	store := ctx.KVStore(k.storeKey)
	contractAddressesKey := types.GetExecutedContractsKey()
	b := store.Get(contractAddressesKey)
	if b == nil {
		return contracts, false
	}

	k.cdc.MustUnmarshal(b, &contracts)
	return contracts, true
}

func (k Keeper) SetExecutedContractAddresses(ctx sdk.Context, contracts types.ExecutedContracts) error {
	store := ctx.KVStore(k.storeKey)
	contractAddressesKey := types.GetExecutedContractsKey()
	b := k.cdc.MustMarshal(&contracts)
	store.Set(contractAddressesKey, b)
	return nil
}

func (k Keeper) DeleteExecutedContractAddresses(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	contractAddressesKey := types.GetExecutedContractsKey()
	store.Delete(contractAddressesKey)
}
