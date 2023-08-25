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

# Deploy the smart contract on chain to test the callbacks. (find the source code under the following url: `~/scripts/tests/ibc-hooks/counter/src/contract.rs`)
echo "Deploying counter contract"
TX_HASH=$($BINARY tx wasm store $(pwd)/scripts/tests/ibc-hooks/counter/artifacts/counter.wasm --from $WALLET_2 --chain-id test-2 --home $CHAIN_DIR/test-2 --node tcp://localhost:26657 --keyring-backend test -y --gas 10000000 -o json | jq -r '.txhash')
sleep 3
CODE_ID=$($BINARY query tx $TX_HASH -o josn --chain-id test-2 --home $CHAIN_DIR/test-2 --node tcp://localhost:26657 | jq -r '.logs[0].events[1].attributes[1].value')


# Use Instantiate2 to instantiate the previous smart contract with a random hash to enable multiple instances of the same contract (when needed).
echo "Instantiating counter contract"
RANDOM_HASH=$(hexdump -vn16 -e'4/4 "%08X" 1 "\n"' /dev/urandom)
TX_HASH=$($BINARY tx wasm instantiate2 $CODE_ID '{"count": 0}' $RANDOM_HASH --no-admin --label="Label with $RANDOM_HASH" --from $WALLET_2 --chain-id test-2 --home $CHAIN_DIR/test-2 --node tcp://localhost:26657 --keyring-backend test -y --gas 10000000 -o json | jq -r '.txhash')
sleep 3
CONTRACT_ADDRESS=$($BINARY query tx $TX_HASH -o josn --chain-id test-2 --home $CHAIN_DIR/test-2 --node tcp://localhost:26657 | jq -r '.logs[0].events[1].attributes[0].value')

echo "Executing the IBC Hook to increment the counter"
# First execute an IBC transfer to create the entry in the smart contract with the sender address ...
IBC_HOOK_RES=$($BINARY tx ibc-transfer transfer transfer channel-0 $CONTRACT_ADDRESS 1uluna --memo='{"wasm":{"contract": "'"$CONTRACT_ADDRESS"'" ,"msg": {"increment": {}}}}' --chain-id test-1 --home $CHAIN_DIR/test-1 --node tcp://localhost:16657 --keyring-backend test --from $WALLET_1 -y -o json)
sleep 3
# ... then send another transfer to increments the count value from 0 to 1, send 1 more uluna to the contract address to validate that it increased the value correctly.
IBC_HOOK_RES=$($BINARY tx ibc-transfer transfer transfer channel-0 $CONTRACT_ADDRESS 1uluna --memo='{"wasm":{"contract": "'"$CONTRACT_ADDRESS"'" ,"msg": {"increment": {}}}}' --chain-id test-1 --home $CHAIN_DIR/test-1 --node tcp://localhost:16657 --keyring-backend test --from $WALLET_1 -y -o json)
export WALLET_1_WASM_SENDER=$($BINARY q ibchooks wasm-sender channel-0 "$WALLET_1" --chain-id test-1 --home $CHAIN_DIR/test-1 --node tcp://localhost:16657)

COUNT_RES=""
COUNT_FUNDS_RES=""
while [ "$COUNT_RES" != "1" ] || [ "$COUNT_FUNDS_RES" != "2" ]; do
    sleep 3
    # Query to assert that the counter value is 1 and the fund send are 2uluna (remeber that the first time fund are send to the contract the counter is set to 0 instead of 1)
    COUNT_RES=$($BINARY query wasm contract-state smart "$CONTRACT_ADDRESS" '{"get_count": {"addr": "'"$WALLET_1_WASM_SENDER"'"}}' --chain-id test-2 --home $CHAIN_DIR/test-2 --node tcp://localhost:26657 -o json |  jq -r '.data.count')
    COUNT_FUNDS_RES=$($BINARY query wasm contract-state smart "$CONTRACT_ADDRESS" '{"get_total_funds": {"addr": "'"$WALLET_1_WASM_SENDER"'"}}' --chain-id test-2 --home $CHAIN_DIR/test-2 --node tcp://localhost:26657 -o json |  jq -r '.data.total_funds[0].amount')
    echo "transaction relayed count: $COUNT_RES and relayed funds: $COUNT_FUNDS_RES"
done

echo "Executing the IBC Hook to increment the counter on callback"
# Execute an IBC transfer with ibc_callback to test the callback acknowledgement twice.
IBC_HOOK_RES=$($BINARY tx ibc-transfer transfer transfer channel-0 $WALLET_1_WASM_SENDER 1uluna --memo='{"ibc_callback":"'"$CONTRACT_ADDRESS"'"}' --chain-id test-2 --home $CHAIN_DIR/test-2 --node tcp://localhost:26657 --keyring-backend test --from $WALLET_2  -y -o json)
sleep 3
IBC_HOOK_RES=$($BINARY tx ibc-transfer transfer transfer channel-0 $WALLET_1_WASM_SENDER 1uluna --memo='{"ibc_callback":"'"$CONTRACT_ADDRESS"'"}' --chain-id test-2 --home $CHAIN_DIR/test-2 --node tcp://localhost:26657 --keyring-backend test --from $WALLET_2  -y -o json)
export WALLET_2_WASM_SENDER=$($BINARY q ibchooks wasm-sender channel-0 "$WALLET_2" --chain-id test-2 --home $CHAIN_DIR/test-2 --node tcp://localhost:26657)

COUNT_RES=""
while [ "$COUNT_RES" != "2" ]; do
    sleep 3
    # Query the smart contract to validate that it received the callback twice (notice that the queried addess is the contract address itself).
    COUNT_RES=$($BINARY query wasm contract-state smart "$CONTRACT_ADDRESS" '{"get_count": {"addr": "'"$CONTRACT_ADDRESS"'"}}' --chain-id test-2 --home $CHAIN_DIR/test-2 --node tcp://localhost:26657 -o json |  jq -r '.data.count')
    echo "relayed callback transaction count: $COUNT_RES"
done

echo "Executing the IBC Hook to increment the counter on callback with timeout"
# Prepare two callback queries but this time with a timeout height that is unreachable (0-1) to test the timeout callback.
IBC_HOOK_RES=$($BINARY tx ibc-transfer transfer transfer channel-0 $WALLET_1_WASM_SENDER 1uluna --packet-timeout-height="0-1" --memo='{"ibc_callback":"'"$CONTRACT_ADDRESS"'"}' --chain-id test-2 --home $CHAIN_DIR/test-2 --node tcp://localhost:26657 --keyring-backend test --from $WALLET_2  -y -o json)
sleep 3
IBC_HOOK_RES=$($BINARY tx ibc-transfer transfer transfer channel-0 $WALLET_1_WASM_SENDER 1uluna --packet-timeout-height="0-1" --memo='{"ibc_callback":"'"$CONTRACT_ADDRESS"'"}' --chain-id test-2 --home $CHAIN_DIR/test-2 --node tcp://localhost:26657 --keyring-backend test --from $WALLET_2  -y -o json)
export WALLET_2_WASM_SENDER=$($BINARY q ibchooks wasm-sender channel-0 "$WALLET_2" --chain-id test-2 --home $CHAIN_DIR/test-2 --node tcp://localhost:26657)

COUNT_RES=""
while [ "$COUNT_RES" != "22" ]; do
    sleep 3
    # Query the smart contract to validate that it received the timeout callback twice and keep in mind that per each timeout the contract increases 10 counts (notice that the queried addess is the contract address itself).
    COUNT_RES=$($BINARY query wasm contract-state smart "$CONTRACT_ADDRESS" '{"get_count": {"addr": "'"$CONTRACT_ADDRESS"'"}}' --chain-id test-2 --home $CHAIN_DIR/test-2 --node tcp://localhost:26657 -o json |  jq -r '.data.count')
    echo "relayed timeout callback transaction count: $COUNT_RES"
done

echo ""
echo "##########################"
echo "# SUCCESS: IBC Hook call #"
echo "##########################"
echo ""
