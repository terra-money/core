# Register a contract

Use this guide to register your contract. For a more in-depth guide on registering, visit the [Terra Docs feeshare tutorial](https://docs.terra.money/develop/guides/register-feeshare)

## Using terrad

`terrad tx feeshare register [contract_bech32] [withdraw_bech32] --from [key]`

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

# Update a Contract's Withdrawal Address

This can be changed at any time, as long as the sender still the admin or creator of the contract:

`terrad tx feeshare update [contract] [new_withdraw_address]`

## Update Exception

```text
This can not be done if the contract was created from or is administered by another contract (a contract factory). There is not currently a way for a contract to change its own withdrawal address directly.
```
