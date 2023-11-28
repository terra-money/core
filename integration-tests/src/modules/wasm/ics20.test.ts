
import { MsgExecuteContract, MsgInstantiateContract, MsgStoreCode } from "@terra-money/feather.js";
import { getMnemonics, getLCDClient, blockInclusion } from "../../helpers";
import fs from "fs";
import path from 'path';

describe("Wasm Module (https://github.com/CosmWasm/wasmd/releases/tag/v0.45.0) ", () => {
    // Prepare environment clients, accounts and wallets
    const LCD = getLCDClient();
    const accounts = getMnemonics();
    const wallet = LCD.chain1.wallet(accounts.wasmContracts);
    const walletAddress = accounts.wasmContracts.accAddress("terra");
    let cw20BaseCodeId: number;
    let ics20CodeId: number;
    let cw20ContractAddr: string;
    let ics20ContractAddr: string;

    // Validate that wasm module has the correct params
    test('Must have the correct module params', async () => {
        const moduleParams = await LCD.chain1.wasm.params("test-1");

        expect(moduleParams)
            .toStrictEqual({
                params: {
                    "code_upload_access": {
                        "addresses": [],
                        "permission": "Everybody"
                    },
                    "instantiate_default_permission": "Everybody",
                }
            });
    })

    // Validate that wasm module has the correct params
    test('Must deploy *cw20_base* and *cw20_ics20* contracts', async () => {
        let tx = await wallet.createAndSignTx({
            msgs: [
                new MsgStoreCode(walletAddress, fs.readFileSync(path.join(__dirname, "/../../contracts/cw20_base.wasm")).toString("base64")),
                new MsgStoreCode(walletAddress, fs.readFileSync(path.join(__dirname, "/../../contracts/cw20_ics20.wasm")).toString("base64"))
            ],
            chainID: "test-1",
        });

        let result = await LCD.chain1.tx.broadcastSync(tx, "test-1");
        await blockInclusion();
        let txResult = await LCD.chain1.tx.txInfo(result.txhash, "test-1") as any;
        cw20BaseCodeId = Number(txResult.logs[0].events[1].attributes[1].value);
        expect(cw20BaseCodeId).toBeDefined();
        ics20CodeId = Number(txResult.logs[1].events[1].attributes[1].value);
        expect(ics20CodeId).toBeDefined();
    })

    describe("after contracts have been deployed", () => {
        test("Must instantiate *cw20_base* and *cw20_ics20* contract", async () => {
            let tx = await wallet.createAndSignTx({
                msgs: [new MsgInstantiateContract(
                    walletAddress,
                    walletAddress,
                    cw20BaseCodeId,
                    {
                        name: "Bitcoin",
                        symbol: "BTC",
                        decimals: 8,
                        initial_balances: [{
                            address: walletAddress,
                            amount: "100000000",
                        }]
                    },
                    undefined,
                    "A cw20 contract" + Math.random(),
                )],
                chainID: "test-1",
            });
            let result = await LCD.chain1.tx.broadcastSync(tx, "test-1");
            await blockInclusion();
            let txResult = await LCD.chain1.tx.txInfo(result.txhash, "test-1") as any;
            cw20ContractAddr = txResult.logs[0].events[1].attributes[0].value;
            expect(cw20ContractAddr).toBeDefined();

            tx = await wallet.createAndSignTx({
                msgs: [new MsgInstantiateContract(
                    walletAddress,
                    walletAddress,
                    ics20CodeId,
                    {
                        default_timeout: 60,
                        gov_contract: cw20ContractAddr,
                        allowlist: [{
                            contract: cw20ContractAddr,
                            gas_limit: 1000000,
                        }],
                        default_gas_limit: 1000000
                    },
                    undefined,
                    "A cw20 contract" + Math.random(),
                )],
                chainID: "test-1",
            });
            result = await LCD.chain1.tx.broadcastSync(tx, "test-1");
            await blockInclusion();
            txResult = await LCD.chain1.tx.txInfo(result.txhash, "test-1") as any;
            ics20ContractAddr = txResult.logs[0].events[1].attributes[0].value;
            expect(ics20ContractAddr).toBeDefined();
        })
    })

    describe("after contracts have been deployed", () => {
        test("Must instantiate *cw20_base* and *cw20_ics20* contract", async () => {
            try {
                // SubMessage to Transfer the funds thoguht the IBC channel
                // which must be parsed to base64 and embeded into the "send"
                // message. (we're not using JSON.stringify(object) because it causes and error)
                let subMsg = Buffer.from(`{"channel":"channel-0","remote_address":"${walletAddress}"}`).toString("base64");

                let tx = await wallet.createAndSignTx({
                    msgs: [new MsgExecuteContract(
                        walletAddress,
                        cw20ContractAddr,
                        {
                            send: {
                                contract: ics20ContractAddr,
                                amount: "100000",
                                msg: subMsg
                            }
                        },
                    )],
                    chainID: "test-1",
                });
                let result = await LCD.chain1.tx.broadcastSync(tx, "test-1");
                await blockInclusion();
                let txResult = await LCD.chain1.tx.txInfo(result.txhash, "test-1") as any;
                console.log("txResult", txResult)
            }
            catch (e) {
                console.log(e)
            }
        })
    })
});