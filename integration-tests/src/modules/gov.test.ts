import { getMnemonics } from "../helpers/mnemonics";
import { getLCDClient } from "../helpers/lcd.connection";
import { Coins, MsgVote, Fee, MsgSubmitProposal } from "@terra-money/feather.js";
import { blockInclusion, votingPeriod } from "../helpers/const";
import { VoteOption } from "@terra-money/terra.proto/cosmos/gov/v1beta1/gov";

describe("Alliance Module (https://github.com/terra-money/alliance/tree/release/v0.3.x) ", () => {
    // Prepare environment clients, accounts and wallets
    const LCD = getLCDClient();
    const accounts = getMnemonics();
    const val2Wallet = LCD.chain2.wallet(accounts.val2);
    const val2WalletAddress = val2Wallet.key.accAddress("terra");

    test('Must contain the expected module params', async () => {
        try {
            // Query Alliance module params
            const moduleParams = await LCD.chain2.gov.params("test-2");

            // Validate that the params were set correctly on genesis
            expect(moduleParams)
                .toStrictEqual({
                    "deposit_params": {
                        "max_deposit_period": "172800s",
                        "min_deposit": [
                            {
                                "amount": "10000000",
                                "denom": "uluna",
                            },
                        ],
                    },
                    "tally_params": {
                        "quorum": "0.334000000000000000",
                        "threshold": "0.500000000000000000",
                        "veto_threshold": "0.334000000000000000",
                    },
                    "voting_params": {
                        "voting_period": "4s",
                    },
                });
        }
        catch (e) {
            console.log(e)
            expect(e).toBeUndefined();
        }
    });

    test('Must submit a proposal on chain', async () => {
        try {
            const msgProposal = new MsgSubmitProposal(
                [],
                Coins.fromString("1000000000uluna"),
                val2WalletAddress,
                "metadata",
                "title",
                "summary"
            );
            // Create an alliance proposal sign and submit on chain-2
            let tx = await val2Wallet.createAndSignTx({
                msgs: [msgProposal],
                chainID: "test-2",
            });
            let result = await LCD.chain2.tx.broadcastSync(tx, "test-2");
            await blockInclusion();

            // Check that the proposal was created successfully
            let txResult = await LCD.chain2.tx.txInfo(result.txhash, "test-2") as any;
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
                chainID: "test-2",
            });
            result = await LCD.chain2.tx.broadcastSync(tx, "test-2");
            await votingPeriod();
            txResult = await LCD.chain2.tx.txInfo(result.txhash, "test-2")
            expect(txResult.code).toBe(0);
        }
        catch (e: any) {
            expect(e.response.data.message).toContain("alliance asset already exists");
        }

        // Query the alliance and check if it exists
        const res = await LCD.chain2.gov.proposals("test-2");
        expect(res).toBeDefined();
    })
});