<!--
order: 5
-->

# Post

The FeeShare module uses the post handler to distribute fees between developers and the community.

## FeeShare 

The [Post Decorator](/x/feeshare/post/post.go) executes custom logic after each successful WasmExecuteMsg transaction. All fees paid by a user for transaction execution are sent to the `FeeShare` module account before being redistributed to the registered contracts.

If the `x/feeshare` module is disabled or the Wasm Execute Msg transaction targets an unregistered contract, the handler returns `nil`, without performing any actions. In this case, 100% of the transaction fees remain in the `FeeCollector` module, to be distributed elsewhere.

If the `x/feeshare` module is enabled and a Wasm Execute Msg transaction targets a registered contract, the handler sends a percentage of the transaction fees (paid by the user) to the withdraw address set for that contract, or splits the fee equally among any contract involved in the transaction.

1. The user submits an Execute transaction (`MsgExecuteContract`) to a smart contract and the transaction is executed successfully
2. Check if
   * fees module is enabled
   * the smart contract is registered to receive fee split
  
3. Calculate developer fees according to the `DeveloperShares` parameter.
4. Check which denominations governance allows fees to be paid in.
5. Check which contracts the user executed that also have been registered.
6. Calculate the total amount of fees to be paid to the developer(s). If multiple contracts are involved in a transaction, the 50% reward is split evenly between all registered withdrawal addresses.
7. Distribute the remaining amount in the `FeeCollector` to validators according to the [SDK  Distribution Scheme](https://docs.cosmos.network/main/modules/distribution/03_begin_block.html#the-distribution-scheme).


## Custom Wasm

In order to distribute an equal share of fees to all contrcts involved in a transaction, a [custom Wasm module](../../wasm/README.md) was developed. 

The custom Wasm module keeps track of all contracts involved in a transaction. When a contract is executed, the custom Wasm module keeps track of each participating contract address in a list. When the transaction is completed, the `PostHandler` from the FeeShare module distributes the rewards between the listed participants, and the `PostHandler` from the custom Wasm module removes the contract addresses from the store.
