import { Coins, MsgSend, MsgSubmitProposal, MsgVote } from "@terra-money/feather.js";
import { Params } from "@terra-money/feather.js/dist/core/feemarket/params";
import { MsgFeeDenomParam, MsgParams } from "@terra-money/feather.js/dist/core/feemarket/proposals";
import { VoteOption } from "@terra-money/terra.proto/cosmos/gov/v1beta1/gov";
import BigNumber from 'bignumber.js';
import { blockInclusion, getLCDClient, getMnemonics, votingPeriod } from "../../helpers";

describe("Feemarket Module (https://github.com/terra-money/feemarket/tree/v0.0.1-alpha.2-terra.0) ", () => {
    // Prepare environment clients, accounts and wallets
    const LCD = getLCDClient();
    const accounts = getMnemonics();
    const val1Wallet = LCD.chain1.wallet(accounts.val1);
    const val1WalletAddress = val1Wallet.key.accAddress("terra");
    const rly1Wallet = LCD.chain1.wallet(accounts.rly1);
    const rly1WalletAddress = rly1Wallet.key.accAddress("terra");

    test('Must create a new global eip1559 fees param + test dynamic fees', async () => {
        try {
            const params = new Params(
                '0',    
                '1000000000000000000',
                '0',
                '135000000000000000',
                '135000000000000000',
                '5000',
                '30000000',
                '1',
                true,
                'uluna',
            )
            const msgProposal = new MsgSubmitProposal(
                [new MsgParams(
                    params,
                    'terra10d07y265gmmuvt4z0w9aw880jnsr700juxf95n',
                )],
                Coins.fromString("1000000000uluna"),
                val1WalletAddress,
                "metadata",
                "title",
                "summary"
            );

            // Create an update params proposal sign and submit on chain-1
            let tx = await val1Wallet.createAndSignTx({
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
            let foundParams = await LCD.chain1.feemarket.params("test-1") as any;
            checkParams(foundParams.params, params)
            
            // Start fee tests
            for (let i=0;i<4;i++) {
                let sendTx = await val1Wallet.createAndSignTx({
                    msgs: [
                        new MsgSend(
                            val1WalletAddress,
                            rly1WalletAddress,
                            Coins.fromString("1uluna"),
                        ),
                    ],
                    chainID: "test-1",
                });
                result = await LCD.chain1.tx.broadcastSync(sendTx, "test-1");
                await blockInclusion();
            }

            const minGasPrice = BigNumber("0.0015");
            let congested = true;
            let counter = 0;
            for (let i = 0; i < 100; i++) {
                const gasPrice = await getGasPrice("test-1", "uluna")
                if (congested) {
                    if (gasPrice.isEqualTo(minGasPrice)) {
                        congested = false;
                    } else {
                        expect(gasPrice.isGreaterThan(minGasPrice)).toBe(true);
                        console.log(`congested gasPrice: ${gasPrice.toString()}`)
                    }
                } else {
                    if (counter > 5) break;
                    if (gasPrice.isGreaterThan(minGasPrice)) {
                        congested = true;
                        counter = 0;
                    } else {
                        expect(gasPrice.eq(minGasPrice)).toBe(true);
                        counter++;
                        console.log(`non-congested gasPrice: ${gasPrice.toString()} counter: ${counter}`)
                    }
                }
                // wait for 1 sec
                await new Promise(resolve => setTimeout(resolve, 1000));
            }
        }
        catch (e: any) {
            expect(e).toBeFalsy();
        }
    });

    const getGasPrice = async (chainId: string, feeDenom: string): Promise<BigNumber> => {
        const foundFdp = await LCD.chain1.feemarket.feeDenomParam(chainId, feeDenom) as  any;
        const fdp = foundFdp.fee_denom_params[0] as any;
        const gasPrice = BigNumber(fdp.base_fee)
        return gasPrice
    }

    test('Must update feedenomparam for uluna', async () => {
        try {
            const msgProposal = new MsgSubmitProposal(
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
            );
            // Create an state update proposal sign and submit on chain-1
            let tx = await val1Wallet.createAndSignTx({
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

            // Query the feemarket state for uluna and validate the new values
            let foundFdp = await LCD.chain1.feemarket.feeDenomParam("test-1", "uluna") as any;
            expect(foundFdp.fee_denom_params[0].min_base_fee).toEqual("0.001550000000000000")
        }
        catch (e: any) {
            expect(e).toBeFalsy();
        }
    });
});

const checkParams = (foundParams: any, params: Params) => {
    const exponent = BigNumber(10).exponentiatedBy(18);
    expect(BigNumber(foundParams.alpha).multipliedBy(exponent).isEqualTo(BigNumber(params.alpha))).toBe(true);
    expect(BigNumber(foundParams.beta).multipliedBy(exponent).isEqualTo(BigNumber(params.beta))).toBe(true);
    expect(BigNumber(foundParams.theta).multipliedBy(exponent).isEqualTo(BigNumber(params.theta))).toBe(true);
    expect(BigNumber(foundParams.min_learning_rate).multipliedBy(exponent).isEqualTo(BigNumber(params.minLearningRate))).toBe(true);
    expect(BigNumber(foundParams.max_learning_rate).multipliedBy(exponent).isEqualTo(BigNumber(params.maxLearningRate))).toBe(true);
    expect(foundParams.target_block_utilization).toBe(params.targetBlockUtilization);
    expect(foundParams.max_block_utilization).toBe(params.maxBlockUtilization);
    expect(foundParams.window).toBe(params.window);
    expect(foundParams.enabled).toBe(params.enabled);
    expect(foundParams.default_fee_denom).toBe(params.defaultFeeDenom);
}