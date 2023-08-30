package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	accountkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	custombankkeeper "github.com/terra-money/alliance/custom/bank/keeper"
	customterratypes "github.com/terra-money/core/v2/custom/bank/types"
)

type Keeper struct {
	custombankkeeper.Keeper
	hooks customterratypes.BankHooks
}

var _ bankkeeper.Keeper = Keeper{}

func NewBaseKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	ak accountkeeper.AccountKeeper,
	blockedAddrs map[string]bool,
	authority string,
) Keeper {
	keeper := Keeper{
		Keeper: custombankkeeper.NewBaseKeeper(cdc, storeKey, ak, blockedAddrs, authority),
		hooks:  nil,
	}

	return keeper
}
