#!/bin/bash

OLD_VERSION=v2.3
UPGRADE_HEIGHT=30
CHAIN_ID=pisco-1
CHAIN_HOME=.testnet
ROOT=$(pwd)
DENOM=uluna
SOFTWARE_UPGRADE_NAME="v2.5"
GOV_PERIOD="10s"

VAL_MNEMONIC_1="clock post desk civil pottery foster expand merit dash seminar song memory figure uniform spice circle try happy obvious trash crime hybrid hood cushion"
WALLET_MNEMONIC_1="banner spread envelope side kite person disagree path silver will brother under couch edit food venture squirrel civil budget number acquire point work mass"


# underscore so that go tool will not take gocache into account
mkdir -p _build/gocache

# install old binary
if ! command -v _build/old/terrad &> /dev/null
then
    mkdir -p _build/old
    wget -c "https://github.com/terra-money/core/archive/refs/tags/${OLD_VERSION}.zip" -O _build/${OLD_VERSION}.zip
    unzip _build/${OLD_VERSION}.zip -d _build
    cd ./_build/core-${OLD_VERSION:1}
    make build
    cp build/terrad ../old
    cd ../..
fi

# install new binary
if ! command -v _build/new/terrad &> /dev/null
then
  mkdir -p _build/new
  make build
  cp build/terrad _build/new
fi

export OLD_BINARY=_build/old/terrad
export NEW_BINARY=_build/new/terrad

rm -rf $CHAIN_HOME
# init genesis
$OLD_BINARY init test --home $CHAIN_HOME --chain-id=$CHAIN_ID
echo $VAL_MNEMONIC_1 | $OLD_BINARY keys add val1 --home $CHAIN_HOME --recover --keyring-backend=test
echo $WALLET_MNEMONIC_1 | $OLD_BINARY keys add wallet1 --home $CHAIN_HOME --recover --keyring-backend=test
$OLD_BINARY genesis add-genesis-account $($OLD_BINARY --home $CHAIN_HOME keys show val1 --keyring-backend test -a) 100000000000uluna  --home $CHAIN_HOME
$OLD_BINARY genesis gentx val1 7000000000uluna --home $CHAIN_HOME --chain-id $CHAIN_ID --keyring-backend test
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
$OLD_BINARY tx gov submit-proposal software-upgrade "$SOFTWARE_UPGRADE_NAME" --upgrade-height $UPGRADE_HEIGHT --upgrade-info "temp" --title "upgrade" --description "upgrade"  --from val1 --keyring-backend test --chain-id $CHAIN_ID --home $CHAIN_HOME --broadcast-mode block -y
$OLD_BINARY tx gov deposit 1 "20000000${DENOM}" --from val1 --keyring-backend test --chain-id $CHAIN_ID --home $CHAIN_HOME --broadcast-mode block -y
$OLD_BINARY tx gov vote 1 yes --from val1 --keyring-backend test --chain-id $CHAIN_ID --home $CHAIN_HOME --broadcast-mode block -y
#
## determine block_height to halt
while true; do
    BLOCK_HEIGHT=$($OLD_BINARY status | jq '.SyncInfo.latest_block_height' -r)
    if [ $BLOCK_HEIGHT = "$UPGRADE_HEIGHT" ]; then
        # assuming running only 1 terrad
        echo "BLOCK HEIGHT = $UPGRADE_HEIGHT REACHED, KILLING OLD ONE"
        pkill terrad
        break
    else
        $OLD_BINARY q gov proposal 1 --output=json | jq ".status"
        echo "BLOCK_HEIGHT = $BLOCK_HEIGHT"
        sleep 5
    fi
done
#
sleep 5
#
## run new node
$NEW_BINARY start --home $CHAIN_HOME
