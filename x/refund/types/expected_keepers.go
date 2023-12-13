package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Define the expected keeper interface
type BankKeeper interface {
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
}
