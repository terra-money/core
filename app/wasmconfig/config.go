package wasmconfig

import (
	"github.com/spf13/cast"
	"github.com/spf13/cobra"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
)

// config default values
const (
	DefaultContractQueryGasLimit      = uint64(3_000_000)
	DefaultContractSimulationGasLimit = uint64(50_000_000)
	DefaultContractDebugMode          = false
	DefaultContractMemoryCacheSize    = uint32(2048)
)

// Config is the extra config required for wasm
type Config struct {
	// SimulationGasLimit is the max gas to be used in a smart query contract call
	ContractQueryGasLimit uint64 `mapstructure:"contract-query-gas-limit"`

	// SimulationGasLimit is the max gas to be used in a tx simulation call.
	// When not set the consensus max block gas is used instead
	ContractSimulationGasLimit uint64 `mapstructure:"contract-query-gas-limit"`

	// ContractDebugMode log what contract print
	ContractDebugMode bool `mapstructure:"contract-debug-mode"`

	// MemoryCacheSize in MiB not bytes
	ContractMemoryCacheSize uint32 `mapstructure:"contract-memory-cache-size"`
}

// ToWasmConfig convert config to wasmd's config
func (c Config) ToWasmConfig() wasmtypes.WasmConfig {
	return wasmtypes.WasmConfig{
		SimulationGasLimit: &c.ContractSimulationGasLimit,
		SmartQueryGasLimit: c.ContractQueryGasLimit,
		MemoryCacheSize:    c.ContractMemoryCacheSize,
		ContractDebugMode:  c.ContractDebugMode,
	}
}

// DefaultConfig returns the default settings for WasmConfig
func DefaultConfig() *Config {
	return &Config{
		ContractQueryGasLimit:      DefaultContractQueryGasLimit,
		ContractSimulationGasLimit: DefaultContractSimulationGasLimit,
		ContractDebugMode:          DefaultContractDebugMode,
		ContractMemoryCacheSize:    DefaultContractMemoryCacheSize,
	}
}

// GetConfig load config values from the app options
func GetConfig(appOpts servertypes.AppOptions) *Config {
	return &Config{
		ContractQueryGasLimit:      cast.ToUint64(appOpts.Get("wasm.contract-query-gas-limit")),
		ContractSimulationGasLimit: cast.ToUint64(appOpts.Get("wasm.contract-simulation-gas-limit")),
		ContractDebugMode:          cast.ToBool(appOpts.Get("wasm.contract-debug-mode")),
		ContractMemoryCacheSize:    cast.ToUint32(appOpts.Get("wasm.contract-memory-cache-size")),
	}
}

const (
	flagContractQueryGasLimit      = "wasm.contract-query-gas-limit"
	flagContractSimulationGasLimit = "wasm.contract-simulation-gas-limit"
	flagContractDebugMode          = "wasm.contract-debug-mode"
	flagContractMemoryCacheSize    = "wasm.contract-memory-cache-size"
)

// AddConfigFlags implements servertypes.WasmConfigFlags interface.
func AddConfigFlags(startCmd *cobra.Command) {
	startCmd.Flags().Uint64(flagContractQueryGasLimit, DefaultContractQueryGasLimit, "Set the max gas that can be spent on executing a query with a Wasm contract")
	startCmd.Flags().Uint64(flagContractSimulationGasLimit, DefaultContractSimulationGasLimit, "Set the max gas that can be spent when executing a simulation TX")
	startCmd.Flags().Bool(flagContractDebugMode, DefaultContractDebugMode, "The flag to specify whether print contract logs or not")
	startCmd.Flags().Uint32(flagContractMemoryCacheSize, DefaultContractMemoryCacheSize, "Sets the size in MiB (NOT bytes) of an in-memory cache for Wasm modules. Set to 0 to disable.")
}

// DefaultConfigTemplate default config template for wasm module
const DefaultConfigTemplate = `
[wasm]
# The maximum gas amount can be spent for contract query.
# The contract query will invoke contract execution vm,
# so we need to restrict the max usage to prevent DoS attack
contract-query-gas-limit = "{{ .WASMConfig.ContractQueryGasLimit }}"

# The maximum gas amount can be used in a tx simulation call.
contract-simulation-gas-limit= "{{ .WASMConfig.ContractSimulationGasLimit }}"

# The flag to specify whether print contract logs or not
contract-debug-mode = "{{ .WASMConfig.ContractDebugMode }}"

# The WASM VM memory cache size in MiB not bytes
contract-memory-cache-size = "{{ .WASMConfig.ContractMemoryCacheSize }}"
`
