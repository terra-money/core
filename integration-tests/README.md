# Terra Core Integration Tests

This project is meant to increase the success ratio for new core releases, improve reliability and features for [@terra-money/feather.js](https://github.com/terra-money/feather.js). This tests are written using TypeScript with [jest](https://jestjs.io/) and executed in parallel to improve time execution.


### Development

Keep in mind that tests are executed in paralel when using the same account with two different tests it can misslead test results with errors like "account missmatch sequence" when submitting two transactions with the same nonce, missmatching balances, etc... 

Another good practice with this framework is to isolate and assert values within a test considering that the data is not wiped out each time a new test is executed.


Folders structure:
```sh
├── jest.config.js
├── tsconfig.json
├── package.json
├── package-lock.json
├── README.md
└── src
    ├── contracts             # WASM Contracts to be used in the tests.
    │   └── reflect.wasm
    ├── helpers               # Functions to improve code readability and avoid duplications.
    │   ├── const.ts
    │   ├── lcd.connection.ts
    │   └── mnemonics.ts
    └── modules               # Tests splited by module that intent to test.
        ├── auth.test.ts
        ├── feeshare.test.ts
        └── pob.test.ts
 
```