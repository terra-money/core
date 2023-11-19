import { Coins, Fee, MsgSend } from "@terra-money/feather.js";
import { getMnemonics } from "../helpers/mnemonics";
import { getLCDClient } from "../helpers/lcd.connection";
import { MsgAuctionBid } from "@terra-money/feather.js/dist/core/pob/MsgAuctionBid";
import { blockInclusion } from "../helpers/const";

describe("Proposer Builder Module (https://github.com/skip-mev/pob) ", () => {
    // Prepare environment clients, accounts and wallets
    const LCD = getLCDClient();
    const accounts = getMnemonics();
    const wallet = LCD.chain1.wallet(accounts.pobMnemonic);
    const wallet11 = LCD.chain1.wallet(accounts.pobMnemonic1);

    test('Must contain the correct module params', async () => {
        try {
            // Query POB module params
            const moduleParams = await LCD.chain1.pob.params("test-1");

            expect(moduleParams)
                .toMatchObject({
                    "params": {
                        "escrow_account_address": "32sHF2qbF8xMmvwle9QEcy59Cbc=",
                        "front_running_protection": true,
                        "max_bundle_size": 2,
                        "min_bid_increment": {
                            "amount": "1",
                            "denom": "uluna",
                        },
                        "proposer_fee": "0.000000000000000000",
                        "reserve_fee": {
                            "amount": "1",
                            "denom": "uluna",
                        },
                    },
                });
        }
        catch (e) {
            expect(e).toBeUndefined();
        }
    });

    test('Must create and order two transactions in block', async () => {
        try {
            // Query block height and assert that the value is greater than 1.
            // This blockHeight will be used later on timeoutHeight
            const blockHeight = (await LCD.chain1.tendermint.blockInfo("test-1")).block.header.height;
            expect(parseInt(blockHeight)).toBeGreaterThan(1);

            // Query account info to sign the transactions offline 
            // to be included in the MsgAuctionBid
            const accInfo = await LCD.chain1.auth.accountInfo(wallet.key.accAddress("terra"));

            // **First** message to be signed using **wallet**
            const firstMsg = MsgSend.fromData({
                "@type": "/cosmos.bank.v1beta1.MsgSend",
                "from_address": accounts.pobMnemonic.accAddress("terra"),
                "to_address": accounts.pobMnemonic1.accAddress("terra"),
                "amount": [{ "denom": "uluna", "amount": "1" }]
            });
            const firstSignedSendTx = await wallet.createAndSignTx({
                msgs: [firstMsg],
                memo: "First signed tx",
                chainID: "test-1",
                accountNumber: accInfo.getAccountNumber(),
                sequence: accInfo.getSequenceNumber() + 1,
                fee: new Fee(100000, new Coins({ uluna: 100000 })),
                timeoutHeight: parseInt(blockHeight) + 20,
            });

            // **Second** message to be signed using **wallet**
            const secondMsg = MsgSend.fromData({
                "@type": "/cosmos.bank.v1beta1.MsgSend",
                "from_address": accounts.pobMnemonic.accAddress("terra"),
                "to_address": accounts.pobMnemonic1.accAddress("terra"),
                "amount": [{ "denom": "uluna", "amount": "2" }]
            });
            const secondSignedSendTx = await wallet.createAndSignTx({
                msgs: [secondMsg],
                memo: "Second signed tx",
                chainID: "test-1",
                accountNumber: accInfo.getAccountNumber(),
                sequence: accInfo.getSequenceNumber(),
                fee: new Fee(100000, new Coins({ uluna: 100000 })),
                timeoutHeight: parseInt(blockHeight) + 20,
            });

            // Create the **MsgAuctionBid** with **wallet11**.
            // 
            // The two signed transactions included in MsgAuctionBid
            // ordered as the **secondSingedTransaction** in the first position
            // and the **firstSignedTransaction** in the second position
            let buildTx = await wallet11.createAndSignTx({
                msgs: [MsgAuctionBid.fromData({
                    "@type": "/pob.builder.v1.MsgAuctionBid",
                    bid: { amount: "100000", denom: "uluna" },
                    bidder: accounts.pobMnemonic1.accAddress("terra"),
                    transactions: [secondSignedSendTx.toBytes(), firstSignedSendTx.toBytes()]
                })],
                memo: "Build block",
                chainID: "test-1",
                fee: new Fee(100000, new Coins({ uluna: 100000 })),
                timeoutHeight: parseInt(blockHeight) + 20,
            });
            const result = await LCD.chain1.tx.broadcastSync(buildTx, "test-1");
            await blockInclusion();
            const txResult = await LCD.chain1.tx.txInfo(result.txhash, "test-1");
            expect(txResult.logs).toBeDefined();
            // Recover the transactions hashes from the bundled transactions
            // to query the respective transaction data and check there are two
            const txHashes = (txResult.logs as any)[0].eventsByType.auction_bid.bundled_txs[0].split(",");
            expect(txHashes.length).toBe(2);

            // Define index to check the order of the transactions
            let index = 0;
            for await (const txHash of txHashes) {
                const txResult = await LCD.chain1.tx.txInfo(txHash, "test-1");
                const dataMsg = txResult.tx.body.messages[0].toData();
                // When the index is 0 the expected message is the secondMsg
                // because the MsgAuctionBid orders the transactions that way
                const expectedMsg = index === 0 ? secondMsg : firstMsg;
                expect(dataMsg).toMatchObject(expectedMsg.toData());

                index++;
            }
        }
        catch (e) {
            expect(e).toBeUndefined();
        }
    });
});