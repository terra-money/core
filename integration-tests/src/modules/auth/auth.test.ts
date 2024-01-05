import { getMnemonics, getLCDClient, blockInclusion } from "../../helpers";
import { ContinuousVestingAccount, Coins, MnemonicKey, MsgCreateVestingAccount, Coin } from "@terra-money/feather.js";
import moment from "moment";

describe("Auth Module (https://github.com/terra-money/cosmos-sdk/tree/release/v0.47.x/x/auth)", () => {
    // Prepare environment clients, accounts and wallets
    const LCD = getLCDClient();
    const accounts = getMnemonics();
    const wallet = LCD.chain1.wallet(accounts.genesisVesting1);
    const vestAccAddr1 = accounts.genesisVesting1.accAddress("terra");

    test('Must contain the expected module params', async () => {
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
    });

    test('Must have vesting accounts created on genesis', async () => {
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
        // Validate the original vesting
        expect(vestAcc.base_vesting_account.original_vesting)
            .toStrictEqual(Coins.fromString("10000000000uluna"));

        // Validate other params from base account
        expect(vestAcc.base_vesting_account.base_account.address)
            .toBe(vestAccAddr);
        expect(vestAcc.getAccountNumber())
            .toBe(4);
        expect(vestAcc.getPublicKey())
            .toBeNull();
        expect(vestAcc.getSequenceNumber())
            .toBe(0);

        // Query the non-vested account balance
        const vestAccBalance = await LCD.chain1.bank.balance(vestAccAddr);

        // Validate the unlocked balance is still available
        expect(vestAccBalance[0].get("uluna"))
            .toStrictEqual(Coin.fromString("1000000000000uluna"));
    });

    test('Must create a random vesting account', async () => {
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
        let txResult = await LCD.chain1.tx.txInfo(result.txhash, "test-1") as any;
        expect(txResult.logs[0].events)
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
            }])
    });
});