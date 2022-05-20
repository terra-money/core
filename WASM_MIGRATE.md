## Wasm Migration Guide from Terra Classic
Terra Rebirth is now using wasm module of [wasmd](https://github.com/CosmWasm/wasmd) and it introduces minor compatibility issue with Terra Classic.

### Instantiate 

The contracts are using reply to check instantiated contract address, 
should update the proto file to the following.
```protobuf
// MsgInstantiateContractResponse return instantiation result data
message MsgInstantiateContractResponse {
  // Address is the bech32 address of the new contract instance.
  string address = 1;
  // Data contains base64-encoded bytes to returned from the contract
  bytes data = 2;
}
```

Event key for instantiated contract also should be changed
from `instantiate`.`contract_address` to `instantiate`.`_contract_address`.

