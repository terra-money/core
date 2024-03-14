import { Coin, Coins, MsgInstantiateContract, MsgSend, MsgStoreCode } from "@terra-money/feather.js";
import { MsgCreateSmartAccount, MsgUpdateTransactionHooks } from "@terra-money/feather.js/dist/core/smartaccount";
import fs from "fs";
import path from 'path';
import { blockInclusion, getLCDClient, getMnemonics } from "../../helpers";

describe("Smartaccount Module (https://github.com/terra-money/core/tree/release/v2.6/x/smartaccount) ", () => {
    // Prepare environment clients, accounts and wallets
    const LCD = getLCDClient();
    const accounts = getMnemonics();
    const wallet = LCD.chain1.wallet(accounts.smartaccPreTxMnemonic);
    const smartaccAddress = accounts.smartaccPreTxMnemonic.accAddress("terra");
    const receiver = accounts.smartaccControllerMnemonic.accAddress("terra")
    
    const deployer = LCD.chain1.wallet(accounts.tokenFactoryMnemonic);
    const deployerAddress = accounts.tokenFactoryMnemonic.accAddress("terra");

    let limitContractAddress: string;

    test('Create new smart account', async () => {
        try {
            // create the smartaccount
            let tx = await wallet.createAndSignTx({
                msgs: [new MsgCreateSmartAccount(
                    smartaccAddress
                )],
                chainID: 'test-1',
                gas: '400000',
            });
            await LCD.chain1.tx.broadcastSync(tx, "test-1");
            await blockInclusion();
            // Query smartaccount setting for the smart waccount
            let setting = await LCD.chain1.smartaccount.setting(smartaccAddress);
            expect(setting.toData())
                .toEqual({
                    owner: smartaccAddress,
                    authorization: [],
                    post_transaction: [],
                    pre_transaction: [],
                    fallback: true,
                });
        } catch (e:any) {
            console.log(e);
            expect(e).toBeUndefined();
        }
    });

    test('Deploy smart account limit contract', async () => {
        try {
            let tx = await deployer.createAndSignTx({
                msgs: [new MsgStoreCode(
                    deployerAddress,
                    fs.readFileSync(path.join(__dirname, "/../../../../x/smartaccount/test_helpers/test_data/limit_min_coins_hooks.wasm")).toString("base64"),
                )],
                chainID: "test-1",
            });

            let result = await LCD.chain1.tx.broadcastSync(tx, "test-1");
            await blockInclusion();
            let txResult = await LCD.chain1.tx.txInfo(result.txhash, "test-1") as any;
            let codeId = Number(txResult.logs[0].events[1].attributes[1].value);
            expect(codeId).toBeDefined();

            const msgInstantiateContract = new MsgInstantiateContract(
                deployerAddress,
                deployerAddress,
                codeId,
                {},
                Coins.fromString("1uluna"),
                "limit contract " + Math.random(),
            );

            tx = await deployer.createAndSignTx({
                msgs: [msgInstantiateContract],
                chainID: "test-1",
            });
            result = await LCD.chain1.tx.broadcastSync(tx, "test-1");
            await blockInclusion();
            txResult = await LCD.chain1.tx.txInfo(result.txhash, "test-1") as any;
            limitContractAddress = txResult.logs[0].events[4].attributes[0].value;
            expect(limitContractAddress).toBeDefined();
        } catch(e: any) {
            console.log(e)
            expect(e).toBeUndefined();
        }
    });

    test('Update post tx hooks', async () => {
        try {
            // signing with the controlledAccountAddress should now fail 
            let tx = await wallet.createAndSignTx({
                msgs: [
                    new MsgUpdateTransactionHooks(
                        smartaccAddress,
                        [],
                        [limitContractAddress],
                    ),
                ],
                chainID: "test-1",
                gas: '400000',
            });
            await LCD.chain1.tx.broadcastSync(tx, "test-1");
            await blockInclusion();

            // check if update authorization was successful
            let setting = await LCD.chain1.smartaccount.setting(smartaccAddress);
            expect(setting.postTransaction).toEqual([limitContractAddress]);
        } catch (e:any) {
            console.log(e)
            expect(e).toBeUndefined();
        }
    });

    test('Transaction should pass for sending below limit', async () => {
        try {
            let setting = await LCD.chain1.smartaccount.setting(smartaccAddress);
            expect(setting.postTransaction).toEqual([limitContractAddress]);

            const balance = await LCD.chain1.bank.balance(smartaccAddress);
            const coinsToSend = balance[0].sub(Coins.fromString("1000000uluna"));

            let tx = await wallet.createAndSignTx({
                msgs: [
                    new MsgSend(
                        smartaccAddress,
                        receiver,
                        coinsToSend,
                    ),
                ],
                chainID: "test-1",
                gas: '400000',
            });
            // expect.assertions(1);
            await LCD.chain1.tx.broadcastSync(tx, "test-1");
            await blockInclusion();

            const balanceAfter = await LCD.chain1.bank.balance(smartaccAddress);
            const coinAfter = balanceAfter[0].find((coin: Coin) => coin.denom === "uluna");
            expect(coinAfter).toBeDefined();
            expect(coinAfter!.amount.toNumber()).toBeLessThan(1000000);
            expect(coinAfter!.amount.toNumber()).toBeGreaterThan(900000);
        } catch (e:any) {
            console.log(e)
            expect(e).toBeUndefined();
        }
    });
    
    test.only('Transaction should fail simulation for sending over limit', async () => {
        try {
            // let setting = await LCD.chain1.smartaccount.setting(smartaccAddress);
            // expect(setting.postTransaction).toEqual([limitContractAddress]);

            // should have 940000uluna at this point 60000
            const balance = await LCD.chain1.bank.balance(smartaccAddress);
            const coinsToSend = balance[0].sub(Coins.fromString("60000uluna"));

            let tx = await wallet.createAndSignTx({
                msgs: [
                    new MsgSend(
                        smartaccAddress,
                        receiver,
                        coinsToSend,
                    ),
                ],
                chainID: "test-1",
                gas: '400000',
            });
            await LCD.chain1.tx.simulateTx(tx, "test-1");            
        } catch (e:any) {
            // TODO: simulate transaction will fail but e is an object
            // find a way to make assertion below pass
            console.log(e)
            expect(e.toString()).toContain("Failed post transaction process: Account balance is less than 1000: execute wasm contract failed");
        }
    });

    // test('Transaction should fail simulation for sending over limit', async () => {
    //     try {
    //         // let setting = await LCD.chain1.smartaccount.setting(smartaccAddress);
    //         // expect(setting.postTransaction).toEqual([limitContractAddress]);

    //         // should have 940000uluna at this point 60000
    //         const balance = await LCD.chain1.bank.balance(smartaccAddress);
    //         const coinsToSend = balance[0].sub(Coins.fromString("61000uluna"));

    //         let tx = await wallet.createAndSignTx({
    //             msgs: [
    //                 new MsgSend(
    //                     smartaccAddress,
    //                     receiver,
    //                     coinsToSend,
    //                 ),
    //             ],
    //             chainID: "test-1",
    //             gas: '400000',
    //         });
    //         // simulate should not fail
    //         await LCD.chain1.tx.simulateTx(tx, "test-1");

    //         console.log("asdasd")
    //         const result = await LCD.chain1.tx.broadcastSync(tx, "test-1");
    //         console.log(result)
    //         await blockInclusion();

    //         // check that MsgSend failed
    //         const balanceAfter = await LCD.chain1.bank.balance(smartaccAddress);
    //         const coinAfter = balanceAfter[0].find((coin: Coin) => coin.denom === "uluna");
    //         expect(coinAfter).toBeDefined();
    //         expect(coinAfter!.amount.toNumber()).toBeGreaterThan(60000);
    //     } catch (e:any) {
    //         console.log(e)
    //         expect(e).toBeUndefined();
    //     }
    // });
});