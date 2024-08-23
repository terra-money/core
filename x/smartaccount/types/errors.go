package types

import (
	errorsmod "cosmossdk.io/errors"
)

var (
	ErrFeeShareDisabled = errorsmod.Register(ModuleName, 1, "smartaccount module is disabled by governance")
)
