package keeper

import (
	"fmt"
	"time"

	sdkerrors "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distributionkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	"github.com/cosmos/cosmos-sdk/x/distribution/types"
)

type Keeper struct {
	distributionkeeper.Keeper
	bankKeeper                types.BankKeeper
	communityPoolBlockedUntil time.Time
}

func NewKeeper(
	cdc codec.BinaryCodec,
	key storetypes.StoreKey,
	authKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	stakingKeeper types.StakingKeeper,
	feeCollectorName string,
	authority string,
) Keeper {

	return Keeper{
		Keeper: distributionkeeper.NewKeeper(cdc,
			key,
			authKeeper,
			bankKeeper,
			stakingKeeper,
			feeCollectorName,
			authority,
		),
		bankKeeper:                bankKeeper,
		communityPoolBlockedUntil: time.Date(2025, 1, 1, 0, 0, 0, 0, time.Local),
	}
}

// DistributeFromFeePool distributes funds from the distribution module account to
// a receiver address while updating the community pool
func (k Keeper) DistributeFromFeePool(ctx sdk.Context, amount sdk.Coins, receiveAddr sdk.AccAddress) error {
	if k.communityPoolBlockedUntil.After(ctx.BlockHeader().Time) {
		message := fmt.Sprintf("CommunityPool is blocked until %s", k.communityPoolBlockedUntil)
		return sdkerrors.New(types.ModuleName, 999, message)
	}

	return k.Keeper.DistributeFromFeePool(ctx, amount, receiveAddr)
}
