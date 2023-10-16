import { Coins, Fee, MsgSend } from "@terra-money/feather.js";
import { getAccounts } from "../../helpers/accounts";
import { getLCDClient } from "../../helpers/lcd.connection";
import { MsgAuctionBid } from "@terra-money/feather.js/dist/core/pob/MsgAuctionBid";

describe("POB ", () => {
    const LCD = getLCDClient();
    const accounts = getAccounts();

    test('Should create and order two transactions in block', async () => {
        try {
            const blockHeight = (await LCD.chain1.tendermint.blockInfo("test-1")).block.header.height;
            expect(parseInt(blockHeight)).toBeGreaterThan(1);
            
            const wallet = LCD.chain1.wallet(accounts.wallet1);
            const accInfo = await LCD.chain1.auth.accountInfo(wallet.key.accAddress("terra"));
            const firstSignedSendTx = await wallet.createAndSignTx({
                msgs: [MsgSend.fromData({
                    "@type": "/cosmos.bank.v1beta1.MsgSend",
                    "from_address": accounts.wallet1.accAddress("terra"),
                    "to_address": accounts.wallet11.accAddress("terra"),
                    "amount": [{ "denom": "uluna", "amount": "1" }]
                })],
                memo: "First signed tx",
                chainID: "test-1",
                accountNumber: accInfo.getAccountNumber(),
                sequence: accInfo.getSequenceNumber() + 1,
                fee: new Fee(100000, new Coins({ uluna: 100000 })),
                timeoutHeight: parseInt(blockHeight) + 20,
            });
            const secondSignedSendTx = await wallet.createAndSignTx({
                msgs: [MsgSend.fromData({
                    "@type": "/cosmos.bank.v1beta1.MsgSend",
                    "from_address": accounts.wallet1.accAddress("terra"),
                    "to_address": accounts.wallet11.accAddress("terra"),
                    "amount": [{ "denom": "uluna", "amount": "1" }]
                })],
                memo: "Second signed tx",
                chainID: "test-1",
                accountNumber: accInfo.getAccountNumber(),
                sequence: accInfo.getSequenceNumber() + 2,
                fee: new Fee(100000, new Coins({ uluna: 100000 })),
                timeoutHeight: parseInt(blockHeight) + 20,
            });
    
            const wallet11 = LCD.chain1.wallet(accounts.wallet11);
            let buildTx = await wallet11.createAndSignTx({
                msgs: [MsgAuctionBid.fromData({
                    "@type": "/pob.builder.v1.MsgAuctionBid",
                    bid: {amount: "100000",denom: "uluna"},
                    bidder: accounts.wallet1.accAddress("terra"),
                    transactions: [secondSignedSendTx.toBytes(), firstSignedSendTx.toBytes()]
                })],
                memo: "Build block",
                chainID: "test-1",
                fee: new Fee(100000, new Coins({ uluna: 100000 })),
                timeoutHeight: parseInt(blockHeight) + 20,
            });
            const result = await LCD.chain1.tx.broadcastSync(buildTx, "test-1");
            console.log(result)
        }
        catch(e){
            expect(e).toBeUndefined();
        }
    });
});