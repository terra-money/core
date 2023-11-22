<!--
order: 1
-->

# Concepts

## FeeShare

The FeeShare module enables the splitting of revenue from transaction fees between validators and registered smart contracts. Developers can register their smart contracts and every time someone interacts with a registered smart contract, the contract deployer or their assigned withdrawal account receives a part of the transaction fees. By default, 50% of all transaction fees for Execute Messages are shared. This parameter can be changed by governance and implemented by the `x/feeshare` module.

## Registration

Contracts must register to receive their allocation of fees. Registration is permissionless, and is completed by submitting a signed registration transaction. After the transaction is executed successfully, the withdrawal address stipulated  during registration will start receiving a portion of the transaction fees paid when a user interacts with the registered contract. A withdrawal address can be any address. 

::: tip
 **NOTE**: If your contract is part of a development project, please ensure that the deployer of the contract is an account that is owned by that project and not just an individual contributor, who could become malicious. 
:::

## Fees

Registered contracts will earn a portion of the transaction fee after registering their contracts. Only [Wasm Execute Txs](https://github.com/CosmWasm/wasmd/blob/main/proto/cosmwasm/wasm/v1/tx.proto#L115-L127) (`MsgExecuteContract`) are eligible for feesharing.

Users pay transaction fees to pay to interact with smart contracts and execute transactions. When a transaction is executed, the entire fee amount (`gas limit * gas price`) is sent to the `FeeCollector` module account during the [Cosmos SDK AnteHandler](https://docs.cosmos.network/main/modules/auth/#antehandlers) execution.

If a transaction's fees are not denominated in a coin permitted by the `AllowedDenoms` parameter, there is no payout to involved contracts.  If a user sends a message and it does not interact with any contracts (ex: bankSend), then the entire fee is sent to the `FeeCollector` as expected.

### Fee distribution


After collecting fees, the `FeeCollector` sends 50% of the total collected transaction fees divided equally among the withdrawal addresses of any contract involved in the transaction. 

The original FeeShare module implementation only rewarded registered contracts that took part in the execution of a transaction. However, this approach has been modified to reward *any* registered contract that participates in a transaction. 

All registered contracts involved in a transaction will receive an equal portion of the FeeShare allocation (currently set to 50%). For example, if a transaction involves the participation of 5 contracts, and 3 of them are registered, each registered contract will receive 1/3 of the 50% FeeShare allocation to their withdrawer addresses. 

This equitable distribution is achieved by wrapping the official WASM module in a [custom implementation](../../wasm/README.md) that keeps track of all contracts involved in a transaction. When a contract is executed, the custom WASM module keeps track of each participating contract address in a list. When the transaction is completed, the `PostHandler` from the FeeShare module distributes the rewards between the listed participants, and the `PostHandler` from the custom WASM module removes the contract addresses from the store.






