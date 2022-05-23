<p>&nbsp;</p>
<p align="center">

<img src="core_logo.svg" width=500>

</p>

<p align="center">
The full-node software implementation of the Terra blockchain.<br/><br/>

<a href="https://codecov.io/gh/terra-money/core">
    <img src="https://codecov.io/gh/terra-money/core/branch/main/graph/badge.svg">
</a>
<a href="https://goreportcard.com/report/github.com/terra-money/core">
    <img src="https://goreportcard.com/badge/github.com/terra-money/core">
</a>

</p>

<p align="center">
  <a href="https://docs.terra.money/"><strong>Explore the Docs »</strong></a>
  <br />
  <br/>
  <a href="https://docs.terra.money/docs/develop/module-specifications/README.html">Terra Core reference</a>
  ·
  <a href="https://pkg.go.dev/github.com/terra-money/core?tab=subdirectories">Go API</a>
  ·
  <a href="https://lcd.terra.dev/swagger/#/">Rest API</a>
  ·
  <a href="https://finder.terra.money/">Finder</a>
  ·
  <a href="https://station.terra.money/">Station</a>
</p>

<br/>

## Terra migration guides

Visit the [migration guide](https://migrate.terra.money) to learn how to migrate from Terra Classic to the new Terra blockchain.

## Table of Contents <!-- omit in toc -->

- [What is Terra?](#what-is-terra)
- [Installation](#installation)
  - [From Binary](#from-binary)
  - [From Source](#from-source)
  - [terrad](#terrad)
- [Node Setup](#node-setup)
  - [Terra node quickstart](#terra-node-quickstart)
  - [Join the mainnet](#join-the-mainnet)
  - [Join a testnet](#join-a-testnet)
  - [Run a local testnet](#run-a-local-testnet)
  - [Run a single node testnet](#run-a-single-node-testnet)
- [Set up a production environment](#set-up-a-production-environment)
  - [Increase maximum open files](#increase-maximum-open-files)
  - [Create a dedicated user](#create-a-dedicated-user)
  - [Port configuration](#port-configuration)
  - [Run the server as a daemon](#run-the-server-as-a-daemon)
  - [Register terrad as a service](#register-terrad-as-a-service)
  - [Start, stop, or restart service](#start-stop-or-restart-service)
  - [Access logs](#access-logs)
- [Resources](#resources)
- [Community](#community)
- [Contributing](#contributing)
- [License](#license)

## What is Terra?

**[Terra](https://terra.money)** is a public, open-source, proof-of-stake blockchain. **The Terra Core** is the reference implementation of the Terra protocol written in Golang. The Terra Core is powered by the [Cosmos SDK](https://github.com/cosmos/cosmos-sdk) and [Tendermint](https://github.com/tendermint/tendermint) BFT consensus.

## Installation

### From Binary

The easiest way to install the Terra Core is to download a pre-built binary. You can find the latest binaries on the [releases](https://github.com/terra-money/core/releases) page.

### From Source

**Step 1: Install Golang**

Go v1.18+ or higher is required for The Terra Core.

1. Install [Go 1.18+ from the official site](https://go.dev/dl/). Ensure that your `GOPATH` and `GOBIN` environment variables are properly set up by using the following commands:

   For Windows:

   ```sh
   wget <https://golang.org/dl/go1.18.2.linux-amd64.tar.gz>
   sudo tar -C /usr/local -xzf go1.18.2.linux-amd64.tar.gz
   export PATH=$PATH:/usr/local/go/bin
   export PATH=$PATH:$(go env GOPATH)/bin
   ```

   For Mac:

   ```sh
   export PATH=$PATH:/usr/local/go/bin
   export PATH=$PATH:$(go env GOPATH)/bin
   ```

2. Confirm your Go installation by checking the version:

   ```sh
   go version
   ```


**Step 2: Get Terra Core source code**

Clone the Terra Core from the [official repo](https://github.com/terra-money/core/) and check out the `main` branch for the latest stable release.

```bash
git clone https://github.com/terra-money/core/
cd core
git checkout main
```

**Step 3: Build Terra core**

Run the following command to install `terrad` to your `GOPATH` and build the Terra Core. `terrad` is the node daemon and CLI for interacting with a Terra node.

```bash
# COSMOS_BUILD_OPTIONS=rocksdb make install
make install
```

**Step 4: Verify your installation**

Verify your installation with the following command:

```bash
terrad version --long
```

A successful installation will return the following:

```bash
name: terra
server_name: terrad
version: <x.x.x>
commit: <Commit hash>
build_tags: netgo,ledger
go: go version go1.18.2 darwin/amd64
```

## `terrad`

`terrad` is the all-in-one CLI and node daemon for interacting with the Terra blockchain. 

To view various subcommands and their expected arguments, use the following command:

``` sh
$ terrad --help
```


```
Stargate Terra App

Usage:
  terrad [command]

Available Commands:
  add-genesis-account Add a genesis account to genesis.json
  collect-gentxs      Collect genesis txs and output a genesis.json file
  debug               Tool for helping with debugging your application
  export              Export state to JSON
  gentx               Generate a genesis tx carrying a self delegation
  help                Help about any command
  init                Initialize private validator, p2p, genesis, and application configuration files
  keys                Manage your application's keys
  migrate             Migrate genesis to a specified target version
  query               Querying subcommands
  rosetta             spin up a rosetta server
  start               Run the full node
  status              Query remote node for status
  tendermint          Tendermint subcommands
  testnet             Initialize files for a terrad testnet
  tx                  Transactions subcommands
  unsafe-reset-all    Resets the blockchain database, removes address book files, and resets data/priv_validator_state.json to the genesis state
  validate-genesis    validates the genesis file at the default location or at the location passed as an arg
  version             Print the application binary version information

Flags:
  -h, --help                help for terrad
      --home string         directory for config and data (default "/Users/evan/.terra")
      --log_format string   The logging format (json|plain) (default "plain")
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic) (default "info")
      --trace               print out full stack trace on errors

Use "terrad [command] --help" for more information about a command.
```

Visit the [terrad documentation page](https://docs.terra.money/docs/develop/how-to/terrad/README.html) for more info on usage. 


## Node Setup

Once you have `terrad` installed, you will need to set up your node to be part of the network.

### Join the mainnet

The following requirements are recommended for running a mainnet node:

- Four or more CPU cores
- At least 32 GB of memory
- At least 300 mbps of network bandwidth
- At least 2 TB NVME SSD
- A Linux distribution

### Terra node quickstart
```
terrad init nodename
wget -O ~/.terra/config/genesis.json https://cloudflare-ipfs.com/ipfs/QmZAMcdu85Qr8saFuNpL9VaxVqqLGWNAs72RVFhchL9jWs
curl https://network.terra.dev/addrbook.json > ~/.terrad/config/addrbook.json
terrad start
```

### Join a testnet

Several testnets might exist simultaneously. Ensure that your version of `terrad` is compatible with the network you want to join.

To set up a node on the latest testnet, visit [the testnet repo](https://github.com/terra-money/testnet).

#### Run a local testnet

The easiest way to set up a local testing environment is to run [LocalTerra](https://github.com/terra-money/LocalTerra), a zero-configuration complete testing environment. 

### Run a single node testnet

You can also run a local testnet using a single node. On a local testnet, you will be the sole validator signing blocks.

**Step 1: Create network and account**

First, initialize your genesis file to bootstrap your network. Create a name for your local testnet and provide a moniker to refer to your node:

```bash
terrad init --chain-id=<testnet_name> <node_moniker>
```

Next, create a Terra account by running the following command:

```bash
terrad keys add <account_name>
```

**Step 2: Add account to genesis**

Add your account to genesis and set an initial balance to start. Run the following commands to add your account and set the initial balance:

```bash
terrad add-genesis-account $(terrad keys show <account_name> -a) 100000000uluna
terrad gentx <account_name> 10000000uluna --chain-id=<testnet_name>
terrad collect-gentxs
```

**Step 3: Run terrad**

Now you can start your private Terra network:

```bash
terrad start
```

Your `terrad` node will be running a node on `tcp://localhost:26656`, listening for incoming transactions and signing blocks.

## Set up a production environment

**Note**: This guide only covers general settings for a production-level full node. Visit the [Terra validator's guide](https://docs.terra.money/docs/full-node/manage-a-terra-validator/README.html) for more information.

**This guide has been tested against Linux distributions only. To ensure you successfully set up your production environment, consider setting it up on an Linux system.**

### Increase maximum open files

By default, `terrad` can't open more than 1024 files at once.

You can increase this limit by modifying `/etc/security/limits.conf` and raising the `nofile` capability.

```
*                soft    nofile          65535
*                hard    nofile          65535
```

### Create a dedicated user

It is recommended that you run `terrad` as a normal user. Super-user accounts are only recommended during setup to create and modify files.

### Port configuration

`terrad` uses several TCP ports for different purposes.

- `26656`: The default port for the P2P protocol. Use this port to communicate with other nodes. While this port must be open to join a network, it does not have to be open to the public. Validator nodes should configure `persistent_peers` and close this port to the public.

- `26657`: The default port for the RPC protocol. This port is used for querying / sending transactions and must be open to serve queries from `terrad`. **DO NOT** open this port to the public unless you are planning to run a public node.

- `1317`: The default port for [Lite Client Daemon](https://docs.terra.money/docs/develop/how-to/start-lcd.html) (LCD), which can be enabled in `~/.terra/config/app.toml`. The LCD provides an HTTP RESTful API layer to allow applications and services to interact with your `terrad` instance through RPC. Check the [Terra REST API](https://lcd.terra.dev/swagger/#/) for usage examples. Don't open this port unless you need to use the LCD.

- `26660`: The default port for interacting with the [Prometheus](https://prometheus.io) database. You can use Promethues to monitor an environment. This port is closed by default.

### Run the server as a daemon

**Important**:

Keep `terrad` running at all times. The simplest solution is to register `terrad` as a `systemd` service so that it automatically starts after system reboots and other events.


### Register terrad as a service

First, create a service definition file in `/etc/systemd/system`.

**Sample file: `/etc/systemd/system/terrad.service`**

```
[Unit]
Description=Terra Daemon
After=network.target

[Service]
Type=simple
User=terra
ExecStart=/data/terra/go/bin/terrad start
Restart=on-abort

[Install]
WantedBy=multi-user.target

[Service]
LimitNOFILE=65535
```

Modify the `Service` section from the given sample above to suit your settings.
Note that even if you raised the number of open files for a process, you still need to include `LimitNOFILE`.

After creating a service definition file, you should execute `systemctl daemon-reload`.

### Start, stop, or restart service

Use `systemctl` to control (start, stop, restart)

```bash
# Start
systemctl start terrad
# Stop
systemctl stop terrad
# Restart
systemctl restart terrad
```

### Access logs

```bash
# Entire log
journalctl -t terrad
# Entire log reversed
journalctl -t terrad -r
# Latest and continuous
journalctl -t terrad -f
```

## Resources

Developer Tools:
  - [Terra docs](https://docs.terra.money): Developer documentation.
  - [Faucet](https://faucet.terra.money): Get testnet Luna.
  - [LocalTerra](https://www.github.com/terra-money/LocalTerra): A dockerized local blockchain testnet. 

Developer Forums: 
  - [Terra Developer Discord](https://discord.com/channels/464241079042965516/591812948867940362)
  - [Terra Developer Telegram room](https://t.me/+gCxCPohmVBkyNDRl)

Block Explorer:
  - [Terra Finder](https://finder.terra.money): Terra's basic block explorer.

Wallets:
  - [Terra Station](https://station.terra.money): The official Terra wallet.
  - Terra Station Mobile:
    - [iOS](https://apps.apple.com/us/app/terra-station/id1548434735)
    - [Android](https://play.google.com/store/apps/details?id=money.terra.station&hl=en_US&gl=US)

Research:
  - [Agora](https://agora.terra.money): Research forum

## Community

- [Offical Website](https://terra.money)
- [Discord](https://discord.gg/e29HWwC2Mz)
- [Telegram](https://t.me/terra_announcements)
- [Twitter](https://twitter.com/terra_money)
- [YouTube](https://goo.gl/3G4T1z)

## Contributing

If you are interested in contributing to Terra Core source, please review our [code of conduct](./CODE_OF_CONDUCT.md).

## License

This software is [licensed under the Apache 2.0 license](LICENSE).

© 2022 Terraform Labs, PTE LTD

<p>&nbsp;</p>
<p align="center">
    <a href="https://terra.money/"><img src="https://assets.website-files.com/611153e7af981472d8da199c/61794f2b6b1c7a1cb9444489_symbol-terra-blue.svg" align="center" width=200/></a>
</p>
<div align="center">
</div>


