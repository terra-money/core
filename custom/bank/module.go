package bank

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank/exported"
	"github.com/cosmos/cosmos-sdk/x/bank/types"

	customalliancemod "github.com/terra-money/alliance/custom/bank"
	custombankkeeper "github.com/terra-money/alliance/custom/bank/keeper"
	customtypes "github.com/terra-money/core/v2/custom/bank/types"
)

// AppModule wraps around the bank module and the bank keeper to return the right total supply ignoring bonded tokens
// that the alliance module minted to rebalance the voting power
// It modifies the TotalSupply and SupplyOf GRPC queries
type AppModule struct {
	customalliancemod.AppModule
	hooks customtypes.BankHooks
}

// NewAppModule creates a new AppModule object
func NewAppModule(cdc codec.Codec, keeper custombankkeeper.Keeper, accountKeeper types.AccountKeeper, ss exported.Subspace) AppModule {
	mod := customalliancemod.NewAppModule(cdc, keeper, accountKeeper, ss)

	return AppModule{
		AppModule: mod,
		hooks:     nil,
	}
}

// Set the bank hooks
func (am *AppModule) SetHooks(bh customtypes.BankHooks) *AppModule {
	if am.hooks != nil {
		panic("cannot set bank hooks twice")
	}

	am.hooks = bh

	return am
}

// SendCoins transfers amt coins from a sending account to a receiving account.
// An error is returned upon failure.
func (am AppModule) SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error {
	// BlockBeforeSend hook should always be called before the TrackBeforeSend hook.
	err := am.hooks.BlockBeforeSend(ctx, fromAddr, toAddr, amt)
	if err != nil {
		return err
	}
	am.hooks.TrackBeforeSend(ctx, fromAddr, toAddr, amt)

	return am.SendCoins(ctx, fromAddr, toAddr, amt)
}
