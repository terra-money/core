#!/bin/bash

echo ""
echo "###########################################"
echo "# ICA Cross Chain Delegation to Validator #"
echo "###########################################"
echo ""

BINARY=terrad
CHAIN_DIR=$(pwd)/data

WALLET_1=$($BINARY keys show wallet1 -a --keyring-backend test --home $CHAIN_DIR/test-1)
WALLET_2=$($BINARY keys show wallet2 -a --keyring-backend test --home $CHAIN_DIR/test-2)

echo "Registering ICA on chain test-1"
$BINARY tx interchain-accounts controller register connection-0 --from $WALLET_1 --chain-id test-1 --home $CHAIN_DIR/test-1 --node tcp://localhost:16657 --keyring-backend test --broadcast-mode block -y --gas 10000000 &> /dev/null

ICS_TX_RESULT="Error:"
ICS_TX_ERROR="Error:"
while [[ "$ICS_TX_ERROR" == "$ICS_TX_RESULT"* ]]; do 
    sleep 2
    echo "Waiting for the transaction to be relayed..."
    ICS_TX_RESULT=$($BINARY query interchain-accounts controller interchain-account $WALLET_1 connection-0 --home $CHAIN_DIR/test-1 --chain-id test-1 --node tcp://localhost:16657 -o json | jq -r '.address')
done

echo "Sending tokens to ICA on chain test-2"
$BINARY tx bank send $WALLET_2 $ICS_TX_RESULT 10000000uluna --chain-id test-2 --home $CHAIN_DIR/test-2 --node tcp://localhost:26657 --keyring-backend test --broadcast-mode block -y &> /dev/null
ICS_ACCOUNT_BALANCE=$($BINARY query bank balances $ICS_TX_RESULT --chain-id test-2 --node tcp://localhost:26657 -o json | jq -r '.balances[0].amount')

if [[ "$ICS_ACCOUNT_BALANCE" != "10000000" ]]; then
    echo "Error: ICA Have not received tokens"
    exit 1
fi

echo "Executing Delegation from test-1 to test-2 via ICA"
VAL_ADDR_1=$(cat $CHAIN_DIR/test-2/config/genesis.json | jq -r '.app_state.genutil.gen_txs[0].body.messages[0].validator_address')

$BINARY tx intertx submit \
'{
    "@type":"/cosmos.staking.v1beta1.MsgDelegate",
    "delegator_address": "'"$ICS_TX_RESULT"'",
    "validator_address": "'"$VAL_ADDR_1"'",
    "amount": {
        "denom": "uluna",
        "amount": "'"$ICS_ACCOUNT_BALANCE"'"
    }
}' --connection-id connection-0 --from $WALLET_1 --chain-id test-1 --home $CHAIN_DIR/test-1 --node tcp://localhost:16657 --broadcast-mode block --keyring-backend test -y &> /dev/null

VALIDATOR_DELEGATIONS=""
while [[ "$VALIDATOR_DELEGATIONS" != "$ICS_ACCOUNT_BALANCE" ]]; do 
    sleep 2
    echo "Waiting for the transaction '/cosmos.bank.v1beta1.MsgSend' to be relayed..."
    VALIDATOR_DELEGATIONS=$($BINARY query staking delegations-to $VAL_ADDR_1 --home $CHAIN_DIR/test-2 --node tcp://localhost:26657 -o json | jq -r '.delegation_responses[-1].balance.amount')
done

echo ""
echo "####################################################"
echo "# SUCCESS: ICA Cross Chain Delegation to Validator #"
echo "####################################################"
echo ""
