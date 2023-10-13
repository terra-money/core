package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewParams creates a new Params object
func NewParams(
	enableFeeShare bool,
	developerShares sdk.Dec,
	allowedDenoms []string,
) Params {
	return Params{
		EnableFeeShare:  enableFeeShare,
		DeveloperShares: developerShares,
		AllowedDenoms:   allowedDenoms,
	}
}

func DefaultParams() Params {
	return Params{
		EnableFeeShare:  DefaultEnableFeeShare,
		DeveloperShares: DefaultDeveloperShares,
		AllowedDenoms:   DefaultAllowedDenoms,
	}
}

func validateBool(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validateShares(i interface{}) error {
	v, ok := i.(sdk.Dec)

	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("invalid parameter: nil")
	}

	if v.IsNegative() {
		return fmt.Errorf("value cannot be negative: %T", i)
	}

	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("value cannot be greater than 1: %T", i)
	}

	return nil
}

func validateArray(i interface{}) error {
	_, ok := i.([]string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	for _, denom := range i.([]string) {
		if denom == "" {
			return fmt.Errorf("denom cannot be blank")
		}
	}

	return nil
}

func (p Params) Validate() error {
	if err := validateBool(p.EnableFeeShare); err != nil {
		return err
	}
	if err := validateShares(p.DeveloperShares); err != nil {
		return err
	}
	err := validateArray(p.AllowedDenoms)
	return err
}
