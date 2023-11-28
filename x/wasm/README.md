# Custom WASM module

This module is a wrapper for the official WASM module, used to extend the functionality of the FeeShare module. The original FeeShare module implementation only rewarded registered contracts that took part in the execution of a transaction. However, this approach has been modified using the Custom WASM module wrapper to reward all registered contracts that participate in a transaction. 

When a contract is executed, the custom WASM module keeps track of each participating contract address in a list. When the transaction is completed, the `PostHandler` from the FeeShare module distributes the rewards between the listed participants, and the `PostHandler` from the custom WASM module removes the contract addresses from the store.

For more information on the FeeShare module, visit the [Feeshare spec](../feeshare/spec/README.md). 