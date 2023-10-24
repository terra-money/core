import { getLCDClient } from "../helpers/lcd.connection";

describe("ICQ Module (https://github.com/cosmos/ibc-apps/tree/main/modules/async-icq) ", () => {
    // Prepare environment clients, accounts and wallets
    const LCD = getLCDClient();

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
});
