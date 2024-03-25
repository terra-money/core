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

            const receiverBalanceBefore = await LCD.chain1.bank.balance(receiver);

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

            // check that MsgSend succeeds
            const receiverBalanceAfter = await LCD.chain1.bank.balance(receiver);
            const deltaBalance = receiverBalanceAfter[0].sub(receiverBalanceBefore[0]);
            expect(deltaBalance).toEqual(coinsToSend);
        } catch (e:any) {
            console.log(e)
            expect(e).toBeUndefined();
        }
    });
    
    test('Transaction should fail for sending over limit', async () => {
        try {
            let setting = await LCD.chain1.smartaccount.setting(smartaccAddress);
            expect(setting.postTransaction).toEqual([limitContractAddress]);

            // should have 940000uluna at this point 60000
            const balanceBefore = await LCD.chain1.bank.balance(smartaccAddress);
            // leave 23905uluna for fees so transaction will not fail due to insufficient funds
            // should cost 23705uluna
            const coinsToSend = balanceBefore[0].sub(Coins.fromString("23905uluna"));

            const coinBefore = balanceBefore[0].find((coin: Coin) => coin.denom === "uluna");
            expect(coinBefore).toBeDefined();

            let tx = await wallet.createAndSignTx({
                msgs: [
                    new MsgSend(
                        smartaccAddress,
                        receiver,
                        coinsToSend,
                    ),
                ],
                chainID: "test-1",
            });
            const fee_coins = tx.toData().auth_info.fee.amount;

            // fee_coins[0].amount string to number
            const fee_amount = parseInt(fee_coins[0].amount);

            await LCD.chain1.tx.broadcastSync(tx, "test-1");
            await blockInclusion();

            // check that MsgSend failed
            const balanceAfter = await LCD.chain1.bank.balance(smartaccAddress);
            const coinAfter = balanceAfter[0].find((coin: Coin) => coin.denom === "uluna");
            expect(coinAfter).toBeDefined();
            
            // check that only fees were deducted
            const coinAfter_amount = parseInt(coinAfter!.amount.toString());
            const balaceBeforeMinusFee = coinBefore!.amount.toNumber() - fee_amount;
            expect(balaceBeforeMinusFee).toEqual(coinAfter_amount);
        } catch (e:any) {
            console.log(e)
            expect(e).toBeUndefined();
        }
    });
});