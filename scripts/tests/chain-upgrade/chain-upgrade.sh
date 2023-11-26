#!/bin/bash

OLD_VERSION=release/v2.7
UPGRADE_HEIGHT=30
CHAIN_ID=pisco-1
ROOT=$(pwd)
CHAIN_HOME=$ROOT/_build/.testnet
DENOM=uluna
SOFTWARE_UPGRADE_NAME="v2.8"
GOV_PERIOD="10s"

VAL_MNEMONIC_1="clock post desk civil pottery foster expand merit dash seminar song memory figure uniform spice circle try happy obvious trash crime hybrid hood cushion"
WALLET_MNEMONIC_1="banner spread envelope side kite person disagree path silver will brother under couch edit food venture squirrel civil budget number acquire point work mass"

export OLD_BINARY=$ROOT/_build/terrad_old
export NEW_BINARY=$ROOT/_build/terrad_new

rm -rf /tmp/terra
rm -r $ROOT/_build
mkdir $ROOT/_build

# install old binary
if ! command -v $OLD_BINARY &> /dev/null
then
    mkdir -p /tmp/terra
    cd /tmp/terra
    git clone https://github.com/terra-money/core
    cd core
    git checkout $OLD_VERSION
    make build
    cp /tmp/terra/core/build/terrad $ROOT/_build/terrad_old
    cd $ROOT
fi

# install new binary
if ! command -v $NEW_BINARY &> /dev/null
then
  make build
  cp build/terrad $ROOT/_build/terrad_new
fi

# init genesis
$OLD_BINARY init test --home $CHAIN_HOME --chain-id=$CHAIN_ID
echo $VAL_MNEMONIC_1 | $OLD_BINARY keys add val1 --home $CHAIN_HOME --recover --keyring-backend=test
VAL_ADDR_1=$($OLD_BINARY keys list emi --output=json | jq .[0].address -r)

echo $WALLET_MNEMONIC_1 | $OLD_BINARY keys add wallet1 --home $CHAIN_HOME --recover --keyring-backend=test
WALLET_ADDR_1=$($OLD_BINARY keys list emi --output=json | jq .[0].address -r)

$OLD_BINARY genesis add-genesis-account $($OLD_BINARY --home $CHAIN_HOME keys show val1 --keyring-backend test -a) 100000000000uluna  --home $CHAIN_HOME
$OLD_BINARY genesis gentx val1 1000000000uluna --home $CHAIN_HOME --chain-id $CHAIN_ID --keyring-backend test
$OLD_BINARY genesis collect-gentxs --home $CHAIN_HOME

sed -i -e "s/\"max_deposit_period\": \"172800s\"/\"max_deposit_period\": \"$GOV_PERIOD\"/g" $CHAIN_HOME/config/genesis.json
sed -i -e "s/\"voting_period\": \"172800s\"/\"voting_period\": \"$GOV_PERIOD\"/g" $CHAIN_HOME/config/genesis.json

sed -i -e 's/timeout_commit = "5s"/timeout_commit = "1s"/g' $CHAIN_HOME/config/config.toml
sed -i -e 's/timeout_propose = "3s"/timeout_propose = "1s"/g' $CHAIN_HOME/config/config.toml
sed -i -e 's/index_all_keys = false/index_all_keys = true/g' $CHAIN_HOME/config/config.toml
sed -i -e 's/enable = false/enable = true/g' $CHAIN_HOME/config/app.toml
sed -i -e 's/swagger = false/swagger = true/g' $CHAIN_HOME/config/app.toml


# run old node
echo "Starting old binary on a separate process"
if [[ "$OSTYPE" == "darwin"* ]]; then
    screen -L -dmS node1 $OLD_BINARY start --log_level trace --log_format json --home $CHAIN_HOME --pruning=nothing
else
    screen -L -Logfile $CHAIN_HOME/log-screen.log -dmS node1 $OLD_BINARY start --log_level trace --log_format json --home $CHAIN_HOME --pruning=nothing
fi
#
sleep 15
#
GOV_ADDRESS=$($OLD_BINARY query auth module-account gov --output json | jq .account.base_account.address -r)
echo '{
  "messages": [
    {
      "@type": "/cosmos.upgrade.v1beta1.MsgSoftwareUpgrade",
      "authority" : "'"$GOV_ADDRESS"'",
      "plan" : {
        "name": "'"$SOFTWARE_UPGRADE_NAME"'",
        "time": "0001-01-01T00:00:00Z",
        "height": "'"$UPGRADE_HEIGHT"'",
        "upgraded_client_state": null
      }
    }
  ],
  "metadata": "",
  "deposit": "550000000'$DENOM'",
  "title": "Upgrade to '$SOFTWARE_UPGRADE_NAME'",
  "summary": "Source Code Version https://github.com/terra-money/core"
}' > $PWD/_build/software-upgrade.json

#
$OLD_BINARY tx gov submit-proposal $ROOT/_build/software-upgrade.json --from val1 --keyring-backend test --chain-id $CHAIN_ID --home $CHAIN_HOME  -y
sleep 2
$OLD_BINARY tx gov vote 1 yes --from val1 --keyring-backend test --chain-id $CHAIN_ID --home $CHAIN_HOME  -y
#
## determine block_height to halt
while true; do
    BLOCK_HEIGHT=$($OLD_BINARY status | jq '.SyncInfo.latest_block_height' -r)
    if [ $BLOCK_HEIGHT = "$UPGRADE_HEIGHT" ]; then
        # assuming running only 1 terrad
        echo "BLOCK HEIGHT = $UPGRADE_HEIGHT REACHED, STOPPING OLD ONE"
        pkill terrad_old
        break
    else
        $OLD_BINARY query gov proposal 1 --output=json | jq ".status"
        echo "BLOCK_HEIGHT = $BLOCK_HEIGHT"
        sleep 5
    fi
done
#
sleep 5
#
## run new node
$NEW_BINARY start --home $CHAIN_HOME
