import { Coin, Coins, MsgDelegate, MsgInstantiateContract, MsgStoreCode, ValAddress } from "@terra-money/feather.js";
import { MsgCreateSmartAccount, MsgUpdateTransactionHooks } from "@terra-money/feather.js/dist/core/smartaccount";
import fs from "fs";
import path from 'path';
import { blockInclusion, getLCDClient, getMnemonics } from "../../helpers";

describe("Smartaccount Module (https://github.com/terra-money/core/tree/release/v2.6/x/smartaccount) ", () => {
    // Prepare environment clients, accounts and wallets
    const LCD = getLCDClient();
    const accounts = getMnemonics();
    const wallet = LCD.chain1.wallet(accounts.mnemonic3);
    const controlledAccountAddress = accounts.mnemonic3.accAddress("terra");
    
    const pubkey = accounts.mnemonic4.publicKey;
    expect(pubkey).toBeDefined();
    
    const deployer = LCD.chain1.wallet(accounts.tokenFactoryMnemonic);
    const deployerAddress = accounts.tokenFactoryMnemonic.accAddress("terra");

    let limitContractAddress: string = "terra1nc5tatafv6eyq7llkr2gv50ff9e22mnf70qgjlv737ktmt4eswrquka9l6"

    test('Create new smart account', async () => {
        try {
            // create the smartaccount
            let tx = await wallet.createAndSignTx({
                msgs: [new MsgCreateSmartAccount(
                    controlledAccountAddress
                )],
                chainID: 'test-1',
                gas: '400000',
            });
            await LCD.chain1.tx.broadcastSync(tx, "test-1");
            await blockInclusion();
            // Query smartaccount setting for the smart waccount
            let setting = await LCD.chain1.smartaccount.setting(controlledAccountAddress);
            expect(setting.toData())
                .toEqual({
                    owner: controlledAccountAddress,
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
                    fs.readFileSync(path.join(__dirname, "/../../../../x/smartaccount/test_helpers/test_data/limit_send_only_hooks.wasm")).toString("base64"),
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

    test('Update pre tx hooks', async () => {
        try {
            // signing with the controlledAccountAddress should now fail 
            let tx = await wallet.createAndSignTx({
                msgs: [
                    new MsgUpdateTransactionHooks(
                        controlledAccountAddress,
                        [limitContractAddress],
                        [],
                    ),
                ],
                chainID: "test-1",
                gas: '400000',
            });
            await LCD.chain1.tx.broadcastSync(tx, "test-1");
            await blockInclusion();

            // check if update authorization was successful
            let setting = await LCD.chain1.smartaccount.setting(controlledAccountAddress);
            expect(setting.preTransaction).toEqual([limitContractAddress]);
        } catch (e:any) {
            console.log(e)
            expect(e).toBeUndefined();
        }
    });

    test.only('Transaction should fail for limit', async () => {
        try {
            // signing with the controlledAccountAddress should now fail 
            let tx = await wallet.createAndSignTx({
                msgs: [
                    new MsgDelegate(
                        controlledAccountAddress,
                        ValAddress.fromAccAddress(controlledAccountAddress, "terra"),
                        Coin.fromString("100000000uluna"),
                    ),
                ],
                chainID: "test-1",
                gas: '400000',
            });
            let result = await LCD.chain1.tx.broadcastSync(tx, "test-1");
            console.log(result)
            expect(result.raw_log).toEqual("Unauthorized message type");
        } catch (e:any) {
            console.log(e)
            expect(e).toBeUndefined();
        }
    });
});