package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Event Hooks
// These can be utilized to communicate between a bank keeper and another
// keeper which must take particular actions when sends happen.
// The second keeper must implement this interface, which then the
// bank keeper can call.

// BankHooks event hooks for bank sends
type BankHooks interface {
	TrackBeforeSend(ctx sdk.Context, from, to sdk.AccAddress, amount sdk.Coins)       // Must be before any send is executed
	BlockBeforeSend(ctx sdk.Context, from, to sdk.AccAddress, amount sdk.Coins) error // Must be before any send is executed
}
