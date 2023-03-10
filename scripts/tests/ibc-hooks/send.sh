#!/bin/bash

echo ""
echo "################################"
echo "# IBC Hook mint tokens thu IBC #"
echo "################################"
echo ""

BINARY=terrad
CHAIN_DIR=$(pwd)/data
WALLET_1=$($BINARY keys show wallet1 -a --keyring-backend test --home $CHAIN_DIR/test-1)
WALLET_3=$($BINARY keys show wallet3 -a --keyring-backend test --home $CHAIN_DIR/test-2)

CODE_ID=$($BINARY tx wasm store $(pwd)/scripts/tests/ibc-hooks/cw20_base.wasm --from $WALLET_1 --chain-id test-1 --home $CHAIN_DIR/test-1 --node tcp://localhost:16657 --keyring-backend test --broadcast-mode block  -y --gas 10000000 -o json | jq -r '.logs[0].events[1].attributes[1].value')

RANDOM_HASH=$(hexdump -vn16 -e'4/4 "%08X" 1 "\n"' /dev/urandom)
CONTRACT_ADDRESS=$($BINARY tx wasm instantiate2 $CODE_ID \
'{
    "name": "Bit Money",
    "symbol": "BTM",
    "decimals": 8,
    "initial_balances": [
        { "amount": "1000000", "address": "'"$WALLET_1"'"},
        { "amount": "1000000", "address": "'"$WALLET_3"'"}
    ],
    "mint": {
        "cap": "21000000",
        "minter": "'"$WALLET_1"'"
    },
    "marketing": {
        "description": "Whatever you want to say about your token",
        "logo": {
            "url": ""
        },
        "project": "Bit"
    }
}' $RANDOM_HASH --no-admin --label="Label with $RANDOM_HASH" --from $WALLET_1 --chain-id test-1 --home $CHAIN_DIR/test-1 --node tcp://localhost:16657 --keyring-backend test --broadcast-mode block  -y --gas 10000000 -o json | jq -r '.logs[0].events[0].attributes[0].value')

echo $CODE_ID $CONTRACT_ADDRESS

echo ""
echo "#########################################"
echo "# SUCCESS: IBC Hook mint tokens thu IBC #"
echo "#########################################"
echo ""
