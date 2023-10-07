#!/bin/bash

echo ""
echo "#########################################"
echo "# POB re-ordering transactions on block #"
echo "#########################################"
echo ""

BINARY=terrad
CHAIN_DIR=$(pwd)/data

WALLET_1=$($BINARY keys show wallet1 -a --keyring-backend test --home $CHAIN_DIR/test-1)
WALLET_2=$($BINARY keys show wallet3 -a --keyring-backend test --home $CHAIN_DIR/test-1)

echo "Submit transactions on chain"
FIRST_TX__HASH_ON_BLOCK=$($BINARY tx bank send $WALLET_1 $WALLET_2 1uluna --from $WALLET_1 --chain-id test-1 --home $CHAIN_DIR/test-1 --node tcp://localhost:16657 --keyring-backend test -o json -y --gas 10000000 --aux | jq -r .)
SECOND_TX_HASH__ON_BLOCK=$($BINARY tx bank send $WALLET_1 $WALLET_2 2uluna --from $WALLET_1 --chain-id test-1 --home $CHAIN_DIR/test-1 --node tcp://localhost:16657 --keyring-backend test -o json -y --gas 10000000 --aux | jq -r .)
THIRD_TX__HASH_ON_BLOCK=$($BINARY tx bank send $WALLET_1 $WALLET_2 3uluna --from $WALLET_1 --chain-id test-1 --home $CHAIN_DIR/test-1 --node tcp://localhost:16657 --keyring-backend test -o json -y --gas 10000000 --aux | jq -r .)
FORTH_TX__HASH_ON_BLOCK=$($BINARY tx bank send $WALLET_1 $WALLET_2 4uluna --from $WALLET_1 --chain-id test-1 --home $CHAIN_DIR/test-1 --node tcp://localhost:16657 --keyring-backend test -o json -y --gas 10000000 --aux | jq -r .)
FIFTH_TX__HASH_ON_BLOCK=$($BINARY tx bank send $WALLET_1 $WALLET_2 5uluna --from $WALLET_1 --chain-id test-1 --home $CHAIN_DIR/test-1 --node tcp://localhost:16657 --keyring-backend test -o json -y --gas 10000000 --aux | jq -r .)

echo "Re-ordering transactions on block..."
echo $FIRST_TX__HASH_ON_BLOCK
BLOCK_REORDERING=$($BINARY tx builder auction-bid $WALLET_1 1000000uluna $FIFTH_TX__HASH_ON_BLOCK,$FORTH_TX__HASH_ON_BLOCK,$THIRD_TX__HASH_ON_BLOCK,$SECOND_TX_HASH__ON_BLOCK,$FIRST_TX__HASH_ON_BLOCK --timeout-height 1000 --from $WALLET_1 --chain-id test-1 --home $CHAIN_DIR/test-1 --node tcp://localhost:16657 --keyring-backend test -o json -y --gas 10000000 | jq -r .)

echo "Waiting for the block to be produced..."

echo "$BLOCK_REORDERING"

echo ""
echo "##################################################"
echo "# SUCCESS: POB re-ordering transactions on block #"
echo "##################################################"
echo ""
