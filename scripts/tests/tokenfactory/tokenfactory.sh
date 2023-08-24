#!/bin/bash

echo ""
echo "##########################################"
echo "# Create, Delete, Mint with Tokenfactory #"
echo "##########################################"
echo ""

BINARY=terrad
CHAIN_DIR=$(pwd)/data
TOKEN_DENOM=utoken$RANDOM
MINT_AMOUNT=1000000

WALLET_1=$($BINARY keys show wallet1 -a --keyring-backend test --home $CHAIN_DIR/test-1)
WALLET_2=$($BINARY keys show wallet2 -a --keyring-backend test --home $CHAIN_DIR/test-2)
WALLET_3=$($BINARY keys show wallet3 -a --keyring-backend test --home $CHAIN_DIR/test-1)

echo "Creating token denom $TOKEN_DENOM with $WALLET_1 on chain test-1"
CREATED_RES_DENOM=$($BINARY tx tokenfactory create-denom $TOKEN_DENOM --from $WALLET_1 --home $CHAIN_DIR/test-1 --chain-id test-1 --node tcp://localhost:16657  --keyring-backend test -o json -y | jq -r '.logs[0].events[2].attributes[1].value')
if [ "$CREATED_RES_DENOM" != "factory/$WALLET_1/$TOKEN_DENOM" ]; then
    echo "ERROR: Tokenfactory creating denom error. Expected result 'factory/$WALLET_1/$TOKEN_DENOM', got '$CREATED_RES_DENOM'"
    exit 1
fi

echo "Minting $MINT_AMOUNT units of $TOKEN_DENOM with $WALLET_1 on chain test-1"
MINT_RES=$($BINARY tx tokenfactory mint $MINT_AMOUNT$CREATED_RES_DENOM --from $WALLET_1 --home $CHAIN_DIR/test-1 --chain-id test-1 --node tcp://localhost:16657  --keyring-backend test -o json -y | jq -r '.logs[0].events[2].type')
if [ "$MINT_RES" != "coinbase" ]; then
    echo "ERROR: Tokenfactory minting error. Expected result 'coinbase', got '$CREATED_RES_DENOM'"
    exit 1
fi

echo "Querying $TOKEN_DENOM from $WALLET_1 on chain test-1 to validate the amount minted"
BALANCE_RES_AMOUNT=$($BINARY query bank balances $WALLET_1 --denom $CREATED_RES_DENOM --chain-id test-2 --node tcp://localhost:16657 -o json | jq -r '.amount')
if [ "$BALANCE_RES_AMOUNT" != $MINT_AMOUNT ]; then
    echo "ERROR: Tokenfactory minting error. Expected minted balance '$MINT_AMOUNT', got '$BALANCE_RES_AMOUNT'"
    exit 1
fi

echo "Burning 1 $TOKEN_DENOM from $WALLET_1 on chain test-1"
BURN_RES=$($BINARY tx tokenfactory burn 1$CREATED_RES_DENOM --from $WALLET_1 --home $CHAIN_DIR/test-1 --chain-id test-1 --node tcp://localhost:16657  --keyring-backend test -o json -y | jq -r '.logs[0].events[4].type')
if [ "$BURN_RES" != "tf_burn" ]; then
    echo "ERROR: Tokenfactory burning error. Expected result 'tf_burn', got '$BURN_RES'"
    exit 1
fi

echo "Querying $TOKEN_DENOM from $WALLET_1 on chain test-1 to validate the burned amount"
BALANCES_AFTER_BURNING=$($BINARY query bank balances $WALLET_1 --denom $CREATED_RES_DENOM --chain-id test-2 --node tcp://localhost:16657 -o json | jq -r '.amount')
if [ "$BALANCES_AFTER_BURNING" != 999999 ]; then
    echo "ERROR: Tokenfactory minting error. Expected minted balance '999999', got '$BALANCES_AFTER_BURNING'"
    exit 1
fi

echo "Sending 1 $TOKEN_DENOM from $WALLET_1 to $WALLET_3 on chain test-1"
SEND_RES_MSG_TYPE=$($BINARY tx bank send $WALLET_1 $WALLET_3 1$CREATED_RES_DENOM --from $WALLET_1 --home $CHAIN_DIR/test-1 --chain-id test-1 --node tcp://localhost:16657  --keyring-backend test -o json -y | jq -r '.logs[0].events[2].attributes[0].value')
if [ "$SEND_RES_MSG_TYPE" != "/cosmos.bank.v1beta1.MsgSend" ]; then
    echo "ERROR: Sending expected to be '/cosmos.bank.v1beta1.MsgSend' but got '$SEND_RES_MSG_TYPE'"
    exit 1
fi

echo "Querying $TOKEN_DENOM from $WALLET_3 on chain test-1 to validate the funds were received"
BALANCES_RECEIVED=$($BINARY query bank balances $WALLET_3 --denom $CREATED_RES_DENOM --chain-id test-2 --node tcp://localhost:16657 -o json | jq -r '.amount')
if [ "$BALANCES_RECEIVED" != 1 ]; then
    echo "ERROR: Tokenfactory minting error. Expected minted balance '1', got '$BALANCES_RECEIVED'"
    exit 1
fi


echo "IBC'ing 1 $TOKEN_DENOM from $WALLET_1 chain test-1 to $WALLET_2 chain test-2"
IBC_SEND_RES=$($BINARY tx ibc-transfer transfer transfer channel-0 $WALLET_2 1$CREATED_RES_DENOM --chain-id test-1 --home $CHAIN_DIR/test-1 --node tcp://localhost:16657 --keyring-backend test --from $WALLET_1  -y -o json | jq -r '.logs[0].events[3].attributes[0].value')
if [ "$IBC_SEND_RES" != "/ibc.applications.transfer.v1.MsgTransfer" ]; then
    echo "ERROR: IBC'ing expected type '/ibc.applications.transfer.v1.MsgTransfer' but got '$IBC_SEND_RES'"
    exit 1
fi

IBC_RECEIVED_RES_AMOUNT=$($BINARY query bank balances $WALLET_2 --chain-id test-2 --node tcp://localhost:26657 -o json | jq -r '.balances[0].amount')
IBC_RECEIVED_RES_DENOM=""
while [ "$IBC_RECEIVED_RES_AMOUNT" != "1" ] || [ "${IBC_RECEIVED_RES_DENOM:0:4}" != "ibc/" ]; do
    sleep 2
    IBC_RECEIVED_RES_AMOUNT=$($BINARY query bank balances $WALLET_2 --chain-id test-2 --node tcp://localhost:26657 -o json | jq -r '.balances[0].amount')
    IBC_RECEIVED_RES_DENOM=$($BINARY query bank balances $WALLET_2 --chain-id test-2 --node tcp://localhost:26657 -o json | jq -r '.balances[0].denom')
    echo "Received:" $IBC_RECEIVED_RES_AMOUNT $IBC_RECEIVED_RES_DENOM
done

echo ""
echo "###################################################"
echo "# SUCCESS: Create, Delete, Mint with Tokenfactory #"
echo "###################################################"
echo ""
