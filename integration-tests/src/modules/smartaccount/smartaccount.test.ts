import { Coins, MsgExecuteContract, MsgInstantiateContract, MsgSend, MsgStoreCode, PublicKey } from "@terra-money/feather.js";
import { AuthorizationMsg, MsgCreateSmartAccount, MsgUpdateAuthorization } from "@terra-money/feather.js/dist/core/smartaccount";
import fs from "fs";
import path from 'path';
import { blockInclusion, getLCDClient, getMnemonics } from "../../helpers";

describe("Smartaccount Module (https://github.com/terra-money/core/tree/release/v2.6/x/smartaccount) ", () => {
    // Prepare environment clients, accounts and wallets
    const LCD = getLCDClient();
    const accounts = getMnemonics();
    const wallet = LCD.chain1.wallet(accounts.mnemonic5);
    const controlledAccountAddress = accounts.mnemonic5.accAddress("terra");
    
    // const controller = accounts.mnemonic4.publicKey;
    const pubkey = accounts.mnemonic4.publicKey;
    expect(pubkey).toBeDefined();

    // TODO: convert pubkey to base64 string similar to golang pubkey.Bytes()
    const pubkeybb = pubkey as PublicKey
    const ggg = pubkeybb.toAmino()

    const key = ggg.value as string;
    console.log(key)
    const pubkeyBs = Buffer.from(key);
    const initMsg = {
        initialization: {
            sender: controlledAccountAddress,
            account: controlledAccountAddress,
            public_key: pubkeyBs,
        }
    }
    // marshal initMsg to bytes similar to json.Marshal in golang
    const initMsgBytes = Buffer.from(JSON.stringify(initMsg)).toString('base64');

    const deployer = LCD.chain1.wallet(accounts.tokenFactoryMnemonic);
    const deployerAddress = accounts.tokenFactoryMnemonic.accAddress("terra");

    let authContractAddress: string;

    test.only('Deploy smart account auth contract and initialize priv key for wallet', async () => {
        try {
            let tx = await deployer.createAndSignTx({
                msgs: [new MsgStoreCode(
                    deployerAddress,
                    fs.readFileSync(path.join(__dirname, "/../../contracts/smart_auth_contract.wasm")).toString("base64"),
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

            // add pubkey of controller for smart account
            tx = await wallet.createAndSignTx({
                msgs: [new MsgExecuteContract(
                    controlledAccountAddress,
                    authContractAddress,
                    initMsg,
                )],
                chainID: 'test-1',
                gas: '400000',
            });
            result = await LCD.chain1.tx.broadcastSync(tx, "test-1");
            await blockInclusion();
            txResult = await LCD.chain1.tx.txInfo(result.txhash, "test-1") as any;
            console.log(txResult);
            codeId = Number(txResult.logs[0].events[1].attributes[1].value);
            expect(codeId).toBeDefined();
        } catch(e: any) {
            expect(e).toBeUndefined();
        }
    });

    test('Create new smart account and give control to controller', async () => {
        try {
            // create the smartaccount
            let tx = await wallet.createAndSignTx({
                msgs: [new MsgCreateSmartAccount(
                    controlledAccountAddress
                )],
                chainID: 'test-1',
                gas: '400000',
            });
            let result = await LCD.chain1.tx.broadcastSync(tx, "test-1");
            await blockInclusion();
            let txResult = await LCD.chain1.tx.txInfo(result.txhash, "test-1") as any;
            let codeId = Number(txResult.logs[0].events[1].attributes[1].value);
            expect(codeId).toBeDefined();
            // Query smartaccount setting for the smart waccount
            const setting = await LCD.chain1.smartaccount.setting(controlledAccountAddress);

            expect(setting.toData())
                .toEqual({
                    owner: controlledAccountAddress,
                    authorization: [],
                    post_transaction: [],
                    pre_transaction: [],
                    fallback: true,
                });
            
            // give control to controller
            tx = await wallet.createAndSignTx({
                msgs: [new MsgUpdateAuthorization(
                    controlledAccountAddress,
                    false,
                    [new AuthorizationMsg(authContractAddress, initMsgBytes)],
                )],
                chainID: 'test-1',
                gas: '400000',
            });
            result = await LCD.chain1.tx.broadcastSync(tx, "test-1");
            await blockInclusion();
            txResult = await LCD.chain1.tx.txInfo(result.txhash, "test-1") as any;
            codeId = Number(txResult.logs[0].events[1].attributes[1].value);
            expect(codeId).toBeDefined();

            // signing with the controlledAccountAddress should now fail 
            tx = await wallet.createAndSignTx({
                msgs: [
                    new MsgSend(
                        controlledAccountAddress,
                        controlledAccountAddress,
                        Coins.fromString("1uluna"),
                    ),
                ],
                chainID: "test-1",
            });
            result = await LCD.chain1.tx.broadcastSync(tx, "test-1");
            await blockInclusion();
            txResult = await LCD.chain1.tx.txInfo(result.txhash, "test-1") as any;
            console.log(txResult);
            codeId = Number(txResult.logs[0].events[1].attributes[1].value);
            expect(codeId).toBeDefined();
        } catch (e:any) {
            console.log(e);
            expect(e).toBeUndefined();
        }
    });
});