#!/bin/bash

echo ""
echo "##################"
echo "# Create relayer #"
echo "##################"
echo ""

# Configure predefined mnemonic pharses
BINARY=relayer
CHAIN_DIR=$(pwd)/src/test-data
CHAINID_1=test-1
CHAINID_2=test-2
MNEMONIC_1="alley afraid soup fall idea toss can goose become valve initial strong forward bright dish figure check leopard decide warfare hub unusual join cart"
MNEMONIC_2="record gift you once hip style during joke field prize dust unique length more pencil transfer quit train device arrive energy sort steak upset"

# Ensure rly is installed
if ! [ -x "$(command -v $BINARY)" ]; then
    echo "$BINARY is required to run this script..."
    echo "Installing go relayer https://github.com/cosmos/relayer"
    go install github.com/cosmos/relayer/v2@v2.4.2
fi

echo "Initializing $BINARY..."
$BINARY config init --home $CHAIN_DIR/relayer

echo "Adding configurations for both chains..."
$BINARY chains add-dir ./src/setup/relayer/chains --home $CHAIN_DIR/relayer
$BINARY paths add $CHAINID_1 $CHAINID_2 test1-test2 --file ./src/setup/relayer/paths/test1-test2.json --home $CHAIN_DIR/relayer

echo "Restoring accounts..."
$BINARY keys restore $CHAINID_1 testkey "$MNEMONIC_1" --home $CHAIN_DIR/relayer
$BINARY keys restore $CHAINID_2 testkey "$MNEMONIC_2" --home $CHAIN_DIR/relayer

echo "Creating clients and a connection..."
$BINARY tx connection test1-test2 --home $CHAIN_DIR/relayer

echo "Creating a channel..."
$BINARY tx channel test1-test2 --home $CHAIN_DIR/relayer

echo "Starting to listen relayer..."
$BINARY start test1-test2 -p events -b 100 --flush-interval 1s --time-threshold 1s --home $CHAIN_DIR/relayer > $CHAIN_DIR/relayer.log 2>&1 &

sleep 10

echo ""
echo "############################"
echo "# SUCCESS: Relayer created #"
echo "############################"
echo ""
