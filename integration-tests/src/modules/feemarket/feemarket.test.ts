import { Coins, Fee, MnemonicKey, MsgSend, MsgSubmitProposal, MsgVote } from "@terra-money/feather.js";
import { VoteOption } from "@terra-money/terra.proto/cosmos/gov/v1beta1/gov";
import { blockInclusion, getLCDClient, getMnemonics, votingPeriod } from "../../helpers";
import { FeemarketParams, MsgFeeDenomParam, MsgParams } from "@terra-money/feather.js/dist/core/feemarket";


describe("Feemarket Module (https://github.com/terra-money/feemarket/tree/v0.0.1-alpha.2-terra.0) ", () => {
    // Prepare environment clients, accounts and wallets
    const LCD = getLCDClient();
    const accounts = getMnemonics();
    const val1Wallet = LCD.chain1.wallet(accounts.val1);
    const val1WalletAddress = val1Wallet.key.accAddress("terra");

    test('Must send a proposal to setup highly volatile fees', async () => {
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
                            min_learning_rate: '129000000000000000',
                            max_learning_rate: '129000000000000000',
                            target_block_utilization: '20000',
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
                    "max_learning_rate": "0.129",
                    "min_learning_rate": "0.129",
                    "target_block_utilization": "20000",
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

    test('Must send two transactions and validate fees sensitivity', async () => {
        try {
            // Send a transaction to test the dynamic fees
            let sendTx = await val1Wallet.createAndSignTx({
                msgs: [
                    new MsgSend(
                        val1WalletAddress,
                        new MnemonicKey().accAddress("terra"), // To a random account
                        Coins.fromString("1uluna"),
                    ),
                    new MsgSend(
                        val1WalletAddress,
                        new MnemonicKey().accAddress("terra"), // To a random account
                        Coins.fromString("1uluna"),
                    ),
                ],
                fee: new Fee(150_000, Coins.fromString("1000000uluna")),
                chainID: "test-1",
            });
            const result = await LCD.chain1.tx.broadcastSync(sendTx, "test-1");
            await blockInclusion();
            const txResult = await LCD.chain1.tx.txInfo(result.txhash, "test-1")
            expect(txResult.code).toBe(0);

            let congested = true;
            let counter = 0;
            // This loop will validate that the fees will return to
            // the minimum value over time and stays there since
            // no transactions are submitted.
            // :WARNING: This test can fail when the relayer submits 
            // transactions because it spikes the fees again.
            while (true) {
                // To vaoid spamming too much wait for 250ms
                await new Promise((resolve) => setTimeout(() => resolve(250), 250));
                const res = (await LCD.chain1.feemarket.feeDenomParam("test-1", "uluna"))[0];
                if (counter == 6) break;
                if (congested) {
                    if (res.baseFee.equals(res.minBaseFee)) {
                        congested = false;
                    } else {
                        expect(res.baseFee.greaterThan(res.minBaseFee)).toBe(true);
                    }
                } else {
                    if (res.baseFee.greaterThan(res.minBaseFee)) {
                        congested = true;
                        counter = 0;
                    } else {
                        expect(res.baseFee.eq(res.minBaseFee)).toBe(true);
                        counter++;
                    }
                }
            }
        }
        catch (e: any) {
            console.log(e)
            expect(e).toBeUndefined();
        }
    });

    test('Must send a proposal to update fee denom for uluna and validate it has been updated correctly', async () => {
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
                fee: new Fee(200000, Coins.fromString("1000000uluna"))
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
                fee: new Fee(200000, Coins.fromString("1000000uluna"))
            });
            result = await LCD.chain1.tx.broadcastSync(tx, "test-1");
            await votingPeriod();

            // Validate the tx vote was casted successflully
            txResult = await LCD.chain1.tx.txInfo(result.txhash, "test-1")
            expect(txResult.code).toBe(0);

            // Query the feemarket state for uluna and validate the new values
            const res = (await LCD.chain1.feemarket.feeDenomParam("test-1", "uluna"))[0];
            expect(res.feeDenom).toEqual("uluna");
            expect(res.baseFee.toNumber()).toBeGreaterThanOrEqual(0.00155);
            expect(res.minBaseFee.toString()).toStrictEqual("0.00155");
        }
        catch (e: any) {
            console.log(e)
            expect(e).toBeUndefined();
        }
    });

    test('Must modify the fees back to the default values', async () => {
        try {
            // Create a FeemarketParams proposal change and 
            // submit it on chain 1.
            let tx = await val1Wallet.createAndSignTx({
                msgs: [new MsgSubmitProposal(
                    [new MsgFeeDenomParam(
                        'uluna',
                        '1500000000000000',
                        'terra10d07y265gmmuvt4z0w9aw880jnsr700juxf95n',
                    ),new MsgParams(
                        FeemarketParams.fromData({
                            alpha: '0',
                            beta: '1000000000000000000',
                            theta: '0',
                            min_learning_rate: '125000000000000000',
                            max_learning_rate: '125000000000000000',
                            target_block_utilization: '15000000',
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
                fee: new Fee(200_000, Coins.fromString("100000uluna"))
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
                fee: new Fee(200_000, Coins.fromString("100000uluna"))
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
                    "max_learning_rate": "0.125",
                    "min_learning_rate": "0.125",
                    "target_block_utilization": "15000000",
                    "max_block_utilization": "30000000",
                    "window": "1",
                    "enabled": true,
                    "default_fee_denom": "uluna",
                });

            // Query the feemarket state for uluna and validate the new values
            const res = (await LCD.chain1.feemarket.feeDenomParam("test-1", "uluna"))[0];
            expect(res.feeDenom).toEqual("uluna");
            expect(res.baseFee.toNumber()).toBeGreaterThan(0.0015);
            expect(res.minBaseFee.toString()).toStrictEqual("0.0015");
            // This loop will validate that the fees return to
            // the minimum value and allows the execution of the
            // next test. It is done this way to avoid possible
            // failures in the next test because of the fees.
            while (true) {
                // To vaoid spamming too much wait for 250ms
                await new Promise((resolve) => setTimeout(() => resolve(250), 250));
                const res = (await LCD.chain1.feemarket.feeDenomParam("test-1", "uluna"))[0];
                if (res.baseFee.equals(res.minBaseFee)) break;
            }
        }
        catch (e: any) {
            console.log(e)
            expect(e).toBeUndefined();
        }
    });
});