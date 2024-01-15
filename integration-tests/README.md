# Terra Core Integration Tests

This project is meant to increase the success ratio of new core releases, improve reliability and features for [@terra-money/feather.js](https://github.com/terra-money/feather.js). This tests are written using TypeScript with [jest](https://jestjs.io/) and tries to improve the coverage by asserting as many outputs as possible.

### Development

This set of tests must run out of the box in Linux-based systems installing [GoLang 1.20](https://go.dev/), [jq](https://stedolan.github.io/jq/) and [screen](https://www.geeksforgeeks.org/screen-command-in-linux-with-examples/). The relayer used in the tests is [go relayer](https://github.com/cosmos/relayer). Keep in mind that the data is not wiped out each time a new test is executed.


Folders structure:
```sh
├── jest.config.js
├── tsconfig.json
├── package.json
├── package-lock.json
├── README.md
└── src
    ├── setup                 # Scripts to start the two networks and relayers
    ├── contracts             # WASM Contracts to be used in the tests.
    ├── helpers               # Functions to improve code readability and avoid duplications.
    │   ├── const.ts
    │   ├── lcd.connection.ts
    │   └── mnemonics.ts
    └── modules               # Tests splited by module
 
```