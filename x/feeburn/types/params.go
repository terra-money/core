package types

import (
	"fmt"
)

// NewParams creates a new Params object
func NewParams(enableFeeBurn bool) Params {
	return Params{
		EnableFeeBurn: enableFeeBurn,
	}
}

func DefaultParams() Params {
	return Params{
		EnableFeeBurn: true,
	}
}
func (p Params) Validate() error {
	if err := validateBool(p.EnableFeeBurn); err != nil {
		return err
	}
	return nil
}

func validateBool(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}
