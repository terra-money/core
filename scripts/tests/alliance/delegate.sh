#!/bin/bash

echo ""
echo "#################################################"
echo "# Alliance: bridge funds and create an alliance #"
echo "#################################################"
echo ""

BINARY=terrad
CHAIN_DIR=$(pwd)/data

VAL_WALLET_1=$($BINARY keys show val1 -a --keyring-backend test --home $CHAIN_DIR/test-1)
VAL_WALLET_2=$($BINARY keys show val2 -a --keyring-backend test --home $CHAIN_DIR/test-2)

echo "Sending tokens from validator wallet on test-1 to validator wallet on test-2"
IBC_TRANSFER=$($BINARY tx ibc-transfer transfer transfer channel-0 $VAL_WALLET_2 10000000uluna --chain-id test-1 --from $VAL_WALLET_1 --home $CHAIN_DIR/test-1 --node tcp://localhost:16657 --keyring-backend test --broadcast-mode block -y)

ACCOUNT_BALANCE=""
IBC_DENOM=""
while [ "$ACCOUNT_BALANCE" == "" ]; do
    IBC_DENOM=$($BINARY q bank balances $VAL_WALLET_2 --chain-id test-2 --node tcp://localhost:26657 -o json | jq -r '.balances[0].denom')
    if [ "$IBC_DENOM" != "uluna" ]; then
        ACCOUNT_BALANCE=$($BINARY q bank balances $VAL_WALLET_2 --chain-id test-2 --node tcp://localhost:26657 -o json | jq -r '.balances[0].amount')
    fi
    sleep 1
done

echo "Creating an alliance with the denom $IBC_DENOM"
PROPOSAL=$($BINARY tx gov submit-legacy-proposal create-alliance $IBC_DENOM 0.5 0 1 0 0.1 1s --from=$VAL_WALLET_2 --home $CHAIN_DIR/test-2 --deposit 10000000000uluna --node tcp://localhost:26657 -o json --keyring-backend test --broadcast-mode block --gas 1000000 -y )
PROPOSAL_ID=$($BINARY query gov proposals --home $CHAIN_DIR/test-2 --count-total --node tcp://localhost:26657 -o json --output json --chain-id=test-2 | jq .proposals[-1].id -r)
VOTE_RES=$($BINARY tx gov vote $PROPOSAL_ID yes --from=$VAL_WALLET_2 --home $CHAIN_DIR/test-2 --keyring-backend=test --broadcast-mode=block --gas 1000000 --chain-id=test-2 --node tcp://localhost:26657 -o json -y)

ALLIANCE="null"
while [ "$ALLIANCE" == "null" ]; do
    ALLIANCE=$($BINARY q alliance alliances --chain-id test-2 --node tcp://localhost:26657 -o json | jq -r '.alliances[0]')
    sleep 1
done

echo "Delegating 10000000 to the alliance $IBC_DENOM"
VAL_ADDR=$(allianced query staking validators --output json | jq .validators[0].operator_address --raw-output)
DELEGATE_RES=$($BINARY tx alliance delegate $VAL_ADDR 10000000$IBC_DENOM --from=node0 --from=$VAL_WALLET_2 --home $CHAIN_DIR/test-2 --keyring-backend=test --broadcast-mode=block --gas 1000000 --chain-id=test-2 -o json  -y)

DELEGATION=""
while [ "$DELEGATION" == "" ]; do
    DELEGATION=$($BINARY query alliance delegation $VAL_WALLET_2 $VAL_ADDR $IBC_DENOM --chain-id test-2 --node tcp://localhost:26657 -o json | jq -r '.delegation.balance.amount')
    sleep 1
done

echo ""
echo "#########################################################"
echo "# Success: Alliance bridge funds and create an alliance #"
echo "#########################################################"
echo ""