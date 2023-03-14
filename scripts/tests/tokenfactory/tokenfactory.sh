#!/bin/bash

echo ""
echo "##########################################"
echo "# Create, Delete, Mint with Tokenfactory #"
echo "##########################################"
echo ""

BINARY=terrad
CHAIN_DIR=$(pwd)/data
MINT_DENOM=utoken
MINT_AMOUNT=10000.00

WALLET_3=$($BINARY keys show wallet3 -a --keyring-backend test --home $CHAIN_DIR/test-1)

echo "Minting $MINT_DENOM with $WALLET_3 on chain test-1"
CREATE_RES=$($BINARY tx tokenfactory create-denom $MINT_DENOM --from $WALLET_3 --home $CHAIN_DIR/test-1 --chain-id test-1 --node tcp://localhost:16657)

echo $CREATE_RES

exit 1
echo ""
echo "###################################################"
echo "# SUCCESS: Create, Delete, Mint with Tokenfactory #"
echo "###################################################"
echo ""
