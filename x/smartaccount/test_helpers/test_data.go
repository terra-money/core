package test_helpers

import _ "embed"

//go:embed test_data/limit_send_only_hooks.wasm
var LimitSendOnlyHookWasm []byte

//go:embed test_data/smart_auth_contract.wasm
var SmartAuthContractWasm []byte
