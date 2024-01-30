#!/bin/bash

echo ""
echo "##################"
echo "# Feemarket Test #"
echo "##################"
echo ""

# Configure predefined mnemonic pharses
BINARY=terrad
CHAIN_DIR=$(pwd)/src/test-data
CHAINID_1=test-1
VAL_1="terra14zvlymmmtgt8gwhknsy83gkfdkeaxfylmx39l8"
RLY_1="terra1a698u5rm2x6y50x5m3q37tnn0k6d4rjpfc8e7h"
WALLET_1="terra1v0eee20gjl68fuk0chyrkch2z7suw2mhg3wkxf"
WALLET_2="terra182ylz0xutzy5qs57nwkxmmqgy2trlfejhuzrfx"
WALLET_3="terra1gyf58rxglrzp343d4wkw7vzlcw6d8knp2qmg0t"
WALLET_4="terra1pdap0ppzxwyn3mz49fsy6utq4sqswp9tztgr48"
WALLET_5="terra120rzk7n6cd2vufkmwrat34adqh0rgca9tkyfe5"
WALLET_6="terra155zkh5akfg8a0cpz5cps7ukw7sztdcxynwe85t"
WALLET_7="terra1p4kcrttuxj9kyyvv5px5ccgwf0yrw74yp7jqm6"
WALLET_8="terra19ekhxctawn8vqvppm6m5nl3htkwym3hcdcka2a"

# Spam transactions
spamtx() {
  for i in {1..15}; do
    $BINARY tx bank send $VAL_1 $WALLET_1 1uluna --from val1 --node http://localhost:16657 --keyring-backend test --fees 30000uluna --home $CHAIN_DIR/$CHAINID_1 -y
    $BINARY tx bank send $WALLET_1 $WALLET_1 1uluna --from wallet1 --node http://localhost:16657 --keyring-backend test --fees 30000uluna --home $CHAIN_DIR/$CHAINID_1 -y
    $BINARY tx bank send $WALLET_2 $WALLET_1 1uluna --from wallet2 --node http://localhost:16657 --keyring-backend test --fees 30000uluna --home $CHAIN_DIR/$CHAINID_1 -y
    $BINARY tx bank send $WALLET_3 $WALLET_1 1uluna --from wallet3 --node http://localhost:16657 --keyring-backend test --fees 30000uluna --home $CHAIN_DIR/$CHAINID_1 -y
    $BINARY tx bank send $WALLET_4 $WALLET_1 1uluna --from wallet4 --node http://localhost:16657 --keyring-backend test --fees 30000uluna --home $CHAIN_DIR/$CHAINID_1 -y
    $BINARY tx bank send $WALLET_5 $WALLET_1 1uluna --from wallet5 --node http://localhost:16657 --keyring-backend test --fees 30000uluna --home $CHAIN_DIR/$CHAINID_1 -y
    $BINARY tx bank send $WALLET_6 $WALLET_1 1uluna --from wallet6 --node http://localhost:16657 --keyring-backend test --fees 30000uluna --home $CHAIN_DIR/$CHAINID_1 -y
    $BINARY tx bank send $WALLET_7 $WALLET_1 1uluna --from wallet7 --node http://localhost:16657 --keyring-backend test --fees 30000uluna --home $CHAIN_DIR/$CHAINID_1 -y
    $BINARY tx bank send $WALLET_8 $WALLET_1 1uluna --from wallet8 --node http://localhost:16657 --keyring-backend test --fees 30000uluna --home $CHAIN_DIR/$CHAINID_1 -y
  done
}

spamtx

echo ""
echo "############################"
echo "# SUCCESS: Feemarket Tested #"
echo "############################"
echo ""
