#!/bin/bash

BINARY=terrad
CHAIN_DIR=$(pwd)/src/test-data
CHAINID_1=test-1
CHAINID_2=test-2

# Chain1
RLY_MNEMONIC_1="alley afraid soup fall idea toss can goose become valve initial strong forward bright dish figure check leopard decide warfare hub unusual join cart"
VAL_MNEMONIC_1="clock post desk civil pottery foster expand merit dash seminar song memory figure uniform spice circle try happy obvious trash crime hybrid hood cushion"

WALLET_MNEMONIC_1="banner spread envelope side kite person disagree path silver will brother under couch edit food venture squirrel civil budget number acquire point work mass"
WALLET_MNEMONIC_2="veteran try aware erosion drink dance decade comic dawn museum release episode original list ability owner size tuition surface ceiling depth seminar capable only"
WALLET_MNEMONIC_3="vacuum burst ordinary enact leaf rabbit gather lend left chase park action dish danger green jeans lucky dish mesh language collect acquire waste load"
WALLET_MNEMONIC_4="open attitude harsh casino rent attitude midnight debris describe spare cancel crisp olive ride elite gallery leaf buffalo sheriff filter rotate path begin soldier"
WALLET_MNEMONIC_5="same heavy travel border destroy catalog music manual love festival exile resist always gas off coffee crystal provide random harvest sea cloud child field"
WALLET_MNEMONIC_6="broken title little open demand ladder mimic keen execute word couple door relief rule pulp demand believe cactus swing fluid tired what crop purse"
WALLET_MNEMONIC_7="unit question bulk desk slush answer share bird earth brave book wing special gorilla ozone release permit mercy luxury version advice impact unfair drama"
WALLET_MNEMONIC_8="year aim panel oyster sunny faint dress skin describe chair guilt possible venue pottery inflict mass debate poverty multiply pulse ability purse situate inmate"
WALLET_MNEMONIC_9="leave side blue panel curve ancient suspect slide seminar neutral doctor boring only curious spell surround remind obtain slogan hire giant soccer crunch system"
WALLET_MNEMONIC_10="degree under tray object thought mercy mushroom captain bus work faint basic twice cube noble man ripple close flush bunker dish spare hungry arm"
WALLET_MNEMONIC_11="range struggle season mesh antenna delay sell light yard path risk curve brain nut cabin injury dilemma fun comfort crumble today transfer bring draft"
WALLET_MNEMONIC_12="giraffe trim element wheel cannon nothing enrich shiver upon output iron recall already fix appear produce fix behind scissors artefact excite tennis into side"
WALLET_MNEMONIC_13="run turn cup combine sad toast roof already melt chimney arctic save avocado theory bracket cherry cotton fee once favorite swarm ignore dream element"
WALLET_MNEMONIC_14="script key fold coyote cage squirrel prevent pole auction slide vintage shoot mirror erosion equip goose capable critic test space sketch monkey eight candy"
WALLET_MNEMONIC_15="work clap clarify edit explain exact depth ramp law hard feel beauty stumble occur prevent crush distance purpose scrap current describe skirt panther skirt"

# Chain2
VAL_MNEMONIC_2="angry twist harsh drastic left brass behave host shove marriage fall update business leg direct reward object ugly security warm tuna model broccoli choice"
RLY_MNEMONIC_2="record gift you once hip style during joke field prize dust unique length more pencil transfer quit train device arrive energy sort steak upset"

P2PPORT_1=16656
P2PPORT_2=26656
RPCPORT_1=16657
RPCPORT_2=26657
RESTPORT_1=1316
RESTPORT_2=1317
GRPCPORT_1=8090
GRPCPORT_2=9090
GRPCWEB_1=8091
GRPCWEB_2=9091

# Stop if it is already running 
if pgrep -x "$BINARY" >/dev/null; then
    echo "Terminating $BINARY..."
    killall $BINARY
fi

echo "Removing previous data..."
rm -rf $CHAIN_DIR

echo "Stopping any processes related with the tests..."
pkill terrad
pkill terrad
pkill relayer

# Add directories for both chains, exit if an error occurs
if ! mkdir -p $CHAIN_DIR/$CHAINID_1 2>/dev/null; then
    echo "Failed to create chain folder. Aborting..."
    exit 1
fi

if ! mkdir -p $CHAIN_DIR/$CHAINID_2 2>/dev/null; then
    echo "Failed to create chain folder. Aborting..."
    exit 1
fi

echo "Initializing $CHAINID_1 & $CHAINID_2..."
$BINARY init test --home $CHAIN_DIR/$CHAINID_1 --chain-id=$CHAINID_1 &> /dev/null
$BINARY init test --home $CHAIN_DIR/$CHAINID_2 --chain-id=$CHAINID_2 &> /dev/null

echo "Adding genesis accounts..."
## Special wallets for validator and relayer
echo $VAL_MNEMONIC_1 | $BINARY keys add val1 --home $CHAIN_DIR/$CHAINID_1 --recover --keyring-backend=test
echo $RLY_MNEMONIC_1 | $BINARY keys add rly1 --home $CHAIN_DIR/$CHAINID_1 --recover --keyring-backend=test
echo $VAL_MNEMONIC_2 | $BINARY keys add val2 --home $CHAIN_DIR/$CHAINID_2 --recover --keyring-backend=test
echo $RLY_MNEMONIC_2 | $BINARY keys add rly2 --home $CHAIN_DIR/$CHAINID_2 --recover --keyring-backend=test 

## Other wallets with funds at genesis
echo $WALLET_MNEMONIC_1 | $BINARY keys add wallet1 --home $CHAIN_DIR/$CHAINID_1 --recover --keyring-backend=test
echo $WALLET_MNEMONIC_2 | $BINARY keys add wallet2 --home $CHAIN_DIR/$CHAINID_1 --recover --keyring-backend=test
echo $WALLET_MNEMONIC_3 | $BINARY keys add wallet3 --home $CHAIN_DIR/$CHAINID_1 --recover --keyring-backend=test
echo $WALLET_MNEMONIC_4 | $BINARY keys add wallet4 --home $CHAIN_DIR/$CHAINID_1 --recover --keyring-backend=test
echo $WALLET_MNEMONIC_5 | $BINARY keys add wallet5 --home $CHAIN_DIR/$CHAINID_1 --recover --keyring-backend=test
echo $WALLET_MNEMONIC_6 | $BINARY keys add wallet6 --home $CHAIN_DIR/$CHAINID_1 --recover --keyring-backend=test
echo $WALLET_MNEMONIC_7 | $BINARY keys add wallet7 --home $CHAIN_DIR/$CHAINID_1 --recover --keyring-backend=test
echo $WALLET_MNEMONIC_8 | $BINARY keys add wallet8 --home $CHAIN_DIR/$CHAINID_1 --recover --keyring-backend=test
echo $WALLET_MNEMONIC_9 | $BINARY keys add wallet9 --home $CHAIN_DIR/$CHAINID_1 --recover --keyring-backend=test
echo $WALLET_MNEMONIC_10 | $BINARY keys add wallet10 --home $CHAIN_DIR/$CHAINID_1 --recover --keyring-backend=test
echo $WALLET_MNEMONIC_11 | $BINARY keys add wallet11 --home $CHAIN_DIR/$CHAINID_1 --recover --keyring-backend=test
echo $WALLET_MNEMONIC_12 | $BINARY keys add wallet12 --home $CHAIN_DIR/$CHAINID_1 --recover --keyring-backend=test
echo $WALLET_MNEMONIC_13 | $BINARY keys add wallet13 --home $CHAIN_DIR/$CHAINID_1 --recover --keyring-backend=test
echo $WALLET_MNEMONIC_14 | $BINARY keys add wallet14 --home $CHAIN_DIR/$CHAINID_1 --recover --keyring-backend=test
echo $WALLET_MNEMONIC_15 | $BINARY keys add wallet15 --home $CHAIN_DIR/$CHAINID_1 --recover --keyring-backend=test

echo $WALLET_MNEMONIC_1 | $BINARY keys add wallet1 --home $CHAIN_DIR/$CHAINID_2 --recover --keyring-backend=test
echo $WALLET_MNEMONIC_2 | $BINARY keys add wallet2 --home $CHAIN_DIR/$CHAINID_2 --recover --keyring-backend=test
echo $WALLET_MNEMONIC_3 | $BINARY keys add wallet3 --home $CHAIN_DIR/$CHAINID_2 --recover --keyring-backend=test
echo $WALLET_MNEMONIC_4 | $BINARY keys add wallet4 --home $CHAIN_DIR/$CHAINID_2 --recover --keyring-backend=test
echo $WALLET_MNEMONIC_5 | $BINARY keys add wallet5 --home $CHAIN_DIR/$CHAINID_2 --recover --keyring-backend=test
echo $WALLET_MNEMONIC_6 | $BINARY keys add wallet6 --home $CHAIN_DIR/$CHAINID_2 --recover --keyring-backend=test
echo $WALLET_MNEMONIC_7 | $BINARY keys add wallet7 --home $CHAIN_DIR/$CHAINID_2 --recover --keyring-backend=test
echo $WALLET_MNEMONIC_8 | $BINARY keys add wallet8 --home $CHAIN_DIR/$CHAINID_2 --recover --keyring-backend=test
echo $WALLET_MNEMONIC_9 | $BINARY keys add wallet9 --home $CHAIN_DIR/$CHAINID_2 --recover --keyring-backend=test
echo $WALLET_MNEMONIC_10 | $BINARY keys add wallet10 --home $CHAIN_DIR/$CHAINID_2 --recover --keyring-backend=test
echo $WALLET_MNEMONIC_11 | $BINARY keys add wallet11 --home $CHAIN_DIR/$CHAINID_2 --recover --keyring-backend=test
echo $WALLET_MNEMONIC_12 | $BINARY keys add wallet12 --home $CHAIN_DIR/$CHAINID_2 --recover --keyring-backend=test
echo $WALLET_MNEMONIC_13 | $BINARY keys add wallet13 --home $CHAIN_DIR/$CHAINID_2 --recover --keyring-backend=test
echo $WALLET_MNEMONIC_14 | $BINARY keys add wallet14 --home $CHAIN_DIR/$CHAINID_2 --recover --keyring-backend=test
echo $WALLET_MNEMONIC_15 | $BINARY keys add wallet15 --home $CHAIN_DIR/$CHAINID_2 --recover --keyring-backend=test

## Special Wallets
VAL1_ADDR=$($BINARY keys show val1 --home $CHAIN_DIR/$CHAINID_1 --keyring-backend test -a)
VAL2_ADDR=$($BINARY keys show val2 --home $CHAIN_DIR/$CHAINID_2 --keyring-backend test -a)
RLY1_ADDR=$($BINARY keys show rly1 --home $CHAIN_DIR/$CHAINID_1 --keyring-backend test -a)
RLY2_ADDR=$($BINARY keys show rly2 --home $CHAIN_DIR/$CHAINID_2 --keyring-backend test -a)

## Other wallets with funds at genesis
WALLET1_ADDR=$($BINARY keys show wallet1 --home $CHAIN_DIR/$CHAINID_1 --keyring-backend test -a)
WALLET2_ADDR=$($BINARY keys show wallet2 --home $CHAIN_DIR/$CHAINID_1 --keyring-backend test -a)
WALLET3_ADDR=$($BINARY keys show wallet3 --home $CHAIN_DIR/$CHAINID_1 --keyring-backend test -a)
WALLET4_ADDR=$($BINARY keys show wallet4 --home $CHAIN_DIR/$CHAINID_1 --keyring-backend test -a)
WALLET5_ADDR=$($BINARY keys show wallet5 --home $CHAIN_DIR/$CHAINID_1 --keyring-backend test -a)
WALLET6_ADDR=$($BINARY keys show wallet6 --home $CHAIN_DIR/$CHAINID_1 --keyring-backend test -a)
WALLET7_ADDR=$($BINARY keys show wallet7 --home $CHAIN_DIR/$CHAINID_1 --keyring-backend test -a)
WALLET8_ADDR=$($BINARY keys show wallet8 --home $CHAIN_DIR/$CHAINID_1 --keyring-backend test -a)
WALLET9_ADDR=$($BINARY keys show wallet9 --home $CHAIN_DIR/$CHAINID_1 --keyring-backend test -a)
WALLET10_ADDR=$($BINARY keys show wallet10 --home $CHAIN_DIR/$CHAINID_1 --keyring-backend test -a)
WALLET11_ADDR=$($BINARY keys show wallet11 --home $CHAIN_DIR/$CHAINID_1 --keyring-backend test -a)
WALLET12_ADDR=$($BINARY keys show wallet12 --home $CHAIN_DIR/$CHAINID_1 --keyring-backend test -a)
WALLET13_ADDR=$($BINARY keys show wallet13 --home $CHAIN_DIR/$CHAINID_1 --keyring-backend test -a)
WALLET14_ADDR=$($BINARY keys show wallet14 --home $CHAIN_DIR/$CHAINID_1 --keyring-backend test -a)
WALLET15_ADDR=$($BINARY keys show wallet15 --home $CHAIN_DIR/$CHAINID_1 --keyring-backend test -a)


$BINARY genesis add-genesis-account $VAL1_ADDR 1000000000000uluna --home $CHAIN_DIR/$CHAINID_1
$BINARY genesis add-genesis-account $VAL2_ADDR 1000000000000uluna --home $CHAIN_DIR/$CHAINID_2
$BINARY genesis add-genesis-account $RLY1_ADDR 1000000000000uluna --home $CHAIN_DIR/$CHAINID_1
$BINARY genesis add-genesis-account $RLY2_ADDR 1000000000000uluna --home $CHAIN_DIR/$CHAINID_2


$BINARY genesis add-genesis-account $WALLET1_ADDR 1000000000000uluna --home $CHAIN_DIR/$CHAINID_1
$BINARY genesis add-genesis-account $WALLET2_ADDR 1000000000000uluna --home $CHAIN_DIR/$CHAINID_1
$BINARY genesis add-genesis-account $WALLET3_ADDR 1000000000000uluna --vesting-amount 10000000000uluna --vesting-start-time $(date +%s) --vesting-end-time $(($(date '+%s') + 100000023)) --home $CHAIN_DIR/$CHAINID_1
$BINARY genesis add-genesis-account $WALLET4_ADDR 1000000000000uluna --vesting-amount 10000000000uluna --vesting-start-time $(date +%s) --vesting-end-time $(($(date '+%s') + 100000023)) --home $CHAIN_DIR/$CHAINID_1
$BINARY genesis add-genesis-account $WALLET5_ADDR 1000000000000uluna --home $CHAIN_DIR/$CHAINID_1
$BINARY genesis add-genesis-account $WALLET6_ADDR 1000000000000uluna --home $CHAIN_DIR/$CHAINID_1
$BINARY genesis add-genesis-account $WALLET7_ADDR 1000000000000uluna --home $CHAIN_DIR/$CHAINID_1
$BINARY genesis add-genesis-account $WALLET8_ADDR 1000000000000uluna --home $CHAIN_DIR/$CHAINID_1
$BINARY genesis add-genesis-account $WALLET9_ADDR 1000000000000uluna --home $CHAIN_DIR/$CHAINID_1
$BINARY genesis add-genesis-account $WALLET10_ADDR 1000000000000uluna --home $CHAIN_DIR/$CHAINID_1
$BINARY genesis add-genesis-account $WALLET11_ADDR 1000000000000uluna --home $CHAIN_DIR/$CHAINID_1
$BINARY genesis add-genesis-account $WALLET12_ADDR 1000000000000uluna --home $CHAIN_DIR/$CHAINID_1
$BINARY genesis add-genesis-account $WALLET13_ADDR 1000000000000uluna --home $CHAIN_DIR/$CHAINID_1
$BINARY genesis add-genesis-account $WALLET14_ADDR 1000000000000uluna --home $CHAIN_DIR/$CHAINID_1
$BINARY genesis add-genesis-account $WALLET15_ADDR 1000000000000uluna --home $CHAIN_DIR/$CHAINID_1

$BINARY genesis add-genesis-account $WALLET1_ADDR 1000000000000uluna --home $CHAIN_DIR/$CHAINID_2
$BINARY genesis add-genesis-account $WALLET2_ADDR 1000000000000uluna --home $CHAIN_DIR/$CHAINID_2
$BINARY genesis add-genesis-account $WALLET3_ADDR 1000000000000uluna --vesting-amount 10000000000uluna --vesting-start-time $(date +%s) --vesting-end-time $(($(date '+%s') + 100000023)) --home $CHAIN_DIR/$CHAINID_2
$BINARY genesis add-genesis-account $WALLET4_ADDR 1000000000000uluna --vesting-amount 10000000000uluna --vesting-start-time $(date +%s) --vesting-end-time $(($(date '+%s') + 100000023)) --home $CHAIN_DIR/$CHAINID_2
$BINARY genesis add-genesis-account $WALLET5_ADDR 1000000000000uluna --home $CHAIN_DIR/$CHAINID_2
$BINARY genesis add-genesis-account $WALLET6_ADDR 1000000000000uluna --home $CHAIN_DIR/$CHAINID_2
$BINARY genesis add-genesis-account $WALLET7_ADDR 1000000000000uluna --home $CHAIN_DIR/$CHAINID_2
$BINARY genesis add-genesis-account $WALLET8_ADDR 1000000000000uluna --home $CHAIN_DIR/$CHAINID_2
$BINARY genesis add-genesis-account $WALLET9_ADDR 1000000000000uluna --home $CHAIN_DIR/$CHAINID_2
$BINARY genesis add-genesis-account $WALLET10_ADDR 1000000000000uluna --home $CHAIN_DIR/$CHAINID_2
$BINARY genesis add-genesis-account $WALLET11_ADDR 1000000000000uluna --home $CHAIN_DIR/$CHAINID_2
$BINARY genesis add-genesis-account $WALLET12_ADDR 1000000000000uluna --home $CHAIN_DIR/$CHAINID_2
$BINARY genesis add-genesis-account $WALLET13_ADDR 1000000000000uluna --home $CHAIN_DIR/$CHAINID_2
$BINARY genesis add-genesis-account $WALLET14_ADDR 1000000000000uluna --home $CHAIN_DIR/$CHAINID_2
$BINARY genesis add-genesis-account $WALLET15_ADDR 1000000000000uluna --home $CHAIN_DIR/$CHAINID_2

echo "Creating and collecting gentx..."
$BINARY genesis gentx val1 7000000000uluna --home $CHAIN_DIR/$CHAINID_1 --chain-id $CHAINID_1 --keyring-backend test
$BINARY genesis gentx val2 7000000000uluna --home $CHAIN_DIR/$CHAINID_2 --chain-id $CHAINID_2 --keyring-backend test
$BINARY genesis collect-gentxs --home $CHAIN_DIR/$CHAINID_1 &> /dev/null
$BINARY genesis collect-gentxs --home $CHAIN_DIR/$CHAINID_2 &> /dev/null

echo "Changing defaults and ports in app.toml and config.toml files..."
sed -i -e 's#"tcp://0.0.0.0:26656"#"tcp://localhost:'"$P2PPORT_1"'"#g' $CHAIN_DIR/$CHAINID_1/config/config.toml
sed -i -e 's#"tcp://127.0.0.1:26657"#"tcp://localhost:'"$RPCPORT_1"'"#g' $CHAIN_DIR/$CHAINID_1/config/config.toml
sed -i -e 's#"tcp://localhost:26657"#"tcp://localhost:'"$RPCPORT_1"'"#g' $CHAIN_DIR/$CHAINID_1/config/client.toml
sed -i -e 's/timeout_commit = "5s"/timeout_commit = "1s"/g' $CHAIN_DIR/$CHAINID_1/config/config.toml
sed -i -e 's/timeout_propose = "3s"/timeout_propose = "1s"/g' $CHAIN_DIR/$CHAINID_1/config/config.toml
sed -i -e 's/index_all_keys = false/index_all_keys = true/g' $CHAIN_DIR/$CHAINID_1/config/config.toml
sed -i -e 's/enable = false/enable = true/g' $CHAIN_DIR/$CHAINID_1/config/app.toml
sed -i -e 's/swagger = false/swagger = true/g' $CHAIN_DIR/$CHAINID_1/config/app.toml
sed -i -e 's#"tcp://localhost:1317"#"tcp://localhost:'"$RESTPORT_1"'"#g' $CHAIN_DIR/$CHAINID_1/config/app.toml
sed -i -e 's/streamers = \[\]/streamers = \["fastquery"\]/g' $CHAIN_DIR/$CHAINID_1/config/app.toml
sed -i '/# Enable defines if the Rosetta API server should be enabled\./ {N; s/enable = true/enable = false/}' $CHAIN_DIR/$CHAINID_1/config/app.toml


sed -i -e 's#"tcp://0.0.0.0:26656"#"tcp://localhost:'"$P2PPORT_2"'"#g' $CHAIN_DIR/$CHAINID_2/config/config.toml
sed -i -e 's#"tcp://127.0.0.1:26657"#"tcp://localhost:'"$RPCPORT_2"'"#g' $CHAIN_DIR/$CHAINID_2/config/config.toml
sed -i -e 's#"tcp://localhost:26657"#"tcp://localhost:'"$RPCPORT_2"'"#g' $CHAIN_DIR/$CHAINID_2/config/client.toml
sed -i -e 's/timeout_commit = "5s"/timeout_commit = "1s"/g' $CHAIN_DIR/$CHAINID_2/config/config.toml
sed -i -e 's/timeout_propose = "3s"/timeout_propose = "1s"/g' $CHAIN_DIR/$CHAINID_2/config/config.toml
sed -i -e 's/index_all_keys = false/index_all_keys = true/g' $CHAIN_DIR/$CHAINID_2/config/config.toml
sed -i -e 's/enable = false/enable = true/g' $CHAIN_DIR/$CHAINID_2/config/app.toml
sed -i -e 's/swagger = false/swagger = true/g' $CHAIN_DIR/$CHAINID_2/config/app.toml
sed -i -e 's#"tcp://localhost:1317"#"tcp://localhost:'"$RESTPORT_2"'"#g' $CHAIN_DIR/$CHAINID_2/config/app.toml
sed -i -e 's/streamers = \[\]/streamers = \["fastquery"\]/g' $CHAIN_DIR/$CHAINID_2/config/app.toml
sed -i '/# Enable defines if the Rosetta API server should be enabled\./ {N; s/enable = true/enable = false/}' $CHAIN_DIR/$CHAINID_2/config/app.toml

echo "Chaning genesis.json..."
sed -i -e 's/"voting_period": "172800s"/"voting_period": "2s"/g' $CHAIN_DIR/$CHAINID_1/config/genesis.json
sed -i -e 's/"voting_period": "172800s"/"voting_period": "2s"/g' $CHAIN_DIR/$CHAINID_2/config/genesis.json
sed -i -e 's/"reward_delay_time": "604800s"/"reward_delay_time": "0s"/g' $CHAIN_DIR/$CHAINID_1/config/genesis.json
sed -i -e 's/"reward_delay_time": "604800s"/"reward_delay_time": "0s"/g' $CHAIN_DIR/$CHAINID_2/config/genesis.json

exit

echo "Starting $CHAINID_1 in $CHAIN_DIR..."
echo "Creating log file at $CHAIN_DIR/$CHAINID_1.log"
$BINARY start --log_level trace --log_format json --home $CHAIN_DIR/$CHAINID_1 --pruning=nothing --grpc.address="0.0.0.0:$GRPCPORT_1" --grpc-web.address="0.0.0.0:$GRPCWEB_1" > $CHAIN_DIR/$CHAINID_1.log 2>&1 &

echo "Starting $CHAINID_2 in $CHAIN_DIR..."
echo "Creating log file at $CHAIN_DIR/$CHAINID_2.log"
$BINARY start --log_level trace --log_format json --home $CHAIN_DIR/$CHAINID_2 --pruning=nothing --grpc.address="0.0.0.0:$GRPCPORT_2" --grpc-web.address="0.0.0.0:$GRPCWEB_2" > $CHAIN_DIR/$CHAINID_2.log 2>&1 &
