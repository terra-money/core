package types

import (
	fmt "fmt"

	yaml "gopkg.in/yaml.v2"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Ante params default values
const (
	// DefaultMinimumCommissionEnforced minimum commission Enforced flag
	DefaultMinimumCommissionEnforced bool = true
)

// Parameter keys
var (
	ParamStoreKeyMinimumCommissionEnforced = []byte("MinimumCommissionEnforced")
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable for ante
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams return new ante Params instance
func NewParams(minimumCommissionEnforced bool) Params {
	return Params{
		MinimumCommissionEnforced: minimumCommissionEnforced,
	}
}

// DefaultParams return default ante Params
func DefaultParams() Params {
	return NewParams(
		DefaultMinimumCommissionEnforced,
	)
}

// ParamSetPairs implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(ParamStoreKeyMinimumCommissionEnforced, &p.MinimumCommissionEnforced, validateMiniumCommissionEnforced),
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

	return nil
}

func validateMiniumCommissionEnforced(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}
