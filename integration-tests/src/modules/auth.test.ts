import { getMnemonics } from "../helpers/mnemonics";
import { getLCDClient } from "../helpers/lcd.connection";
import { Coins, MnemonicKey, MsgCreateVestingAccount} from "@terra-money/feather.js";
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
            const accountInfo = await LCD.chain1.auth.accountInfo(vestAccAddr);

            expect(accountInfo.toData())
                .toMatchObject({
                    "@type": "/cosmos.vesting.v1beta1.ContinuousVestingAccount",
                    "base_vesting_account": {
                        "@type": "/cosmos.vesting.v1beta1.BaseVestingAccount",
                        "base_account": {
                            "@type": "/cosmos.auth.v1beta1.BaseAccount",
                            "account_number": "3",
                            "address": "terra1gyf58rxglrzp343d4wkw7vzlcw6d8knp2qmg0t",
                            "pub_key": {
                                "@type": "/cosmos.crypto.secp256k1.PubKey",
                                "key": "A4VkfYoPDY1Ku7PxPU5LZDYdQE3OS/liDNmCJxsVvQxW",
                            },
                            "sequence": "1",
                        },
                        "delegated_free": [],
                        "delegated_vesting": [{
                            "amount": "10000000000",
                            "denom": "uluna",
                        },
                        ],
                        "end_time": "1797557044",
                        "original_vesting": [{
                            "amount": "10000000000",
                            "denom": "uluna",
                        },
                        ],
                    },
                    "start_time": "1697557021",
                });
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
            console.log(JSON.stringify(txResult))
        }
        catch (e) {
            console.log(JSON.stringify(e))
            expect(e).toBeUndefined();
        }
    });
});