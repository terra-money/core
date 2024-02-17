package types

import (
	errorsmod "cosmossdk.io/errors"
)

var (
	ErrFeeShareDisabled = errorsmod.Register(ModuleName, 1, "smartaccounts module is disabled by governance")
)
