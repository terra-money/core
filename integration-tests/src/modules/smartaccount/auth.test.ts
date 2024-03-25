import { Coins, MsgInstantiateContract, MsgSend, MsgStoreCode, SimplePublicKey } from "@terra-money/feather.js";
import { AuthorizationMsg, Initialization, MsgCreateSmartAccount, MsgDisableSmartAccount, MsgUpdateAuthorization } from "@terra-money/feather.js/dist/core/smartaccount";
import fs from "fs";
import path from 'path';
import { blockInclusion, getLCDClient, getMnemonics } from "../../helpers";

describe("Smartaccount Module (https://github.com/terra-money/core/tree/release/v2.6/x/smartaccount) ", () => {
    // Prepare environment clients, accounts and wallets
    const LCD = getLCDClient();
    const accounts = getMnemonics();
    const wallet = LCD.chain1.wallet(accounts.smartaccAuthMnemonic);
    const controlledAccountAddress = accounts.smartaccAuthMnemonic.accAddress("terra");
    
    const controller = LCD.chain1.wallet(accounts.smartaccControllerMnemonic);
    const pubkey = accounts.smartaccControllerMnemonic.publicKey;
    expect(pubkey).toBeDefined();

    const pubkeybb = pubkey as SimplePublicKey;
    const pubkeyStr = pubkeybb.key;
    const initMsg =  Initialization.fromData({
        account: controlledAccountAddress,
        msg: pubkeyStr,
    });
    
    const deployer = LCD.chain1.wallet(accounts.tokenFactoryMnemonic);
    const deployerAddress = accounts.tokenFactoryMnemonic.accAddress("terra");

    let authContractAddress: string;

    test('Deploy smart account auth contract', async () => {
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

    test('Only controller should be able to sign', async () => {
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

    test('Disable smart contract', async () => {
        try {
            // signing with the controlledAccountAddress should now fail 
            let tx = await controller.createAndSignTx({
                msgs: [
                    new MsgDisableSmartAccount(
                        controlledAccountAddress,
                    ),
                ],
                chainID: "test-1",
                gas: '400000',
            });
            let result = await LCD.chain1.tx.broadcastSync(tx, "test-1");
            expect(result.raw_log).toEqual("[]");
            await blockInclusion();

            // check if update authorization was successful
            try {
                await LCD.chain1.smartaccount.setting(controlledAccountAddress);
            } catch (e:any) {
                expect(e).toBeDefined();
            }
        } catch (e:any) {
            console.log(e)
            expect(e).toBeUndefined();
        }
    });

    test('Only original pk should be able to sign', async () => {
        try {
            // signing with the controller should now fail 
            let tx = await controller.createAndSignTx({
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
            expect(result.raw_log).toEqual("pubKey does not match signer address terra1wm6wwmnsdkrdugw507q4ngak589t4alq7uaqhf with signer index: 0: invalid pubkey");

            // signing with the original pk should now succeed
            tx = await wallet.createAndSignTx({
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