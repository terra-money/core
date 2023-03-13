# Tests

This folder contain the integration tests that should run successfully each time a new core release is created.

At the moment there are four tests defined that can be run with `make integration-test-all`. The brakdown of each tests:

- [init-test-framework](./start.sh): build the core and spinup two nodes with their own genesis event and some accounts preloaded with funds,
- [test-relayer](./relayer/): connect the two blockchains with a relayer openning a channel between both of them,
    - [test-ica](./ica/delegate.sh): using the relayer this test create an interchain account and delegate funds using an interchain message,
    - [test-ibc-hooks](./ibc-hooks/increment.sh): deploy a slightly modified [counter contract](./ibc-hooks/counter/) to chain test-2 and submits two request from chain test-1 to validate the contract received the funds and executed the wasm code correctly,
- remove-ica-data: removes the data and kills the process when the integration tests are completed.

## Development process

This set of tests must run out of the box in Linux based systems installing [GoLang 1.20](https://go.dev/), [jq](https://stedolan.github.io/jq/), [screen](https://www.geeksforgeeks.org/screen-command-in-linux-with-examples/) and [rly](https://github.com/cosmos/relayer).