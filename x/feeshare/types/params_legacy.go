package types

// TODO: Remove this and params_legacy_test.go after v0.47.x (v16) upgrade

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Parameter store key
var (
	DefaultEnableFeeShare  = true
	DefaultDeveloperShares = sdk.NewDecWithPrec(50, 2) // 50%
	DefaultAllowedDenoms   = []string(nil)             // all allowed

	ParamStoreKeyEnableFeeShare  = []byte("EnableFeeShare")
	ParamStoreKeyDeveloperShares = []byte("DeveloperShares")
	ParamStoreKeyAllowedDenoms   = []byte("AllowedDenoms")
)

// ParamKeyTable returns the parameter key table.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// ParamSetPairs returns the parameter set pairs.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(ParamStoreKeyEnableFeeShare, &p.EnableFeeShare, validateBool),
		paramtypes.NewParamSetPair(ParamStoreKeyDeveloperShares, &p.DeveloperShares, validateShares),
		paramtypes.NewParamSetPair(ParamStoreKeyAllowedDenoms, &p.AllowedDenoms, validateArray),
	}
}
