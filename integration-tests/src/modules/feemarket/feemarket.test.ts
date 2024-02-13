import { Coins, Fee, MnemonicKey, MsgSend, MsgSubmitProposal, MsgVote } from "@terra-money/feather.js";
import { VoteOption } from "@terra-money/terra.proto/cosmos/gov/v1beta1/gov";
import { blockInclusion, getLCDClient, getMnemonics, votingPeriod } from "../../helpers";
import { FeemarketParams, MsgFeeDenomParam, MsgParams } from "@terra-money/feather.js/dist/core/feemarket";
import { execSync, exec } from 'child_process';
import path from 'path';


describe("Feemarket Module (https://github.com/terra-money/feemarket/tree/v0.0.1-alpha.2-terra.0) ", () => {
    // Prepare environment clients, accounts and wallets
    const LCD = getLCDClient();
    const accounts = getMnemonics();
    const val1Wallet = LCD.chain1.wallet(accounts.val1);
    const val1WalletAddress = val1Wallet.key.accAddress("terra");

    // Stop the relayer that way it does not affect the 
    // feemarket due the high amount of transactions
    beforeAll(() => {
        try {
            execSync("pkill relayer")
        }
        catch (e) {
            console.log(e)
        }
    });


    test('Must send a proposal to setup new global eip1559 fees param', async () => {
        try {
            // Create a FeemarketParams proposal change and 
            // submit it on chain 1.
            let tx = await val1Wallet.createAndSignTx({
                msgs: [new MsgSubmitProposal(
                    [new MsgParams(
                        FeemarketParams.fromData({
                            alpha: '0',
                            beta: '1000000000000000000',
                            theta: '0',
                            min_learning_rate: '135000000000000000',
                            max_learning_rate: '135000000000000000',
                            target_block_utilization: '5000',
                            max_block_utilization: '30000000',
                            window: '1',
                            enabled: true,
                            default_fee_denom: 'uluna'
                        }),
                        'terra10d07y265gmmuvt4z0w9aw880jnsr700juxf95n',
                    )],
                    Coins.fromString("1000000000uluna"),
                    val1WalletAddress,
                    "metadata",
                    "title",
                    "summary"
                )],
                chainID: "test-1",
            });
            let result = await LCD.chain1.tx.broadcastSync(tx, "test-1");
            await blockInclusion();

            // Check that the proposal was created successfully
            let txResult = await LCD.chain1.tx.txInfo(result.txhash, "test-1") as any;
            expect(txResult.code).toBe(0);

            // Get the proposal id and validate exists
            const proposalId = Number(txResult.logs[0].eventsByType.submit_proposal.proposal_id[0]);
            expect(proposalId).toBeDefined();

            // Vote for the proposal
            tx = await val1Wallet.createAndSignTx({
                msgs: [new MsgVote(
                    proposalId,
                    val1WalletAddress,
                    VoteOption.VOTE_OPTION_YES
                )],
                chainID: "test-1",
            });
            result = await LCD.chain1.tx.broadcastSync(tx, "test-1");
            await votingPeriod();
            txResult = await LCD.chain1.tx.txInfo(result.txhash, "test-1")
            expect(txResult.code).toBe(0);

            // Query the feemarket params and validate the new values
            let foundParams = await LCD.chain1.feemarket.params("test-1");
            expect(foundParams.toData())
                .toStrictEqual({
                    "alpha": "0",
                    "beta": "1",
                    "theta": "0",
                    "max_learning_rate": "0.135",
                    "min_learning_rate": "0.135",
                    "target_block_utilization": "5000",
                    "max_block_utilization": "30000000",
                    "window": "1",
                    "enabled": true,
                    "default_fee_denom": "uluna",
                });
        }
        catch (e: any) {
            console.log(e)
            expect(e).toBeUndefined();
        }
    });

    test('Must send a transaction and validate the fee', async () => {
        try {
            // Send a transaction to test the dynamic fees
            let sendTx = await val1Wallet.createAndSignTx({
                msgs: [
                    new MsgSend(
                        val1WalletAddress,
                        new MnemonicKey().accAddress("terra"), // To a random account
                        Coins.fromString("1uluna"),
                    ),
                ],
                chainID: "test-1",
            });
            const result = await LCD.chain1.tx.broadcastSync(sendTx, "test-1");
            await blockInclusion();
            const txResult = await LCD.chain1.tx.txInfo(result.txhash, "test-1")
            expect(txResult.code).toBe(0);

            let congested = true;
            let counter = 0;
            while (true) {
                const res = (await LCD.chain1.feemarket.feeDenomParam("test-1", "uluna"))[0];
                let gasPrice = res.baseFee;
                let minBaseFee = res.minBaseFee;
                if (congested) {
                    if (gasPrice.equals(minBaseFee)) {
                        congested = false;
                    } else {
                        expect(gasPrice.greaterThan(minBaseFee)).toBe(true);
                        console.log(`congested gasPrice: ${gasPrice.toString()}`)
                    }
                } else {
                    if (counter > 3) break;
                    if (gasPrice.greaterThan(minBaseFee)) {
                        congested = true;
                        counter = 0;
                    } else {
                        expect(gasPrice.eq(minBaseFee)).toBe(true);
                        counter++;
                        console.log(`non-congested gasPrice: ${gasPrice.toString()} counter: ${counter}`)
                    }
                }
                // wait for a new block
                await blockInclusion();
            }
        }
        catch (e: any) {
            console.log(e)
            expect(e).toBeUndefined();
        }
    });

    test('Must update fee denom param for uluna', async () => {
        try {
            // Create an state update proposal sign and submit on chain-1
            let tx = await val1Wallet.createAndSignTx({
                msgs: [new MsgSubmitProposal(
                    [new MsgFeeDenomParam(
                        'uluna',
                        '1550000000000000',
                        'terra10d07y265gmmuvt4z0w9aw880jnsr700juxf95n',
                    )],
                    Coins.fromString("1000000000uluna"),
                    val1WalletAddress,
                    "metadata",
                    "title",
                    "summary"
                )],
                chainID: "test-1",
                fee: new Fee(1000000, Coins.fromString("1000000uluna"))
            });
            let result = await LCD.chain1.tx.broadcastSync(tx, "test-1");
            await blockInclusion();

            // Check that the proposal was created successfully
            let txResult = await LCD.chain1.tx.txInfo(result.txhash, "test-1") as any;
            expect(txResult.code).toBe(0);

            // Get the proposal id and validate exists
            const proposalId = Number(txResult.logs[0].eventsByType.submit_proposal.proposal_id[0]);
            expect(proposalId).toBeDefined();

            // Vote for the proposal
            tx = await val1Wallet.createAndSignTx({
                msgs: [new MsgVote(
                    proposalId,
                    val1WalletAddress,
                    VoteOption.VOTE_OPTION_YES
                )],
                chainID: "test-1",
                fee: new Fee(1000000, Coins.fromString("1000000uluna"))
            });
            result = await LCD.chain1.tx.broadcastSync(tx, "test-1");
            await votingPeriod();

            // Validate the tx vote was casted successflully
            txResult = await LCD.chain1.tx.txInfo(result.txhash, "test-1")
            expect(txResult.code).toBe(0);

            // Query the feemarket state for uluna and validate the new values
            const res = (await LCD.chain1.feemarket.feeDenomParam("test-1", "uluna"))[0];
            expect(res.feeDenom).toEqual("uluna");
            expect(res.baseFee.toNumber()).toBeGreaterThan(0.00155);
            expect(res.minBaseFee.toString()).toStrictEqual("0.00155");
        }
        catch (e: any) {
            console.log(e)
            expect(e).toBeUndefined();
        }
    });

    // After all the tests are executed, 
    // start the relayer again
    afterAll(() => {
        // Create the path
        const pathToRelayDir = path.join(__dirname, "/../../test-data/relayer");
        // Start the relayer again
        const relayerStart = exec(`relayer start "test1-test2" -p="events" -b=100 --flush-interval="1s" --time-threshold="1s" --home="${pathToRelayDir}" > ${pathToRelayDir}/relayer.log 2>&1`)
        relayerStart.unref();
    });
});