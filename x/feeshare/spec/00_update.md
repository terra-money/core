# Update a Contract's Withdrawal Address

This can be changed at any time so long as you are still the admin or creator of a contract with the command:

`junod tx feeshare update [contract] [new_withdraw_address]`

## Update Exception

```text
This can not be done if the contract was created from or is administered by another contract (a contract factory). There is not currently a way for a contract to change its own withdrawal address directly.
```
