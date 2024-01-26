package app

import (
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	"github.com/cosmos/cosmos-sdk/std"

	"github.com/terra-money/core/v2/app/params"
	pobtypes "github.com/terra-money/core/v2/x/builder/types"
)

var legacyCodecRegistered = false

// MakeEncodingConfig creates an EncodingConfig for testing
func MakeEncodingConfig() params.EncodingConfig {
	encodingConfig := params.MakeEncodingConfig()
	std.RegisterLegacyAminoCodec(encodingConfig.Amino)
	std.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	ModuleBasics.RegisterLegacyAminoCodec(encodingConfig.Amino)
	ModuleBasics.RegisterInterfaces(encodingConfig.InterfaceRegistry)

	// Register POB interfaces and concrete types
	// for the POB module that was used back in the day.
	pobtypes.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	pobtypes.RegisterLegacyAminoCodec(encodingConfig.Amino)

	if !legacyCodecRegistered {
		// authz module use this codec to get signbytes.
		// authz MsgExec can execute all message types,
		// so legacy.Cdc need to register all amino messages to get proper signature
		ModuleBasics.RegisterLegacyAminoCodec(legacy.Cdc)
		legacyCodecRegistered = true
	}

	return encodingConfig
}
