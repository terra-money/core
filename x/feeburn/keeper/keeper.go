package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/cometbft/cometbft/libs/log"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/terra-money/core/v2/x/tokenfactory/types"
)

type Keeper struct {
	storeKey  storetypes.StoreKey
	cdc       codec.BinaryCodec
	authority string
}

// NewKeeper returns a new instance of the x/tokenfactory keeper
func NewKeeper(storeKey storetypes.StoreKey, cdc codec.BinaryCodec, authority string) Keeper {

	return Keeper{
		storeKey:  storeKey,
		cdc:       cdc,
		authority: authority,
	}
}

// Logger returns a logger for the x/tokenfactory module
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) GetAuthority() string {
	return k.authority
}
