import { getMnemonics } from "../helpers/mnemonics";
import { getLCDClient } from "../helpers/lcd.connection";
import { ContinuousVestingAccount, Coins, MnemonicKey, MsgCreateVestingAccount } from "@terra-money/feather.js";
import moment from "moment";
import { blockInclusion } from "../helpers/const";

describe("Auth Module (https://github.com/terra-money/cosmos-sdk/tree/release/v0.47.x/x/auth)", () => {
    // Prepare environment clients, accounts and wallets
    const LCD = getLCDClient();
    const accounts = getMnemonics();
    const wallet = LCD.chain1.wallet(accounts.genesisVesting1);
    const vestAccAddr1 = accounts.genesisVesting1.accAddress("terra");

    test('Must contain the expected module params', async () => {
        try {
            // Query Auth module params
            const moduleParams = await LCD.chain1.auth.parameters("test-1");

            expect(moduleParams)
                .toMatchObject({
                    "max_memo_characters": 256,
                    "tx_sig_limit": 7,
                    "tx_size_cost_per_byte": 10,
                    "sig_verify_cost_ed25519": 590,
                    "sig_verify_cost_secp256k1": 1000
                });
        }
        catch (e) {
            expect(e).toBeUndefined();
        }
    });

    test('Must have vesting accounts created on genesis', async () => {
        try {
            // Query genesis vesting account info
            const vestAccAddr = accounts.genesisVesting.accAddress("terra");
            const vestAcc = (await LCD.chain1.auth.accountInfo(vestAccAddr)) as ContinuousVestingAccount;

            // Validate the instance of the object
            expect(vestAcc)
                .toBeInstanceOf(ContinuousVestingAccount);
            // Validate the vesting start has been set in the past
            expect(vestAcc.start_time)
                .toBeLessThan(moment().unix());
            // Validate the vesting end has been set in the past
            expect(vestAcc.base_vesting_account.end_time)
                .toBeGreaterThan(moment().unix());
            // Validate the original vesting and delegated vesting
            expect(vestAcc.base_vesting_account.original_vesting)
                .toStrictEqual(Coins.fromString("10000000000uluna"));
            expect(vestAcc.base_vesting_account.delegated_vesting)
                .toStrictEqual(Coins.fromString("10000000000uluna"));

            // Validate other params from base account
            expect(vestAcc.base_vesting_account.base_account.address)
                .toBe(vestAccAddr);
            expect(vestAcc.getAccountNumber())
                .toBe(3);
            expect(vestAcc.getPublicKey())
                .toBeNull();
            expect(vestAcc.getSequenceNumber())
                .toBe(0);

            // Query the non-vested account balance
            const vestAccBalance = await LCD.chain1.bank.balance(vestAccAddr);

            // Validate the unlocked balance is still available
            expect(vestAccBalance[0])
                .toStrictEqual(Coins.fromString("990000000000uluna"));
        }
        catch (e) {
            expect(e).toBeUndefined();
        }
    });

    test('Must create a random vesting account', async () => {
        try {
            const randomAccountAddress = new MnemonicKey().accAddress("terra");
            // Register a new vesting account
            let tx = await wallet.createAndSignTx({
                msgs: [new MsgCreateVestingAccount(
                    vestAccAddr1,
                    randomAccountAddress,
                    Coins.fromString("100uluna"),
                    moment().add(1, "minute").unix(),
                    false,
                )],
                chainID: "test-1",
            });

            let result = await LCD.chain1.tx.broadcastSync(tx, "test-1");
            await blockInclusion();
            let txResult = await LCD.chain1.tx.txInfo(result.txhash, "test-1");
            expect(JSON.parse(txResult.raw_log)[0].events)
                .toEqual([{
                    "type": "message",
                    "attributes": [{
                        "key": "action",
                        "value": "/cosmos.vesting.v1beta1.MsgCreateVestingAccount"
                    }, {
                        "key": "sender",
                        "value": vestAccAddr1
                    }, {
                        "key": "module",
                        "value": "vesting"
                    }]
                },
                {
                    "type": "coin_spent",
                    "attributes": [{
                        "key": "spender",
                        "value": vestAccAddr1
                    }, {
                        "key": "amount",
                        "value": "100uluna"
                    }]
                },
                {
                    "type": "coin_received",
                    "attributes": [{
                        "key": "receiver",
                        "value": randomAccountAddress
                    }, {
                        "key": "amount",
                        "value": "100uluna"
                    }]
                },
                {
                    "type": "transfer",
                    "attributes": [{
                        "key": "recipient",
                        "value": randomAccountAddress
                    }, {
                        "key": "sender",
                        "value": vestAccAddr1
                    }, {
                        "key": "amount",
                        "value": "100uluna"
                    }]
                },
                {
                    "type": "message",
                    "attributes": [{
                        "key": "sender",
                        "value": vestAccAddr1
                    }]
                }
                ])
        }
        catch (e) {
            expect(e).toBeUndefined();
        }
    });
});