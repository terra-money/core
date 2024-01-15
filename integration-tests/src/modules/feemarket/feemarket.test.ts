import { Coins, Fee, MsgSubmitProposal, MsgVote } from "@terra-money/feather.js";
import { MsgParams } from "@terra-money/feather.js/dist/core/feemarket/msgs";
import { Params } from "@terra-money/feather.js/dist/core/feemarket/params";
import { VoteOption } from "@terra-money/terra.proto/cosmos/gov/v1beta1/gov";
import { blockInclusion, getLCDClient, getMnemonics, votingPeriod } from "../../helpers";

fdescribe("Feemarket Module (https://github.com/terra-money/feemarket/tree/v0.0.1-alpha.2-terra.0) ", () => {
    // Prepare environment clients, accounts and wallets
    const LCD = getLCDClient();
    const accounts = getMnemonics();
    const val2Wallet = LCD.chain1.wallet(accounts.val2);
    const val2WalletAddress = val2Wallet.key.accAddress("terra");

    test.only('Must create a new global eip1559 fees param', async () => {
        try {
            const msgProposal = new MsgSubmitProposal(
                [new MsgParams(
                    new Params(
                        '0',
                        '1',
                        '0',
                        '0',
                        '0.125',
                        '0.125',
                        '15000000',
                        '30000000',
                        '1',
                        true,
                    ),
                    'asd'
                    )],
                Coins.fromString("1000000000uluna"),
                val2WalletAddress,
                "metadata",
                "title",
                "summary"
            );
            // Create an alliance proposal sign and submit on chain-1
            let tx = await val2Wallet.createAndSignTx({
                msgs: [msgProposal],
                chainID: "test-1",
            });
            let result = await LCD.chain1.tx.broadcastSync(tx, "test-1");
            await blockInclusion();

            // Check that the proposal was created successfully
            let txResult = await LCD.chain1.tx.txInfo(result.txhash, "test-1") as any;
            expect(txResult.code).toBe(0);

            // Get the proposal id and validate exists
            let proposalId = Number(txResult.logs[0].eventsByType.submit_proposal.proposal_id[0]);
            expect(proposalId)

            // Vote for the proposal
            tx = await val2Wallet.createAndSignTx({
                msgs: [new MsgVote(
                    proposalId,
                    val2WalletAddress,
                    VoteOption.VOTE_OPTION_YES
                )],
                fee: new Fee(100_000, "0uluna"),
                chainID: "test-1",
            });
            result = await LCD.chain1.tx.broadcastSync(tx, "test-1");
            await votingPeriod();
            txResult = await LCD.chain1.tx.txInfo(result.txhash, "test-1")
            expect(txResult.code).toBe(0);
        }
        catch (e: any) {
            expect(e).toBeUndefined();
        }

        // Query the feemarket params and validate the new values
    });
});