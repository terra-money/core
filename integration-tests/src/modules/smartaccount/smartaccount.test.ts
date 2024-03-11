import { Coins, MsgInstantiateContract, MsgSend, MsgStoreCode, SimplePublicKey } from "@terra-money/feather.js";
import { AuthorizationMsg, Initialization, MsgCreateSmartAccount, MsgUpdateAuthorization } from "@terra-money/feather.js/dist/core/smartaccount";
import fs from "fs";
import path from 'path';
import { blockInclusion, getLCDClient, getMnemonics } from "../../helpers";

describe("Smartaccount Module (https://github.com/terra-money/core/tree/release/v2.6/x/smartaccount) ", () => {
    // Prepare environment clients, accounts and wallets
    const LCD = getLCDClient();
    const accounts = getMnemonics();
    const wallet = LCD.chain1.wallet(accounts.mnemonic5);
    const controlledAccountAddress = accounts.mnemonic5.accAddress("terra");
    
    const controller = LCD.chain1.wallet(accounts.mnemonic4);
    const pubkey = accounts.mnemonic4.publicKey;
    expect(pubkey).toBeDefined();

    // TODO: convert pubkey to base64 string similar to golang pubkey.Bytes()
    const pubkeybb = pubkey as SimplePublicKey;
    const pubkeyStr = pubkeybb.key;
    // AsCe1GUUuW2cT63a35JRpGYaJ6/xIZXvrZRfRGsyxIhK
    const initMsg =  Initialization.fromData({
        account: controlledAccountAddress,
        msg: pubkeyStr,
    });
    
    const deployer = LCD.chain1.wallet(accounts.tokenFactoryMnemonic);
    const deployerAddress = accounts.tokenFactoryMnemonic.accAddress("terra");

    let authContractAddress = "terra14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9ssrc8au";

    test('Deploy smart account auth contract and initialize priv key for wallet', async () => {
        try {
            let tx = await deployer.createAndSignTx({
                msgs: [new MsgStoreCode(
                    deployerAddress,
                    fs.readFileSync(path.join(__dirname, "/../../../../x/smartaccount/test_helpers/test_data/smart_auth_contract.wasm")).toString("base64"),
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
                "Smart auth contract " + Math.random(),
            );

            tx = await deployer.createAndSignTx({
                msgs: [msgInstantiateContract],
                chainID: "test-1",
            });
            result = await LCD.chain1.tx.broadcastSync(tx, "test-1");
            await blockInclusion();
            txResult = await LCD.chain1.tx.txInfo(result.txhash, "test-1") as any;
            authContractAddress = txResult.logs[0].events[4].attributes[0].value;
            expect(authContractAddress).toBeDefined();
        } catch(e: any) {
            console.log(e)
            expect(e).toBeUndefined();
        }
    });

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

    test('Give smart account control to controller', async () => {
        try {
            // give control to controller
            const authMsg = new AuthorizationMsg(authContractAddress, initMsg);
            console.log(authMsg.toData())
            let tx = await wallet.createAndSignTx({
                msgs: [new MsgUpdateAuthorization(
                    controlledAccountAddress,
                    false,
                    [authMsg],
                )],
                chainID: 'test-1',
                gas: '400000',
            });
            await LCD.chain1.tx.broadcastSync(tx, "test-1");
            await blockInclusion();

            // check if update authorization was successful
            let setting = await LCD.chain1.smartaccount.setting(controlledAccountAddress);
            expect(setting.toData())
                .toEqual({
                    owner: controlledAccountAddress,
                    authorization: [authMsg.toData()],
                    post_transaction: [],
                    pre_transaction: [],
                    fallback: false,
                });
        } catch (e:any) {
            console.log(e)
            expect(e).toBeUndefined();
        }
    });

    test.only('Only controller should be able to sign', async () => {
        try {
            // signing with the controlledAccountAddress should now fail 
            let tx = await wallet.createAndSignTx({
                msgs: [
                    new MsgSend(
                        controlledAccountAddress,
                        controlledAccountAddress,
                        Coins.fromString("1uluna"),
                    ),
                ],
                chainID: "test-1",
                gas: '400000',
            });
            let result = await LCD.chain1.tx.broadcastSync(tx, "test-1");
            expect(result.raw_log).toEqual("authorization failed: unauthorized");

            // signing with the controller should now succeed
            tx = await controller.createAndSignTx({
                msgs: [
                    new MsgSend(
                        controlledAccountAddress,
                        deployerAddress,
                        Coins.fromString("1uluna"),
                    ),
                ],
                chainID: "test-1",
                gas: '400000',
            });
            const deployerBalanceBefore = await LCD.chain1.bank.balance(deployerAddress);
            result = await LCD.chain1.tx.broadcastSync(tx, "test-1");
            await blockInclusion();
            const deployerBalanceAfter = await LCD.chain1.bank.balance(deployerAddress);
            const deltaBalance = deployerBalanceAfter[0].sub(deployerBalanceBefore[0]);
            expect(deltaBalance.toString()).toEqual("1uluna");
        } catch (e:any) {
            console.log(e)
            expect(e).toBeUndefined();
        }
    });
});