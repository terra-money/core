# Counter contract from [Osmosis Labs](https://github.com/osmosis-labs/osmosis/commit/64393a14e18b2562d72a3892eec716197a3716c7)

This contract is a modification of the standard cosmwasm `counter` contract.
Namely, it tracks a counter, _by sender_.
This is a better way to test wasmhooks.

This contract tracks any funds sent to it by adding it to the state under the `sender` key.

This way we can verify that, independently of the sender, the funds will end up under the 
`WasmHooksModuleAccount` address when the contract is executed via an IBC send that goes 
through the wasmhooks module.
