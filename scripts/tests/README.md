# Tests

This folder contains the integration tests that should run successfully each time a new core release is created.

At the moment, there are four tests defined that can be run with `make integration-test-all`. The breakdown of each test is as follows:

- [init-test-framework](./start.sh): build the core and spin up two nodes with their own genesis event and some accounts preloaded with funds.
- [test-relayer](./relayer/): connect the two blockchains with a relayer opening a channel between both of them.
    - [test-ica](./ica/delegate.sh): using the relayer, this test creates an interchain account and delegates funds using an interchain message.
    - [test-ibc-hooks](./ibc-hooks/increment.sh): deploys a slightly modified [counter contract](./ibc-hooks/counter/) to chain test-2 and submits two requests from chain test-1 to validate that the contract received the funds and executed the wasm code correctly.
- remove-ica-data: removes the data and kills the process when the integration tests are completed.

## Development process

This set of tests must run out of the box in Linux-based systems installing [GoLang 1.20](https://go.dev/), [jq](https://stedolan.github.io/jq/), [screen](https://www.geeksforgeeks.org/screen-command-in-linux-with-examples/) and [rly](https://github.com/cosmos/relayer).