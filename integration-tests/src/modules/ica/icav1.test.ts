import { AccAddress, Coin, MsgTransfer, MsgSend, Coins } from "@terra-money/feather.js";
import { blockInclusion, getLCDClient, getMnemonics } from "../../helpers";
import { MsgRegisterInterchainAccount, MsgSendTx } from "@terra-money/feather.js/dist/core/ica/controller/v1/msgs";
import { Height } from "@terra-money/feather.js/dist/core/ibc/core/client/Height";
import Long from "long";
import { InterchainAccountPacketData } from "@terra-money/feather.js/dist/core/ica/controller/v1/InterchainAccountPacketData";
import { CosmosTx } from "@terra-money/feather.js/dist/core/ica/controller/v1/CosmosTx";

describe("ICA Module (https://github.com/cosmos/ibc-go/tree/release/v7.3.x/modules/apps/27-interchain-accounts)", () => {
    // Prepare environment clients, accounts and wallets
    const LCD = getLCDClient();
    const { icaMnemonic } = getMnemonics();
    const chain1Wallet = LCD.chain1.wallet(icaMnemonic);
    const externalAccAddr = icaMnemonic.accAddress("terra");
    let ibcCoinDenom: string | undefined;
    let intechainAccountAddr: AccAddress | undefined;

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

    test('Must query the interchain account to determine its existance', async () => {
        // Query the account address of the interchain account
        let res = await LCD.chain1.icaV1.controllerAccountAddress(externalAccAddr, "connection-0")
            .catch(e => {
                // assert that the expected error is that it failed to retreive the account
                const expectMsg = "failed to retrieve account address for icacontroller-";
                expect(e.response.data.message.startsWith(expectMsg)).toBeTruthy();
            })

        // if res is defined then the account exists
        if (res !== undefined) {
            expect(res.address).toBeDefined();
            intechainAccountAddr = res.address;
            // Check during 5 blocks for the receival 
            // of the IBC coin on chain-2
            for (let i = 0; i <= 5; i++) {
                await blockInclusion();
                let _ibcCoin = (await LCD.chain2.bank.balance(intechainAccountAddr))[0].find(c => c.denom.startsWith("ibc/"));
                if (_ibcCoin) {
                    expect(_ibcCoin.denom.startsWith("ibc/")).toBeTruthy();
                    ibcCoinDenom = _ibcCoin.denom
                    break;
                }
            }
        }
    });

    test('Must creat the interchain account if des not already exist', async () => {
        let tx = await chain1Wallet.createAndSignTx({
            msgs: [new MsgRegisterInterchainAccount(
                externalAccAddr,
                "connection-0",
                ""
            )],
            chainID: "test-1",
        }).catch(e => {
            console.log(e)
            const expectedMsg = "failed to execute message; message index: 0: existing active channel channel-1 for portID icacontroller-terra1p4kcrttuxj9kyyvv5px5ccgwf0yrw74yp7jqm6 on connection connection-0: active channel already set for this owner";
            expect(e.response.data.message.startsWith(expectedMsg))
                .toBeTruthy();
        });

        if (tx !== undefined) {
            let result = await LCD.chain1.tx.broadcastSync(tx, "test-1");
            await blockInclusion();
            let txResult = await LCD.chain1.tx.txInfo(result.txhash, "test-1") as any;
            let events = txResult.logs[0].events;

            expect(events[0])
                .toStrictEqual({
                    "type": "message",
                    "attributes": [{
                        "key": "action",
                        "value": "/ibc.applications.interchain_accounts.controller.v1.MsgRegisterInterchainAccount"
                    }, {
                        "key": "sender",
                        "value": "terra1p4kcrttuxj9kyyvv5px5ccgwf0yrw74yp7jqm6"
                    }]
                });
            expect(events[2])
                .toStrictEqual({
                    "type": "message",
                    "attributes": [{
                        "key": "module",
                        "value": "ibc_channel"
                    }]
                })

            // Check during 5 blocks for the receival 
            // of the IBC coin on chain-2
            for (let i = 0; i <= 5; i++) {
                await blockInclusion();
                let res = await LCD.chain1.icaV1.controllerAccountAddress(externalAccAddr, "connection-0")
                    .catch((e) => {
                        const expectMsg = "failed to retrieve account address for icacontroller-";
                        expect(e.response.data.message.startsWith(expectMsg)).toBeTruthy();
                    })
                if (res) {
                    expect(res.address).toBeDefined();
                    intechainAccountAddr = res.address;
                    break;
                }
            }
        }
    });

    describe('After assuring the interchain account exists', () => {
        test("Must send funds to the interchain account from chain-1 to chain-2", async () => {
            if (typeof intechainAccountAddr === "string") {
                let blockHeight = (await LCD.chain1.tendermint.blockInfo("test-1")).block.header.height;
                let tx = await chain1Wallet.createAndSignTx({
                    msgs: [new MsgTransfer(
                        "transfer",
                        "channel-0",
                        Coin.fromString("100000000uluna"),
                        externalAccAddr,
                        intechainAccountAddr as string,
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
                    let _ibcCoin = (await LCD.chain2.bank.balance(intechainAccountAddr))[0].find(c => c.denom.startsWith("ibc/"));
                    if (_ibcCoin) {
                        expect(_ibcCoin.denom.startsWith("ibc/")).toBeTruthy();
                        ibcCoinDenom = _ibcCoin.denom
                        break;
                    }
                }
            } else {
                // This case should never happen but if something goes wrong
                // this is a check to fail.
                expect(intechainAccountAddr).toBeDefined()
            }
        });

        test("Must control the interchain account from chain-1 to send funds on chain-2 from the account address to a burnAddress", async () => {
            try {

                const burnAddress = "terra1zdpgj8am5nqqvht927k3etljyl6a52kwqup0je";
                let interchainAccountPacketData = new InterchainAccountPacketData(
                    new CosmosTx([new MsgSend(
                        intechainAccountAddr as string,
                        burnAddress,
                        Coins.fromString("1000" + ibcCoinDenom),
                    )])
                )
                let msgSendTx = new MsgSendTx(
                    externalAccAddr,
                    "connection-0",
                    Long.fromString((new Date().getTime() * 1000000 + 600000000).toString()),
                    interchainAccountPacketData,
                );
                let tx = await chain1Wallet.createAndSignTx({
                    msgs: [msgSendTx],
                    chainID: "test-1",
                });
    
                let result = await LCD.chain1.tx.broadcastSync(tx, "test-1");
                await blockInclusion();
                let txResult = await LCD.chain1.tx.txInfo(result.txhash, "test-1") as any;
                const events = txResult.logs[0].events;
                expect(events[0])
                    .toStrictEqual({
                        "type": "message",
                        "attributes": [{
                            "key": "action",
                            "value": "/ibc.applications.interchain_accounts.controller.v1.MsgSendTx"
                        }, {
                            "key": "sender",
                            "value": "terra1p4kcrttuxj9kyyvv5px5ccgwf0yrw74yp7jqm6"
                        }]
                    });
                expect(events[2])
                    .toStrictEqual({
                        "type": "message",
                        "attributes": [{
                            "key": "module",
                            "value": "ibc_channel"
                        }]
                    })
    
    
                // Check during 5 blocks for the receival 
                // of the IBC coin on chain-2
                for (let i = 0; i <= 5; i++) {
                    await blockInclusion();
                    const bankRes = await LCD.chain2.bank.balance(burnAddress);
                    const coins = bankRes[0].find(c => c.denom === ibcCoinDenom);
                    if (coins) {
                        expect(coins).toBeDefined();
                        expect(coins?.denom).toStrictEqual(ibcCoinDenom);
                        expect(coins?.amount.toNumber()).toBeGreaterThanOrEqual(1000);
                        break;
                    }
                }
            }
            catch(e) {
                console.log(e)
            }
        })
    });
});
