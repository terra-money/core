import { Coins, Fee, MnemonicKey, MsgSend } from "@terra-money/feather.js";
import { getMnemonics, getLCDClient, blockInclusion } from "../../helpers";

describe("FeeBurn Module (https://github.com/terra-money/core/tree/release/v2.9/x/feeburn) ", () => {
    // Prepare environment clients, accounts and wallets
    const LCD = getLCDClient();
    const accounts = getMnemonics();
    const wallet = LCD.chain1.wallet(accounts.mnemonic2);
    const randomAddress = new MnemonicKey().accAddress("terra");

    test('Must burn unused TX Fees', async () => {
        const sendTx = await wallet.createAndSignTx({
            msgs: [new MsgSend(
                wallet.key.accAddress("terra"),
                randomAddress,
                new Coins({ uluna: 1 }),
            )],
            chainID: "test-1",
            fee: new Fee(200_000, new Coins({ uluna: 3_000 })),
        });

        const result = await LCD.chain1.tx.broadcastSync(sendTx, "test-1");
        await blockInclusion();
        const txResult = await LCD.chain1.tx.txInfo(result.txhash, "test-1");
        const eventsLength = txResult.events.length;
        expect([txResult.events[eventsLength - 2], txResult.events[eventsLength - 1]])
            .toStrictEqual([{
                "type": "burn",
                "attributes": [{
                    "key": "burner",
                    "value": "terra17xpfvakm2amg962yls6f84z3kell8c5lkaeqfa",
                    "index": true
                }, {
                    "key": "amount",
                    "value": "1768uluna",
                    "index": true
                }]
            }, {
                "type": "terra.feeburn.v1.FeeBurnEvent",
                "attributes": [{
                    "key": "burn_rate",
                    "value": "\"0.589615000000000000\"",
                    "index": true
                }, {
                    "key": "fees_burn",
                    "value": "[{\"denom\":\"uluna\",\"amount\":\"1768\"}]",
                    "index": true
                }]
            }]);
    });
});