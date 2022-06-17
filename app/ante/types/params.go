package types

import (
	fmt "fmt"

	yaml "gopkg.in/yaml.v2"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Ante params default values
const (
	// DefaultMinimumCommissionEnforced minimum commission Enforced flag
	DefaultMinimumCommissionEnforced bool = true
)

// Ante params default values
var (
	// Default maximum number of bonded validators
	DefaultMinimumCommission sdk.Dec = sdk.NewDecWithPrec(5, 2) // 5%
)

// Parameter keys
var (
	ParamStoreKeyMinimumCommissionEnforced = []byte("MinimumCommissionEnforced")
	ParamStoreKeyMinimumCommission         = []byte("MinimumCommission")
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable for ante
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams return new ante Params instance
func NewParams(minimumCommissionEnforced bool, minimumCommission sdk.Dec) Params {
	return Params{
		MinimumCommissionEnforced: minimumCommissionEnforced,
		MinimumCommission:         minimumCommission,
	}
}

// DefaultParams return default ante Params
func DefaultParams() Params {
	return NewParams(
		DefaultMinimumCommissionEnforced,
		DefaultMinimumCommission,
	)
}

// ParamSetPairs implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(ParamStoreKeyMinimumCommissionEnforced, &p.MinimumCommissionEnforced, validateMiniumCommissionEnforced),
		paramtypes.NewParamSetPair(ParamStoreKeyMinimumCommission, &p.MinimumCommission, validateMinimumCommission),
	}
}

// String returns a human readable string representation of the parameters.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// Validate validate a set of params
func (p Params) Validate() error {
	if err := validateMiniumCommissionEnforced(p.MinimumCommissionEnforced); err != nil {
		return err
	}

	if err := validateMinimumCommission(p.MinimumCommission); err != nil {
		return err
	}

	return nil
}

func validateMiniumCommissionEnforced(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validateMinimumCommission(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.GT(sdk.OneDec()) || v.IsNegative() {
		return fmt.Errorf("minimum commission must be [0, 1]: %d", v)
	}

	return nil
}
