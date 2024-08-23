package types

// NewGenesisState creates a new genesis state.
func NewGenesisState(params Params, settings []*Setting) GenesisState {
	return GenesisState{
		Params:   params,
		Settings: settings,
	}
}

// DefaultGenesisState sets default evm genesis state with empty accounts and
// default params and chain config values.
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params:   DefaultParams(),
		Settings: DefaultSettings(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	if err := gs.Params.Validate(); err != nil {
		return err
	}
	for _, setting := range gs.Settings {
		if err := setting.Validate(); err != nil {
			return err
		}
	}
	return nil
}
