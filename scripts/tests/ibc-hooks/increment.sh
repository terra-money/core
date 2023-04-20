#!/bin/bash

echo ""
echo "#################"
echo "# IBC Hook call #"
echo "#################"
echo ""

BINARY=terrad
CHAIN_DIR=$(pwd)/data
WALLET_1=$($BINARY keys show wallet1 -a --keyring-backend test --home $CHAIN_DIR/test-1)
WALLET_2=$($BINARY keys show wallet2 -a --keyring-backend test --home $CHAIN_DIR/test-2)

echo "Deploying counter contract"
CODE_ID=$($BINARY tx wasm store $(pwd)/scripts/tests/ibc-hooks/counter/artifacts/counter.wasm --from $WALLET_2 --chain-id test-2 --home $CHAIN_DIR/test-2 --node tcp://localhost:26657 --keyring-backend test --broadcast-mode block  -y --gas 10000000 -o json | jq -r '.logs[0].events[1].attributes[1].value')

echo "Instantiating counter contract"
RANDOM_HASH=$(hexdump -vn16 -e'4/4 "%08X" 1 "\n"' /dev/urandom)
CONTRACT_ADDRESS=$($BINARY tx wasm instantiate2 $CODE_ID '{"count": 0}' $RANDOM_HASH --no-admin --label="Label with $RANDOM_HASH" --from $WALLET_2 --chain-id test-2 --home $CHAIN_DIR/test-2 --node tcp://localhost:26657 --keyring-backend test --broadcast-mode block  -y --gas 10000000 -o json | jq -r '.logs[0].events[0].attributes[0].value')

echo "Executing the IBC Hook to increment the counter"
# First increment create the entry in the smart contract with the sender address ...
IBC_HOOK_RES=$($BINARY tx ibc-transfer transfer transfer channel-0 $CONTRACT_ADDRESS 1uluna --memo='{"wasm":{"contract": "'"$CONTRACT_ADDRESS"'" ,"msg": {"increment": {}}}}' --chain-id test-1 --home $CHAIN_DIR/test-1 --node tcp://localhost:16657 --keyring-backend test --from $WALLET_1 --broadcast-mode block -y -o json)
# ... then it increments the value from 0 to 1 and send 1 more luna to the contract address.
IBC_HOOK_RES=$($BINARY tx ibc-transfer transfer transfer channel-0 $CONTRACT_ADDRESS 1uluna --memo='{"wasm":{"contract": "'"$CONTRACT_ADDRESS"'" ,"msg": {"increment": {}}}}' --chain-id test-1 --home $CHAIN_DIR/test-1 --node tcp://localhost:16657 --keyring-backend test --from $WALLET_1 --broadcast-mode block -y -o json)
export WALLET_1_WASM_SENDER=$($BINARY q ibchooks wasm-sender channel-0 "$WALLET_1" --chain-id test-1 --home $CHAIN_DIR/test-1 --node tcp://localhost:16657)

COUNT_RES=""
COUNT_FUNDS_RES=""
while [ "$COUNT_RES" != "1" ] || [ "$COUNT_FUNDS_RES" != "2" ]; do
    sleep 2
    COUNT_RES=$($BINARY query wasm contract-state smart "$CONTRACT_ADDRESS" '{"get_count": {"addr": "'"$WALLET_1_WASM_SENDER"'"}}' --chain-id test-2 --home $CHAIN_DIR/test-2 --node tcp://localhost:26657 -o json |  jq -r '.data.count')
    COUNT_FUNDS_RES=$($BINARY query wasm contract-state smart "$CONTRACT_ADDRESS" '{"get_total_funds": {"addr": "'"$WALLET_1_WASM_SENDER"'"}}' --chain-id test-2 --home $CHAIN_DIR/test-2 --node tcp://localhost:26657 -o json |  jq -r '.data.total_funds[0].amount')
    echo "relayed count: $COUNT_RES relayed funds: $COUNT_FUNDS_RES"
done

echo ""
echo "##########################"
echo "# SUCCESS: IBC Hook call #"
echo "##########################"
echo ""
