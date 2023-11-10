import { Coins,  MsgInstantiateContract, MsgStoreCode } from "@terra-money/feather.js";
import { getMnemonics, getLCDClient, blockInclusion } from "../../helpers";
import fs from "fs";
import path from 'path';

describe("Feeshare Module (https://github.com/terra-money/core/tree/release/v2.7/x/feeshare) ", () => {
    // Prepare environment clients, accounts and wallets
    const LCD = getLCDClient();
    const accounts = getMnemonics();
    const wallet = LCD.chain1.wallet(accounts.tokenFactoryMnemonic);
    const tokenFactoryWalletAddr = accounts.tokenFactoryMnemonic.accAddress("terra");
    let contractAddress: string;

    // Reat the reflect contract, store on chain, 
    // instantiate to be used in the following tests
    // and finally save the contract address.
    beforeAll(async () => {
        let tx = await wallet.createAndSignTx({
            msgs: [new MsgStoreCode(
                tokenFactoryWalletAddr,
                fs.readFileSync(path.join(__dirname, "/../../contracts/no100.wasm")).toString("base64"),
            )],
            chainID: "test-1",
        });

        let result = await LCD.chain1.tx.broadcastSync(tx, "test-1");
        await blockInclusion();
        let txResult = await LCD.chain1.tx.txInfo(result.txhash, "test-1") as any;
        let codeId = Number(txResult.logs[0].events[1].attributes[1].value);
        expect(codeId).toBeDefined();

        const msgInstantiateContract = new MsgInstantiateContract(
            tokenFactoryWalletAddr,
            tokenFactoryWalletAddr,
            codeId,
            {},
            Coins.fromString("1uluna"),
            "no100 contract " + Math.random(),
        );

        tx = await wallet.createAndSignTx({
            msgs: [msgInstantiateContract],
            chainID: "test-1",
        });
        result = await LCD.chain1.tx.broadcastSync(tx, "test-1");
        await blockInclusion();
        txResult = await LCD.chain1.tx.txInfo(result.txhash, "test-1") as any;
        contractAddress = txResult.logs[0].events[4].attributes[0].value;
        expect(contractAddress).toBeDefined();
    });

    test('Must contain the correct module params', async () => {
        const moduleParams = await LCD.chain1.feeshare.params("test-1");

        expect(moduleParams)
            .toMatchObject({
                "params": {
                    "allowed_denoms": [],
                    "developer_shares": "0.500000000000000000",
                    "enable_fee_share": true,
                },
            });
    });

    test('Must query all endpoints before creating a denom', async () => {
        // // Register feeshare
        // let tx = await wallet.createAndSignTx({
        //     msgs: [new MsgCreateDenom(
        //         contractAddress,
        //         subdenom,
        //     )],
        //     chainID: "test-1",
        // });
        // let result = await LCD.chain1.tx.broadcastSync(tx, "test-1");
        // await blockInclusion();
        // let txResult = await LCD.chain1.tx.txInfo(result.txhash, "test-1") as any;
        // console.log(txResult.logs)
        // console.log(randomAccountAddress)
        // expect(txResult.logs).toBeDefined();
    });
});