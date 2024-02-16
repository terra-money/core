package types

import (
	errorsmod "cosmossdk.io/errors"
)

var (
	ErrFeeShareDisabled              = errorsmod.Register(ModuleName, 1, "smartaccounts module is disabled by governance")
	ErrFeeShareAlreadyRegistered     = errorsmod.Register(ModuleName, 2, "smartaccounts already exists for given contract")
	ErrFeeShareNoContractDeployed    = errorsmod.Register(ModuleName, 3, "no contract deployed")
	ErrFeeShareContractNotRegistered = errorsmod.Register(ModuleName, 4, "no smartaccounts registered for contract")
	ErrFeeSharePayment               = errorsmod.Register(ModuleName, 5, "smartaccounts payment error")
	ErrFeeShareInvalidWithdrawer     = errorsmod.Register(ModuleName, 6, "invalid withdrawer address")
)
