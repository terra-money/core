# Store the following account addresses within the current shell env
export WALLET_1=$(terrad keys show wallet1 -a --keyring-backend test --home ./data/test-1) && echo $WALLET_1;
export WALLET_2=$(terrad keys show wallet2 -a --keyring-backend test --home ./data/test-1) && echo $WALLET_2;
export WALLET_3=$(terrad keys show wallet3 -a --keyring-backend test --home ./data/test-2) && echo $WALLET_3;
export WALLET_4=$(terrad keys show wallet4 -a --keyring-backend test --home ./data/test-2) && echo $WALLET_4;

## REGISTER ICA ON chain test-1
terrad tx interchain-accounts controller register connection-0 --from $WALLET_1 --chain-id test-1 --home ./data/test-1 --node tcp://localhost:16657 --keyring-backend test -y --gas 10000000
terrad query interchain-accounts controller interchain-account $WALLET_1 connection-0 --home ./data/test-1 --chain-id test-1 --node tcp://localhost:16657

## SEND TOKENS TO THE ICA ON chain test-2
export ICA_ADDR=$(terrad query interchain-accounts controller interchain-account $WALLET_1 connection-0 --home ./data/test-1 --node tcp://localhost:16657 -o json | jq -r '.address') && echo $ICA_ADDR
terrad q bank balances $ICA_ADDR --chain-id test-2 --node tcp://localhost:26657
terrad tx bank send $WALLET_3 $ICA_ADDR 10000000uluna --chain-id test-2 --home ./data/test-2 --node tcp://localhost:26657 --keyring-backend test -y
terrad q bank balances $ICA_ADDR --chain-id test-2 --node tcp://localhost:26657

### SUBMIT ICA TX FROM CHAIN test-1 TO test-2 
export VAL_ADDR_1=$(cat ./data/test-2/config/genesis.json | jq -r '.app_state.genutil.gen_txs[0].body.messages[0].validator_address') && echo $VAL_ADDR_1

terrad tx intertx submit \
'{
    "@type":"/cosmos.staking.v1beta1.MsgDelegate",
    "delegator_address": "'"$ICA_ADDR"'",
    "validator_address": "'"$VAL_ADDR_1"'",
    "amount": {
        "denom": "uluna",
        "amount": "1000"
    }
}' --connection-id connection-0 --from $WALLET_1 --chain-id test-1 --home ./data/test-1 --node tcp://localhost:16657 --keyring-backend test -y

terrad q staking delegations-to $VAL_ADDR_1 --home ./data/test-2 --node tcp://localhost:26657