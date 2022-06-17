package types

// NewGenesisState return new GenesisState instance
func NewGenesisState(params Params) *GenesisState {
	return &GenesisState{
		Params: params,
	}
}

// DefaultGenesisState return default GenesisState
func DefaultGenesisState() *GenesisState {
	return NewGenesisState(DefaultParams())
}

// Validate performs basic validation of ante genesis data returning an
// error for any failed validation criteria.
func (genState GenesisState) Validate() error {
	return genState.Params.Validate()
}
