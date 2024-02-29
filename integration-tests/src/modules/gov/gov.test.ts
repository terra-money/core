import { getLCDClient, blockInclusion, votingPeriod, getMnemonics } from "../../helpers";
import { Coins, MsgVote, Fee, MsgSubmitProposal, Proposal, Int } from "@terra-money/feather.js";
import { ProposalStatus, VoteOption } from "@terra-money/terra.proto/cosmos/gov/v1beta1/gov";

describe("Governance Module (https://github.com/terra-money/cosmos-sdk/tree/release/v0.47.x/x/gov) ", () => {
    // Prepare environment clients, accounts and wallets
    const LCD = getLCDClient();
    const accounts = getMnemonics();
    const val2Wallet = LCD.chain2.wallet(accounts.val2);
    const val2WalletAddress = val2Wallet.key.accAddress("terra");
    let proposalId = 0; // Will be populated on "Must submit a proposal on chain"

    test('Must contain the expected module params', async () => {
        // Query All gov module params
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
                    "voting_period": "2s",
                },
                "params": {
                    "burn_proposal_deposit_prevote": false,
                    "burn_vote_quorum": false,
                    "burn_vote_veto": true,
                    "max_deposit_period": "172800s",
                    "min_deposit": [{
                        "amount": "10000000",
                        "denom": "uluna",
                    }],
                    "min_initial_deposit_ratio": "0.000000000000000000",
                    "quorum": "0.334000000000000000",
                    "threshold": "0.500000000000000000",
                    "veto_threshold": "0.334000000000000000",
                    "voting_period": "2s",
                },
            });

        // Query tally module params
        const tallyParams = await LCD.chain2.gov.tallyParams("test-2");
        // Validate that the params were set correctly on genesis
        expect(tallyParams)
            .toStrictEqual({
                "deposit_params": null,
                "voting_params": null,
                "tally_params": {
                    "quorum": "0.334000000000000000",
                    "threshold": "0.500000000000000000",
                    "veto_threshold": "0.334000000000000000",
                },
                "params": {
                    "burn_proposal_deposit_prevote": false,
                    "burn_vote_quorum": false,
                    "burn_vote_veto": true,
                    "max_deposit_period": "172800s",
                    "min_deposit": [{
                        "amount": "10000000",
                        "denom": "uluna",
                    }],
                    "min_initial_deposit_ratio": "0.000000000000000000",
                    "quorum": "0.334000000000000000",
                    "threshold": "0.500000000000000000",
                    "veto_threshold": "0.334000000000000000",
                    "voting_period": "2s",
                },
            });

        // Query voting gov module params
        const votingParams = await LCD.chain2.gov.votingParams("test-2");
        // Validate that the params were set correctly on genesis
        expect(votingParams)
            .toStrictEqual({
                "deposit_params": null,
                "tally_params": null,
                "voting_params": {
                    "voting_period": "2s",
                },
                "params": {
                    "burn_proposal_deposit_prevote": false,
                    "burn_vote_quorum": false,
                    "burn_vote_veto": true,
                    "max_deposit_period": "172800s",
                    "min_deposit": [{
                        "amount": "10000000",
                        "denom": "uluna",
                    }],
                    "min_initial_deposit_ratio": "0.000000000000000000",
                    "quorum": "0.334000000000000000",
                    "threshold": "0.500000000000000000",
                    "veto_threshold": "0.334000000000000000",
                    "voting_period": "2s",
                },
            });


        // Query deposit gov module params
        const depositParams = await LCD.chain2.gov.depositParams("test-2");
        // Validate that the params were set correctly on genesis
        expect(depositParams)
            .toStrictEqual({
                "voting_params": null,
                "tally_params": null,
                "deposit_params": {
                    "max_deposit_period": "172800s",
                    "min_deposit": [
                        {
                            "amount": "10000000",
                            "denom": "uluna",
                        },
                    ],
                },
                "params": {
                    "burn_proposal_deposit_prevote": false,
                    "burn_vote_quorum": false,
                    "burn_vote_veto": true,
                    "max_deposit_period": "172800s",
                    "min_deposit": [{
                        "amount": "10000000",
                        "denom": "uluna",
                    }],
                    "min_initial_deposit_ratio": "0.000000000000000000",
                    "quorum": "0.334000000000000000",
                    "threshold": "0.500000000000000000",
                    "veto_threshold": "0.334000000000000000",
                    "voting_period": "2s",
                },
            });
    });

    test('Must submit an empty proposal on chain', async () => {
        try {
            const msgProposal = new MsgSubmitProposal(
                [],
                Coins.fromString("1000000000uluna"),
                val2WalletAddress,
                "METADATA Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec nec varius odio. Phasellus tellus felis, varius ut sapien sit amet, imperdiet vehicula metus. Nullam convallis, erat sit amet ultrices ornare, quam metus ornare elit, quis sollicitudin dolor lorem non risus. Pellentesque pretium augue Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec nec varius odio. Phasellus tellus felis, varius ut sapien sit amet, imperdiet vehicula metus. Nullam convallis, erat sit amet ultrices ornare, quam metus ornare elit, quis sollicitudin dolor lorem non risus. Pellentesque pretium augue Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec nec varius odio. Phasellus tellus felis, varius ut sapien sit amet, imperdiet vehicula metus. Nullam convallis, erat sit amet ultrices ornare, quam metus ornare elit, quis sollicitudin dolor lorem non risus. Pellentesque pretium augue Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec nec varius odio. Phasellus tellus felis, varius ut sapien sit amet, imperdiet vehicula metus. Nullam convallis, erat sit amet ultrices ornare, quam metus ornare elit, quis sollicitudin dolor lorem non risus. Pellentesque pretium augue Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec nec varius odio. Phasellus tellus felis, varius ut sapien sit amet, imperdiet vehicula metus. Nullam convallis, erat sit amet ultrices ornare, quam metus ornare elit, quis sollicitudin dolor lorem non risus. Pellentesque pretium augue Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec nec varius odio. Phasellus tellus felis, varius ut sapien sit amet, imperdiet vehicula metus. Nullam convallis, erat sit amet ultrices ornare, quam metus ornare elit, quis sollicitudin dolor lorem non risus. Pellentesque pretium augue",
                "TITLE Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec nec varius odio. Phasellus tellus felis, varius ut sapien sit amet, imperdiet vehicula metus. Nullam convallis, erat sit amet ultrices ornare, quam metus ornare elit, quis sollicitudin dolor lorem non risus. Pellentesque pretium augue Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec nec varius odio. Phasellus tellus felis, varius ut sapien sit amet, imperdiet vehicula metus. Nullam convallis, erat sit amet ultrices ornare, quam metus ornare elit, quis sollicitudin dolor lorem non risus. Pellentesque pretium augue Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec nec varius odio. Phasellus tellus felis, varius ut sapien sit amet, imperdiet vehicula metus. Nullam convallis, erat sit amet ultrices ornare, quam metus ornare elit, quis sollicitudin dolor lorem non risus. Pellentesque pretium augue Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec nec varius odio. Phasellus tellus felis, varius ut sapien sit amet, imperdiet vehicula metus. Nullam convallis, erat sit amet ultrices ornare, quam metus ornare elit, quis sollicitudin dolor lorem non risus. Pellentesque pretium augue Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec nec varius odio. Phasellus tellus felis, varius ut sapien sit amet, imperdiet vehicula metus. Nullam convallis, erat sit amet ultrices ornare, quam metus ornare elit, quis sollicitudin dolor lorem non risus. Pellentesque pretium augue Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec nec varius odio. Phasellus tellus felis, varius ut sapien sit amet, imperdiet vehicula metus. Nullam convallis, erat sit amet ultrices ornare, quam metus ornare elit, quis sollicitudin dolor lorem non risus. Pellentesque pretium augue",
                "SUMMARY Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec nec varius odio. Phasellus tellus felis, varius ut sapien sit amet, imperdiet vehicula metus. Nullam convallis, erat sit amet ultrices ornare, quam metus ornare elit, quis sollicitudin dolor lorem non risus. Pellentesque pretium augue Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec nec varius odio. Phasellus tellus felis, varius ut sapien sit amet, imperdiet vehicula metus. Nullam convallis, erat sit amet ultrices ornare, quam metus ornare elit, quis sollicitudin dolor lorem non risus. Pellentesque pretium augue Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec nec varius odio. Phasellus tellus felis, varius ut sapien sit amet, imperdiet vehicula metus. Nullam convallis, erat sit amet ultrices ornare, quam metus ornare elit, quis sollicitudin dolor lorem non risus. Pellentesque pretium augue Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec nec varius odio. Phasellus tellus felis, varius ut sapien sit amet, imperdiet vehicula metus. Nullam convallis, erat sit amet ultrices ornare, quam metus ornare elit, quis sollicitudin dolor lorem non risus. Pellentesque pretium augue Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec nec varius odio. Phasellus tellus felis, varius ut sapien sit amet, imperdiet vehicula metus. Nullam convallis, erat sit amet ultrices ornare, quam metus ornare elit, quis sollicitudin dolor lorem non risus. Pellentesque pretium augue Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec nec varius odio. Phasellus tellus felis, varius ut sapien sit amet, imperdiet vehicula metus. Nullam convallis, erat sit amet ultrices ornare, quam metus ornare elit, quis sollicitudin dolor lorem non risus. Pellentesque pretium augue"
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
            proposalId = Number(txResult.logs[0].eventsByType.submit_proposal.proposal_id[0]);
            expect(proposalId)

            // Vote for the proposal
            tx = await val2Wallet.createAndSignTx({
                msgs: [new MsgVote(
                    proposalId,
                    val2WalletAddress,
                    VoteOption.VOTE_OPTION_YES
                )],
                fee: new Fee(100_000, "100000uluna"),
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
    });

    describe("After submitting the proposal on chain", () => {
        test('Must query the proposals', async () => {
            // Query the alliance and check if it exists
            const res = await LCD.chain2.gov.proposals("test-2");
            let proposal;
            for (const prop of res.proposals) {
                if (prop.id === proposalId) {
                    proposal = prop;
                }
            }
            expect(proposal).toBeDefined();
            proposal = proposal as Proposal;
            expect(proposal.id).toBe(proposalId);
            expect(proposal.proposer).toBe(val2WalletAddress);
            expect(proposal.metadata).toBe("METADATA Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec nec varius odio. Phasellus tellus felis, varius ut sapien sit amet, imperdiet vehicula metus. Nullam convallis, erat sit amet ultrices ornare, quam metus ornare elit, quis sollicitudin dolor lorem non risus. Pellentesque pretium augue Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec nec varius odio. Phasellus tellus felis, varius ut sapien sit amet, imperdiet vehicula metus. Nullam convallis, erat sit amet ultrices ornare, quam metus ornare elit, quis sollicitudin dolor lorem non risus. Pellentesque pretium augue Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec nec varius odio. Phasellus tellus felis, varius ut sapien sit amet, imperdiet vehicula metus. Nullam convallis, erat sit amet ultrices ornare, quam metus ornare elit, quis sollicitudin dolor lorem non risus. Pellentesque pretium augue Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec nec varius odio. Phasellus tellus felis, varius ut sapien sit amet, imperdiet vehicula metus. Nullam convallis, erat sit amet ultrices ornare, quam metus ornare elit, quis sollicitudin dolor lorem non risus. Pellentesque pretium augue Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec nec varius odio. Phasellus tellus felis, varius ut sapien sit amet, imperdiet vehicula metus. Nullam convallis, erat sit amet ultrices ornare, quam metus ornare elit, quis sollicitudin dolor lorem non risus. Pellentesque pretium augue Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec nec varius odio. Phasellus tellus felis, varius ut sapien sit amet, imperdiet vehicula metus. Nullam convallis, erat sit amet ultrices ornare, quam metus ornare elit, quis sollicitudin dolor lorem non risus. Pellentesque pretium augue");
            expect(proposal.title).toBe("TITLE Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec nec varius odio. Phasellus tellus felis, varius ut sapien sit amet, imperdiet vehicula metus. Nullam convallis, erat sit amet ultrices ornare, quam metus ornare elit, quis sollicitudin dolor lorem non risus. Pellentesque pretium augue Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec nec varius odio. Phasellus tellus felis, varius ut sapien sit amet, imperdiet vehicula metus. Nullam convallis, erat sit amet ultrices ornare, quam metus ornare elit, quis sollicitudin dolor lorem non risus. Pellentesque pretium augue Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec nec varius odio. Phasellus tellus felis, varius ut sapien sit amet, imperdiet vehicula metus. Nullam convallis, erat sit amet ultrices ornare, quam metus ornare elit, quis sollicitudin dolor lorem non risus. Pellentesque pretium augue Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec nec varius odio. Phasellus tellus felis, varius ut sapien sit amet, imperdiet vehicula metus. Nullam convallis, erat sit amet ultrices ornare, quam metus ornare elit, quis sollicitudin dolor lorem non risus. Pellentesque pretium augue Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec nec varius odio. Phasellus tellus felis, varius ut sapien sit amet, imperdiet vehicula metus. Nullam convallis, erat sit amet ultrices ornare, quam metus ornare elit, quis sollicitudin dolor lorem non risus. Pellentesque pretium augue Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec nec varius odio. Phasellus tellus felis, varius ut sapien sit amet, imperdiet vehicula metus. Nullam convallis, erat sit amet ultrices ornare, quam metus ornare elit, quis sollicitudin dolor lorem non risus. Pellentesque pretium augue");
            expect(proposal.summary).toBe("SUMMARY Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec nec varius odio. Phasellus tellus felis, varius ut sapien sit amet, imperdiet vehicula metus. Nullam convallis, erat sit amet ultrices ornare, quam metus ornare elit, quis sollicitudin dolor lorem non risus. Pellentesque pretium augue Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec nec varius odio. Phasellus tellus felis, varius ut sapien sit amet, imperdiet vehicula metus. Nullam convallis, erat sit amet ultrices ornare, quam metus ornare elit, quis sollicitudin dolor lorem non risus. Pellentesque pretium augue Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec nec varius odio. Phasellus tellus felis, varius ut sapien sit amet, imperdiet vehicula metus. Nullam convallis, erat sit amet ultrices ornare, quam metus ornare elit, quis sollicitudin dolor lorem non risus. Pellentesque pretium augue Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec nec varius odio. Phasellus tellus felis, varius ut sapien sit amet, imperdiet vehicula metus. Nullam convallis, erat sit amet ultrices ornare, quam metus ornare elit, quis sollicitudin dolor lorem non risus. Pellentesque pretium augue Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec nec varius odio. Phasellus tellus felis, varius ut sapien sit amet, imperdiet vehicula metus. Nullam convallis, erat sit amet ultrices ornare, quam metus ornare elit, quis sollicitudin dolor lorem non risus. Pellentesque pretium augue Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec nec varius odio. Phasellus tellus felis, varius ut sapien sit amet, imperdiet vehicula metus. Nullam convallis, erat sit amet ultrices ornare, quam metus ornare elit, quis sollicitudin dolor lorem non risus. Pellentesque pretium augue");
            expect(proposal.messages.length).toBe(0);
            expect(proposal.total_deposit).toStrictEqual(Coins.fromString("1000000000uluna"));
            expect(proposal.submit_time.getTime()).toBeLessThan(Date.now());
            expect(proposal.voting_start_time?.getTime()).toBeLessThan(Date.now());
            expect(proposal.deposit_end_time?.getTime()).toBeGreaterThan(Date.now());
            expect(proposal.voting_end_time?.getTime()).toBeLessThan(Date.now());
            expect(proposal.status).toBe(ProposalStatus.PROPOSAL_STATUS_PASSED)
        });

        test('Must query the tally', async () => {
            // Query the alliance and check if it exists
            const res = await LCD.chain2.gov.tally("test-2", proposalId);
            expect(res).toBeDefined();
            expect(res.yes_count.gte(7000000000)).toBeTruthy();
            expect(res.abstain_count).toStrictEqual(new Int(0));
            expect(res.no_count).toStrictEqual(new Int(0));
            expect(res.no_with_veto_count).toStrictEqual(new Int(0));
        });
    })
});