package types

import (
	sdk "github.com/cosmos/cosmos-sdk/codec/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

// RegisterInterfaces registers the sdk.Tx interface.
func RegisterInterfaces(registry sdk.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*govtypes.Content)(nil),
		&MinimumCommissionUpdateProposal{},
	)
}
