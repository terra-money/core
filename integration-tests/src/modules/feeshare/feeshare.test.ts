import { Coins, Fee, MnemonicKey, MsgExecuteContract, MsgInstantiateContract, MsgRegisterFeeShare, MsgStoreCode } from "@terra-money/feather.js";
import fs from "fs";
import path from 'path';
import { blockInclusion, getLCDClient, getMnemonics } from "../../helpers";

describe("Feeshare Module (https://github.com/terra-money/core/tree/release/v2.6/x/feeshare) ", () => {
    // Prepare environment clients, accounts and wallets
    const LCD = getLCDClient();
    const accounts = getMnemonics();
    const wallet = LCD.chain1.wallet(accounts.feeshareMnemonic);
    const feeshareAccountAddress = accounts.feeshareMnemonic.accAddress("terra");
    const randomAccountAddress = new MnemonicKey().accAddress("terra");
    let contractAddress: string;

    // Read the reflect contract, store on chain, 
    // instantiate to be used in the following tests
    // and finally save the contract address.
    beforeAll(async () => {
        let tx = await wallet.createAndSignTx({
            msgs: [new MsgStoreCode(
                feeshareAccountAddress,
                fs.readFileSync(path.join(__dirname, "/../../contracts/reflect.wasm")).toString("base64"),
            )],
            chainID: "test-1",
        });

        let result = await LCD.chain1.tx.broadcastSync(tx, "test-1");
        await blockInclusion();
        let txResult = await LCD.chain1.tx.txInfo(result.txhash, "test-1") as any;
        let codeId = Number(txResult.logs[0].events[1].attributes[1].value);
        expect(codeId).toBeDefined();

        const msgInstantiateContract = new MsgInstantiateContract(
            feeshareAccountAddress,
            feeshareAccountAddress,
            codeId,
            {},
            Coins.fromString("1uluna"),
            "Reflect contract " + Math.random(),
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

    test('Must contain the expected module params', async () => {
        // Query feeshare module params
        const moduleParams = await LCD.chain1.feeshare.params("test-1");

        expect(moduleParams)
            .toMatchObject({
                "params": {
                    "enable_fee_share": true,
                    "developer_shares": "0.500000000000000000",
                    "allowed_denoms": []
                }
            });
    });

    test('Must register fee share', async () => {
        // Register feeshare
        let tx = await wallet.createAndSignTx({
            msgs: [new MsgRegisterFeeShare(
                contractAddress,
                feeshareAccountAddress,
                randomAccountAddress,
            )],
            chainID: "test-1",
        });

        let result = await LCD.chain1.tx.broadcastSync(tx, "test-1");
        await blockInclusion();

        // Check the tx logs
        let txResult = await LCD.chain1.tx.txInfo(result.txhash, "test-1") as any;
        expect(txResult.logs[0].events)
            .toMatchObject([{
                "type": "message",
                "attributes": [{
                    "key": "action",
                    "value": "/juno.feeshare.v1.MsgRegisterFeeShare"
                }, {
                    "key": "sender",
                    "value": feeshareAccountAddress,
                }, {
                    "key": "module",
                    "value": "feeshare"
                }]
            },
            {
                "type": "register_feeshare",
                "attributes": [{
                    "key": "contract",
                    "value": contractAddress
                }, {
                    "key": "withdrawer_address",
                    "value": randomAccountAddress,
                }]
            }])

        // Check the registered feeshares by contractAddress
        let feesharesBy = await LCD.chain1.feeshare.feeshares("test-1", contractAddress);
        expect(feesharesBy)
            .toMatchObject({
                "feeshare": {
                    "contract_address": contractAddress,
                    "deployer_address": feeshareAccountAddress,
                    "withdrawer_address": randomAccountAddress,
                }
            })
        // Check that querying all feeshares returns at least one feeshares
        let feesharesByWallet = await LCD.chain1.feeshare.feeshares("test-1");
        expect(feesharesByWallet.feeshare.length).toBeGreaterThan(0);
        await blockInclusion();

        // Send an execute message to the reflect contract
        let msgExecute = new MsgExecuteContract(
            feeshareAccountAddress,
            contractAddress,
            {
                change_owner: {
                    owner: randomAccountAddress,
                }
            },
        );
        tx = await wallet.createAndSignTx({
            msgs: [msgExecute],
            chainID: "test-1",
            fee: new Fee(200_000, "400000uluna"),
        });
        result = await LCD.chain1.tx.broadcastSync(tx, "test-1");
        await blockInclusion();

        // Check the tx logs have the expected events
        txResult = await LCD.chain1.tx.txInfo(result.txhash, "test-1") as any;
        expect(txResult.logs[0].events)
            .toMatchObject([{
                "type": "message",
                "attributes": [{
                    "key": "action",
                    "value": "/cosmwasm.wasm.v1.MsgExecuteContract"
                }, {
                    "key": "sender",
                    "value": feeshareAccountAddress
                }, {
                    "key": "module",
                    "value": "wasm"
                }]
            },
            {
                "type": "execute",
                "attributes": [{
                    "key": "_contract_address",
                    "value": contractAddress
                }]
            },
            {
                "type": "wasm",
                "attributes": [{
                    "key": "_contract_address",
                    "value": contractAddress
                }, {
                    "key": "action",
                    "value": "change_owner"
                }, {
                    "key": "owner",
                    "value": randomAccountAddress
                }]
            }
            ])
        await blockInclusion()

        // Query the random account (new owner of the contract)
        // and validate that the account has received 50% of the fees
        const bankAmount = await LCD.chain1.bank.balance(randomAccountAddress);
        expect(bankAmount[0])
            .toMatchObject(Coins.fromString("200000uluna"))
    });
});