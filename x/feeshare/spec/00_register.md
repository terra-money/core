# Register a contract

`junod tx feeshare register [contract_bech32] [withdraw_bech32] --from [key]`

Registers the withdrawal address for the given contract.

## Parameters

`contract_bech32 (string, required)`: The bech32 address of the contract whose interaction fees will be shared.

`withdraw_bech32 (string, required)`: The bech32 address where the interaction fees will be sent every block.

## Description

This command registers the withdrawal address for the given contract. Any time a user interacts with your contract, the funds will be sent to the withdrawal address. It can be any valid address, such as a DAO, normal account, another contract, or a multi-sig.

## Permissions

This command can only be run by the admin of the contract. If there is no admin, then it can only be run by the contract creator.

## Exceptions

```text
withdraw_bech32 can not be the community pool (distribution) address. This is a limitation of the way the SDK handles this module account
```

```text
For contracts created or administered by a contract factory, the withdrawal address can only be the same as the contract address. This can be registered by anyone, but it's unchangeable. This is helpful for SubDAOs or public goods to save fees in the treasury.

If you create a contract like this, it's best to create an execution method for withdrawing fees to an account. To do this, you'll need to save the withdrawal address in the contract's state before uploading a non-migratable contract.
```
