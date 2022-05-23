## Wasm Migration Guide from Terra Classic
Terra Rebirth is now using wasm module of [wasmd](https://github.com/CosmWasm/wasmd) and it introduces minor compatibility issue with Terra Classic.

### Contract Address
Contract Address legnth will be different from normal account.
```go
// VerifyAddressLen ensures that the address matches the expected length
// ContractAddrLen = 32
// SDKAddrLen = 20
func VerifyAddressLen() func(addr []byte) error {
	return func(addr []byte) error {
		if len(addr) != ContractAddrLen && len(addr) != SDKAddrLen {
			return sdkerrors.ErrInvalidAddress
		}
		return nil
	}
}
```

### Store

#### Permission
A code uploader can specify the permission of code for instantiation
```go
const (
	// AccessTypeUnspecified placeholder for empty value
	AccessTypeUnspecified AccessType = 0
	// AccessTypeNobody forbidden
	AccessTypeNobody AccessType = 1
	// AccessTypeOnlyAddress restricted to an address
	AccessTypeOnlyAddress AccessType = 2
	// AccessTypeEverybody unrestricted
	AccessTypeEverybody AccessType = 3
)
```

### Instantiate 

#### Reply
The contracts, which are using reply to check instantiated contract address, 
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

Ex) https://github.com/terraswap/terraswap/pull/47

#### Event
Event key for instantiated contract also should be changed from `instantiate`.`contract_address` to `instantiate`.`_contract_address`.

#### Label
Now label is used to represent the contract info

#### Burn Operation
`CosmosMsg::Bank(BankMsg::Burn)` is enabled

### Execute

#### Event
Event key for instantiated contract also should be changed from `execute`.`contract_address` to `execute`.`_contract_address`.


### Migrate

#### Event
Event key for instantiated contract also should be changed from `migrate`.`contract_address` to `migrate`.`_contract_address`.

### Reply

#### Event
Event key for instantiated contract also should be changed from `reply`.`contract_address` to `reply`.`_contract_address`.

