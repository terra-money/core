package types

import (
	fmt "fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Ante params default values
var (
	// Default minimum commission value
	DefaultMinimumCommission sdk.Dec = sdk.NewDecWithPrec(10, 2) // 10%
)

// NewGenesisState return new GenesisState instance
func NewGenesisState(params Params, minimumCommission sdk.Dec) *GenesisState {
	return &GenesisState{
		Params:            params,
		MinimumCommission: minimumCommission,
	}
}

// DefaultGenesisState return default GenesisState
func DefaultGenesisState() *GenesisState {
	return NewGenesisState(DefaultParams(), DefaultMinimumCommission)
}

// Validate performs basic validation of ante genesis data returning an
// error for any failed validation criteria.
func (genState GenesisState) Validate() error {

	if err := genState.Params.Validate(); err != nil {
		return err
	}

	if err := validateMinimumCommission(genState.MinimumCommission); err != nil {
		return err
	}

	return nil
}

func validateMinimumCommission(minimumCommission sdk.Dec) error {
	if minimumCommission.GT(sdk.OneDec()) || minimumCommission.IsNegative() {
		return fmt.Errorf("minimum commission must be [0, 1]: %d", minimumCommission)
	}

	return nil
}
