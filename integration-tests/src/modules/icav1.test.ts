//import { Coin, MsgTransfer } from "@terra-money/feather.js";
import { blockInclusion, getLCDClient, getMnemonics } from "../helpers";
//import { Height } from "@terra-money/feather.js/dist/core/ibc/core/client/Height";
import { MsgRegisterInterchainAccount } from "@terra-money/feather.js/dist/core/ica/controller/v1/msgs";

describe("ICA Module (https://github.com/cosmos/ibc-go/tree/release/v7.3.x/modules/apps/27-interchain-accounts)", () => {
    // Prepare environment clients, accounts and wallets
    const LCD = getLCDClient();
    const { icaMnemonic } = getMnemonics();
    const chain1Wallet = LCD.chain1.wallet(icaMnemonic);
    const externalAccAddr = icaMnemonic.accAddress("terra");

    test('Must contain the expected module params', async () => {
        // Query ica host module params
        const hostResParams = await LCD.chain2.icaV1.hostParams("test-2");
        expect(hostResParams.params)
            .toStrictEqual({
                "host_enabled": true,
                "allow_messages": [
                    "/cosmos.authz.v1beta1.MsgExec",
                    "/cosmos.authz.v1beta1.MsgGrant",
                    "/cosmos.authz.v1beta1.MsgRevoke",
                    "/cosmos.bank.v1beta1.MsgSend",
                    "/cosmos.bank.v1beta1.MsgMultiSend",
                    "/cosmos.distribution.v1beta1.MsgSetWithdrawAddress",
                    "/cosmos.distribution.v1beta1.MsgWithdrawValidatorCommission",
                    "/cosmos.distribution.v1beta1.MsgFundCommunityPool",
                    "/cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward",
                    "/cosmos.feegrant.v1beta1.MsgGrantAllowance",
                    "/cosmos.feegrant.v1beta1.MsgRevokeAllowance",
                    "/cosmos.gov.v1beta1.MsgVoteWeighted",
                    "/cosmos.gov.v1beta1.MsgSubmitProposal",
                    "/cosmos.gov.v1beta1.MsgDeposit",
                    "/cosmos.gov.v1beta1.MsgVote",
                    "/cosmos.staking.v1beta1.MsgEditValidator",
                    "/cosmos.staking.v1beta1.MsgDelegate",
                    "/cosmos.staking.v1beta1.MsgUndelegate",
                    "/cosmos.staking.v1beta1.MsgBeginRedelegate",
                    "/cosmos.staking.v1beta1.MsgCreateValidator",
                    "/cosmos.vesting.v1beta1.MsgCreateVestingAccount",
                    "/ibc.applications.transfer.v1.MsgTransfer",
                    "/cosmwasm.wasm.v1.MsgStoreCode",
                    "/cosmwasm.wasm.v1.MsgInstantiateContract",
                    "/cosmwasm.wasm.v1.MsgExecuteContract",
                    "/cosmwasm.wasm.v1.MsgMigrateContract"
                ]
            });

        // Query contoller module params
        const controllerResParams = await LCD.chain2.icaV1.controllerParams("test-2");
        expect(controllerResParams.params)
            .toStrictEqual({
                controller_enabled: true,
            });
    });

    test('Must create an interchain account from chain1 to chain2', async () => {
        try {

        let tx = await chain1Wallet.createAndSignTx({
            msgs: [new MsgRegisterInterchainAccount(
                externalAccAddr,
                "connection-0",
                ""
            )],
            chainID: "test-1",
        });
        console.log(tx);
        let result = await LCD.chain1.tx.broadcastSync(tx, "test-1");
        console.log("result",JSON.stringify(result));
        await blockInclusion();
        let txResult = await LCD.chain1.tx.txInfo(result.txhash, "test-1") as any;
        expect(txResult).toBeDefined();
        console.log("txResult",JSON.stringify(txResult));

        let res = await LCD.chain1.icaV1.controllerAccountAddress(externalAccAddr, "connection-0");
        console.log("Res",JSON.stringify(res))
        }
        catch(e) {
            console.log("Error",e)
            expect(e).toBeUndefined();
        }
    });
});
