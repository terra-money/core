package app

import (
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
)

const (
	// DefaultTerraInstanceCost is initially set the same as in wasmd
	DefaultTerraInstanceCost uint64 = 60_000
	// DefaultTerraCompileCost set to a large number for testing
	DefaultTerraCompileCost uint64 = 100
)

// TerraGasRegisterConfig is defaults plus a custom compile amount
func TerraGasRegisterConfig() wasmkeeper.WasmGasRegisterConfig {
	gasConfig := wasmkeeper.DefaultGasRegisterConfig()
	gasConfig.InstanceCost = DefaultTerraInstanceCost
	gasConfig.CompileCost = DefaultTerraCompileCost

	return gasConfig
}

// NewTerraWasmGasRegister return gas register for wasm module
func NewTerraWasmGasRegister() wasmkeeper.WasmGasRegister {
	return wasmkeeper.NewWasmGasRegister(TerraGasRegisterConfig())
}
