#!/bin/bash

echo ""
echo "#################################################"
echo "# Alliance: bridge funds and create an alliance #"
echo "#################################################"
echo ""

BINARY=terrad
CHAIN_DIR=$(pwd)/data

AMOUNT_TO_DELEGATE=10000000000
ULUNA_DENOM=uluna
VAL_WALLET_1=$($BINARY keys show val1 -a --keyring-backend test --home $CHAIN_DIR/test-1)
VAL_WALLET_2=$($BINARY keys show val2 -a --keyring-backend test --home $CHAIN_DIR/test-2)

echo "Sending tokens from validator wallet on test-1 to validator wallet on test-2"
IBC_TRANSFER=$($BINARY tx ibc-transfer transfer transfer channel-0 $VAL_WALLET_2 $AMOUNT_TO_DELEGATE$ULUNA_DENOM --chain-id test-1 --from $VAL_WALLET_1 --home $CHAIN_DIR/test-1 --node tcp://localhost:16657 --keyring-backend test --broadcast-mode block -y -o json | jq -r '.raw_log' )

if [[ "$IBC_TRANSFER" == "failed to execute message"* ]]; then
    echo "Error: IBC transfer failed, with error: $IBC_TRANSFER"
    exit 1
fi

ACCOUNT_BALANCE=""
IBC_DENOM=""
while [ "$ACCOUNT_BALANCE" == "" ]; do
    IBC_DENOM=$($BINARY q bank balances $VAL_WALLET_2 --chain-id test-2 --node tcp://localhost:26657 -o json | jq -r '.balances[0].denom')
    if [ "$IBC_DENOM" != "$ULUNA_DENOM" ]; then
        ACCOUNT_BALANCE=$($BINARY q bank balances $VAL_WALLET_2 --chain-id test-2 --node tcp://localhost:26657 -o json | jq -r '.balances[0].amount')
    fi
    sleep 2
done

echo "Creating an alliance with the denom $IBC_DENOM"
PROPOSAL_HEIGHT=$($BINARY tx gov submit-legacy-proposal create-alliance $IBC_DENOM 5 0 5 0 0.99 1s --from=$VAL_WALLET_2 --home $CHAIN_DIR/test-2 --deposit 10000000000$ULUNA_DENOM --node tcp://localhost:26657 -o json --keyring-backend test --broadcast-mode block --gas 1000000 -y | jq -r '.height')
PROPOSAL_ID=$($BINARY query gov proposals --home $CHAIN_DIR/test-2 --count-total --node tcp://localhost:26657 -o json --output json --chain-id=test-2 | jq .proposals[-1].id -r)
VOTE_RES=$($BINARY tx gov vote $PROPOSAL_ID yes --from=$VAL_WALLET_2 --home $CHAIN_DIR/test-2 --keyring-backend=test --broadcast-mode=block --gas 1000000 --chain-id=test-2 --node tcp://localhost:26657 -o json -y)

ALLIANCE="null"
while [ "$ALLIANCE" == "null" ]; do
    echo "Waiting for alliance with denom $IBC_DENOM to be created"
    ALLIANCE=$($BINARY q alliance alliances --chain-id test-2 --node tcp://localhost:26657 -o json | jq -r '.alliances[0]')
    sleep 2
done

echo "Delegating $AMOUNT_TO_DELEGATE to the alliance $IBC_DENOM"
VAL_ADDR=$(allianced query staking validators --output json | jq .validators[0].operator_address --raw-output)
DELEGATE_RES=$($BINARY tx alliance delegate $VAL_ADDR $AMOUNT_TO_DELEGATE$IBC_DENOM --from=node0 --from=$VAL_WALLET_2 --home $CHAIN_DIR/test-2 --keyring-backend=test --broadcast-mode=block --gas 1000000 --chain-id=test-2 -o json  -y)

DELEGATIONS=$($BINARY query alliance delegation $VAL_WALLET_2 $VAL_ADDR $IBC_DENOM --chain-id test-2 --node tcp://localhost:26657 -o json | jq -r '.delegation.balance.amount')
if [[ "$DELEGATIONS" == "0" ]]; then
    echo "Error: Alliance delegations expected to be bigger than 0"
    exit 1
fi

echo "Query bank balance after alliance creation"
TOTAL_SUPPLY_BEFORE_ALLIANCE=$($BINARY query bank total --denom $ULUNA_DENOM --height $PROPOSAL_HEIGHT -o json | jq -r '.amount')
TOTAL_SUPPLY_AFTER_ALLIANCE=$($BINARY query bank total --denom $ULUNA_DENOM -o json | jq -r '.amount')
TOTAL_SUPPLY_INCREMENT=$(($TOTAL_SUPPLY_BEFORE_ALLIANCE - $TOTAL_SUPPLY_AFTER_ALLIANCE))


if [ "$TOTAL_SUPPLY_INCREMENT" -gt 100000 ] && [ "$TOTAL_SUPPLY_INCREMENT" -lt 1000000 ]; then
    echo "Error: Something went wrong, total supply of $ULUNA_DENOM has increased out of range 100_000 between 1_000_000. current value $TOTAL_SUPPLY_INCREMENT"
    exit 1
fi

echo ""
echo "#########################################################"
echo "# Success: Alliance bridge funds and create an alliance #"
echo "#########################################################"
echo ""