import { getLCDClient, getMnemonics, blockInclusion, votingPeriod } from "../../helpers";
import { Coin, MsgTransfer, MsgCreateAlliance, Coins, MsgVote, Fee, MsgAllianceDelegate, MsgClaimDelegationRewards, MsgAllianceUndelegate, MsgDeleteAlliance, MsgSubmitProposal } from "@terra-money/feather.js";
import { VoteOption } from "@terra-money/terra.proto/cosmos/gov/v1beta1/gov";
import { Height } from "@terra-money/feather.js/dist/core/ibc/core/client/Height";

describe("Alliance Module (https://github.com/terra-money/alliance/tree/release/v0.3.x) ", () => {
    // Prepare environment clients, accounts and wallets
    const LCD = getLCDClient();
    const accounts = getMnemonics();
    const chain1Wallet = LCD.chain1.wallet(accounts.allianceMnemonic);
    const val2Wallet = LCD.chain2.wallet(accounts.val2);
    const val2WalletAddress = val2Wallet.key.accAddress("terra");
    const val2Address = val2Wallet.key.valAddress("terra");
    const allianceAccountAddress = accounts.allianceMnemonic.accAddress("terra");
    // This will be populated in the "Must create an alliance"
    let ibcCoin = Coin.fromString("1uluna");

    // Send uluna from chain-1 to chain-2 using 
    // the same wallet on both chains and start
    // an Alliance creation process
    beforeAll(async () => {
        let blockHeight = (await LCD.chain1.tendermint.blockInfo("test-1")).block.header.height;
        let tx = await chain1Wallet.createAndSignTx({
            msgs: [new MsgTransfer(
                "transfer",
                "channel-0",
                Coin.fromString("100000000uluna"),
                allianceAccountAddress,
                allianceAccountAddress,
                new Height(2, parseInt(blockHeight) + 100),
                undefined,
                ""
            )],
            chainID: "test-1",
        });

        let result = await LCD.chain1.tx.broadcastSync(tx, "test-1");
        await blockInclusion();
        let txResult = await LCD.chain1.tx.txInfo(result.txhash, "test-1") as any;
        expect(txResult).toBeDefined();

        // Check during 5 blocks for the receival 
        // of the IBC coin on chain-2
        for (let i = 0; i <= 5; i++) {
            await blockInclusion();
            let _ibcCoin = (await LCD.chain2.bank.balance(allianceAccountAddress))[0].find(c => c.denom.startsWith("ibc/"));
            if (_ibcCoin) {
                expect(_ibcCoin.denom.startsWith("ibc/")).toBeTruthy();
                break;
            }
        }
    });

    test('Must contain the expected module params', async () => {
        // Query Alliance module params
        const moduleParams = await LCD.chain2.alliance.params("test-2");

        // Validate that the params were set correctly on genesis
        expect(moduleParams.params.take_rate_claim_interval)
            .toBe("300s");
        expect(moduleParams.params.reward_delay_time)
            .toBe("0s");
    });

    test('Must create an alliance', async () => {
        // Query the IBC coin and check if there is any
        // which menas that the IBC transfer was successful
        for (let i = 0; i <= 5; i++) {
            await blockInclusion();
            let _ibcCoin = (await LCD.chain2.bank.balance(allianceAccountAddress))[0].find(c => c.denom.startsWith("ibc/"));
            if (_ibcCoin) {
                ibcCoin = _ibcCoin;
                break;
            }
        }
        expect(ibcCoin.denom.startsWith("ibc/")).toBeTruthy();

        try {
            const msgProposal = new MsgSubmitProposal(
                [new MsgCreateAlliance(
                    "terra10d07y265gmmuvt4z0w9aw880jnsr700juxf95n",
                    ibcCoin.denom,
                    "100000000000000000",
                    "0",
                    "1000000000000000000",
                    undefined,
                    {
                        "min": "100000000000000000",
                        "max": "100000000000000000"
                    })],
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

        const res = await LCD.chain2.alliance.alliance("test-2", encodeURIComponent(encodeURIComponent(ibcCoin.denom)));
        expect(res).toBeDefined();
        expect(res.alliance.denom).toBe(ibcCoin.denom);
        expect(res.alliance.reward_weight).toBe("0.100000000000000000");
        expect(res.alliance.take_rate).toBe("0.000000000000000000");
        expect(res.alliance.reward_weight_range?.min).toBe("0.100000000000000000")
        expect(res.alliance.reward_weight_range?.max).toBe("0.100000000000000000")
        expect(res.alliance.is_initialized).toBeTruthy();
    });

    describe("After Alliance has been created", () => {
        test('Must delegate to the alliance', async () => {
            const allianceWallet2 = LCD.chain2.wallet(accounts.allianceMnemonic);
            let ibcCoin = (await LCD.chain2.bank.balance(allianceAccountAddress))[0].find(c => c.denom.startsWith("ibc/")) as Coin;
            let tx = await allianceWallet2.createAndSignTx({
                msgs: [
                    new MsgAllianceDelegate(
                        allianceAccountAddress,
                        val2Address,
                        new Coin(ibcCoin.denom, 1000),
                    )
                ],
                chainID: "test-2",
            });
            let result = await LCD.chain2.tx.broadcastSync(tx, "test-2");
            await blockInclusion();

            // Check that the proposal was created successfully
            let txResult = await LCD.chain2.tx.txInfo(result.txhash, "test-2") as any;
            expect(txResult.code).toBe(0);

            // Validate the delegation event
            let delegationEvents = txResult.logs[0].eventsByType["alliance.alliance.DelegateAllianceEvent"];
            expect(delegationEvents.allianceSender)
                .toStrictEqual([`"${allianceAccountAddress}"`])
            expect(delegationEvents.coin)
                .toStrictEqual([`{\"denom\":\"${ibcCoin.denom}\",\"amount\":\"1000\"}`])
            expect(delegationEvents.newShares)
                .toStrictEqual([`"1000.000000000000000000"`])
            expect(delegationEvents.validator)
                .toStrictEqual([`"${val2Address}"`])
        });

        test('Must query one alliance validators', async () => {
            const res = await LCD.chain2.alliance.alliancesByValidator(val2Address);
            expect(res)
                .toStrictEqual({
                    "validator_addr": val2Address,
                    "total_delegation_shares": [{
                        "denom": ibcCoin.denom,
                        "amount": "1000.000000000000000000"
                    }],
                    "validator_shares": [{
                        "denom": ibcCoin.denom,
                        "amount": "1000.000000000000000000"
                    }],
                    "total_staked": [{
                        "denom": ibcCoin.denom,
                        "amount": "1000.000000000000000000"
                    }]
                })
        });

        test('Must query all alliance validators', async () => {
            const res = await LCD.chain2.alliance.alliancesByValidators("test-2");
            expect(res)
                .toStrictEqual({
                    "validators": [{
                        "validator_addr": val2Address,
                        "total_delegation_shares": [{
                            "denom": ibcCoin.denom,
                            "amount": "1000.000000000000000000"
                        }],
                        "validator_shares": [{
                            "denom": ibcCoin.denom,
                            "amount": "1000.000000000000000000"
                        }],
                        "total_staked": [{
                            "denom": ibcCoin.denom,
                            "amount": "1000.000000000000000000"
                        }]
                    }],
                    "pagination": {
                        "next_key": null,
                        "total": "1"
                    }
                })
        });

        describe("After delegation", () => {
            test("Must claim rewards from the alliance", async () => {
                const allianceWallet2 = LCD.chain2.wallet(accounts.allianceMnemonic);
                let ibcCoin = (await LCD.chain2.bank.balance(allianceAccountAddress))[0].find(c => c.denom.startsWith("ibc/")) as Coin;
                let tx = await allianceWallet2.createAndSignTx({
                    msgs: [
                        new MsgClaimDelegationRewards(
                            allianceAccountAddress,
                            val2Address,
                            ibcCoin.denom,
                        ),
                    ],
                    fee: new Fee(300_000, "100000uluna"),
                    chainID: "test-2",
                });
                let result = await LCD.chain2.tx.broadcastSync(tx, "test-2");
                await blockInclusion();

                // Check that the proposal was created successfully
                let txResult = await LCD.chain2.tx.txInfo(result.txhash, "test-2") as any;
                expect(txResult.code).toBe(0);

                // Validate the delegation event
                let claimRewardsEvent = txResult.logs[0].eventsByType["alliance.alliance.ClaimAllianceRewardsEvent"];
                expect(claimRewardsEvent.allianceSender)
                    .toStrictEqual([`"${allianceAccountAddress}"`])
                expect(claimRewardsEvent.validator)
                    .toStrictEqual([`"${val2Address}"`])
            })

            test("Must undelegate from the alliance", async () => {
                await blockInclusion();
                const allianceWallet2 = LCD.chain2.wallet(accounts.allianceMnemonic);
                let ibcCoin = (await LCD.chain2.bank.balance(allianceAccountAddress))[0].find(c => c.denom.startsWith("ibc/")) as Coin;
                let tx = await allianceWallet2.createAndSignTx({
                    msgs: [
                        new MsgAllianceUndelegate(
                            allianceAccountAddress,
                            val2Address,
                            new Coin(ibcCoin.denom, 1000),
                        ),
                    ],
                    fee: new Fee(300_000, "100000uluna"),
                    chainID: "test-2",
                });
                let result = await LCD.chain2.tx.broadcastSync(tx, "test-2");
                await blockInclusion();

                // Check that the proposal was created successfully
                let txResult = await LCD.chain2.tx.txInfo(result.txhash, "test-2") as any;
                expect(txResult.code).toBe(0);

                // Validate the delegation event
                let undelegateEvent = txResult.logs[0].eventsByType["alliance.alliance.UndelegateAllianceEvent"];
                expect(undelegateEvent.allianceSender)
                    .toStrictEqual([`"${allianceAccountAddress}"`])
                expect(undelegateEvent.coin)
                    .toStrictEqual([`{\"denom\":\"${ibcCoin.denom}\",\"amount\":\"1000\"}`])
                expect(undelegateEvent.validator)
                    .toStrictEqual([`"${val2Address}"`])
            })
        })
    })

    describe("After interacting with the Alliance", () => {
        test('Must removed the alliance using gov', async () => {
            let ibcCoin = (await LCD.chain2.bank.balance(allianceAccountAddress))[0].find(c => c.denom.startsWith("ibc/")) as Coin;

            const msgProposal = new MsgSubmitProposal(
                [new MsgDeleteAlliance(
                    "terra10d07y265gmmuvt4z0w9aw880jnsr700juxf95n",
                    ibcCoin.denom,
                )],
                Coins.fromString("1000000000uluna"),
                val2WalletAddress,
                "metadata",
                "title",
                "summary"
            );
            // Create a delete alliance proposal sign and submit on chain-2
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
                fee: new Fee(100_000, "100000uluna"),
                chainID: "test-2",
            });
            result = await LCD.chain2.tx.broadcastSync(tx, "test-2");
            await votingPeriod();
            txResult = await LCD.chain2.tx.txInfo(result.txhash, "test-2")
            expect(txResult.code).toBe(0);

            // Query the alliance and check if it exists
            const res = await LCD.chain2.alliance.alliances("test-2");
            expect(res.alliances.length).toBe(0);
        });
    })
});