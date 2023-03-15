#!/bin/bash

echo ""
echo "############################"
echo "# Validate Vesting Account #"
echo "############################"
echo ""

BINARY=terrad
CHAIN_DIR=$(pwd)/data
VESTING_FILE=$(pwd)/scripts/tests/vesting-accounts/vesting-periods.json
HIDDEN_VESTING_FILE=$(pwd)/scripts/tests/vesting-accounts/.vesting-periods.json

WALLET_3=$($BINARY keys show wallet3 -a --keyring-backend test --home $CHAIN_DIR/test-1)
WALLET_4=$($BINARY keys show wallet4 -a --keyring-backend test --home $CHAIN_DIR/test-2)

echo "Checking the delegated vesting balance of wallet3 on chain test-2 to 90000000000 since 10000000000 is vesting"
WALLET_4_BALANCES=$($BINARY query bank balances $WALLET_4 --chain-id test-2 --node tcp://localhost:26657 -o json | jq -r '.balances[-1].amount')
if [[ "$WALLET_4_BALANCES" != "90000000000" ]]; then
    echo "Error: Expected a balance of 90000000000, got $WALLET_4_BALANCES"
    exit 1
fi

echo "Checking the vesting balance of wallet3 to be staked on chain test-2 to 10000000000"
WALLET_4_DELEGATIONS=$($BINARY query staking delegations $WALLET_4 --home $CHAIN_DIR/test-2 --node tcp://localhost:26657 -o json | jq -r '.delegation_responses[-1].balance.amount')
if [[ "$WALLET_4_DELEGATIONS" != "10000000000" ]]; then
    echo "Error: Expected a total staking of of 10000000000, got $WALLET_4_DELEGATIONS"
    exit 1
fi

echo "Creating a random vesting wallet on chain test-1"
CURRENT_DATE=$(date +%s)
$BINARY keys add wallet$CURRENT_DATE --home $CHAIN_DIR/test-1 --keyring-backend=test &> /dev/null
RANDOM_VESTING_WALLET=$($BINARY keys show wallet$CURRENT_DATE -a --keyring-backend test --home $CHAIN_DIR/test-1)

cp $VESTING_FILE $HIDDEN_VESTING_FILE
sed -i -e 's/"start_time": -1/"start_time": '$CURRENT_DATE'/g' $HIDDEN_VESTING_FILE

echo "Deploying a vesting account on chain test-1 with the address $RANDOM_VESTING_WALLET"
CREATE_VESTING_ACCOUNT_MSG_RES=$($BINARY tx vesting create-periodic-vesting-account $RANDOM_VESTING_WALLET $HIDDEN_VESTING_FILE --from $WALLET_3 --chain-id test-1 --home $CHAIN_DIR/test-1 --node tcp://localhost:16657 --broadcast-mode block --keyring-backend test -y -o json | jq -r '.logs[0].events[2].attributes[0].value')
if [[ "$CREATE_VESTING_ACCOUNT_MSG_RES" != "/cosmos.vesting.v1beta1.MsgCreatePeriodicVestingAccount" ]]; then
    echo "Error: Expected a message type /cosmos.vesting.v1beta1.MsgCreatePeriodicVestingAccount, got $CREATE_VESTING_ACCOUNT_MSG_RES"
    exit 1
fi

echo "Waiting 4 seconds for address $RANDOM_VESTING_WALLET to have spendable balance"
sleep 4
SPENDABLE_BALANCE=$(curl -s -X GET "http://localhost:1316/cosmos/bank/v1beta1/spendable_balances/$RANDOM_VESTING_WALLET" -H "accept: application/json" | jq -r '.balances[-1].amount')
if [[ 0 -gt $SPENDABLE_BALANCE ]]; then
    echo "Error: Expected a vested balance greater than 0, got $SPENDABLE_BALANCE"
    exit 1
fi  
echo "Spendable balance of $RANDOM_VESTING_WALLET is $SPENDABLE_BALANCE"

rm -rf $HIDDEN_VESTING_FILE

echo ""
echo "#####################################"
echo "# SUCCESS: Validate Vesting Account #"
echo "#####################################"
echo ""
