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

var (
	communityPoolBlockedUntil = time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	ErrCommunityPoolBlocked   = sdkerrors.New(
		types.ModuleName,
		9999,
		fmt.Sprintf("CommunityPool is blocked until %s", communityPoolBlockedUntil),
	)
)

type Keeper struct {
	distributionkeeper.Keeper
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
	}
}

// DistributeFromFeePool distributes funds from the distribution module account to
// a receiver address while updating the community pool
func (k Keeper) DistributeFromFeePool(ctx sdk.Context, amount sdk.Coins, receiveAddr sdk.AccAddress) error {
	if communityPoolBlockedUntil.After(ctx.BlockHeader().Time) {
		return ErrCommunityPoolBlocked
	}

	return k.Keeper.DistributeFromFeePool(ctx, amount, receiveAddr)
}
