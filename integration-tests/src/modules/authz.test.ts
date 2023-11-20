import { getMnemonics } from "../helpers/mnemonics";
import { getLCDClient } from "../helpers/lcd.connection";
import { StakeAuthorization, MsgGrantAuthorization, AuthorizationGrant, Coin, MsgExecAuthorized, MsgDelegate } from "@terra-money/feather.js";
import { AuthorizationType } from "@terra-money/terra.proto/cosmos/staking/v1beta1/authz";
import moment from "moment";
import { blockInclusion } from "../helpers/const";

describe("Authz Module (https://github.com/terra-money/cosmos-sdk/tree/release/v0.47.x/x/authz)", () => {
    const LCD = getLCDClient();
    const accounts = getMnemonics();
    // Accounts used in chain2, which means that 
    // will not cause conflicts with txs nonces
    const granterWallet = LCD.chain2.wallet(accounts.feeshareMnemonic);
    const granteeWallet = LCD.chain2.wallet(accounts.pobMnemonic);
    const granterAddr = accounts.feeshareMnemonic.accAddress("terra");
    const granteeAddr = accounts.pobMnemonic.accAddress("terra");
    const val2Addr = accounts.val2.valAddress("terra");

    test('Must register the granter', async () => {
        let tx = await granterWallet.createAndSignTx({
            msgs: [new MsgDelegate(
                granterAddr,
                val2Addr,
                Coin.fromString("1000000uluna"),
            ),new MsgGrantAuthorization(
                granterAddr,
                granteeAddr,
                new AuthorizationGrant(
                    new StakeAuthorization(
                        AuthorizationType.AUTHORIZATION_TYPE_DELEGATE,
                        Coin.fromString("1000000uluna"),
                    ),
                    moment().add(1, "hour").toDate(),
                ),
            )],
            chainID: "test-2",
        });
        let result = await LCD.chain2.tx.broadcastSync(tx, "test-2");
        await blockInclusion();

        // Check the MsgGrantAuthorization executed as expected 
        let txResult = await LCD.chain2.tx.txInfo(result.txhash, "test-2") as any;
        expect(txResult.logs[0].events)
            .toStrictEqual([{
                "type": "message",
                "attributes": [{
                    "key": "action",
                    "value": "/cosmos.authz.v1beta1.MsgGrant"
                }, {
                    "key": "sender",
                    "value": "terra120rzk7n6cd2vufkmwrat34adqh0rgca9tkyfe5"
                }, {
                    "key": "module",
                    "value": "authz"
                }]
            }, {
                "type": "cosmos.authz.v1beta1.EventGrant",
                "attributes": [{
                    "key": "grantee",
                    "value": "\"terra1v0eee20gjl68fuk0chyrkch2z7suw2mhg3wkxf\""
                }, {
                    "key": "granter",
                    "value": "\"terra120rzk7n6cd2vufkmwrat34adqh0rgca9tkyfe5\""
                }, {
                    "key": "msg_type_url",
                    "value": "\"/cosmos.staking.v1beta1.MsgDelegate\""
                }]
            }]);
    });

    describe("Grantee must execute", () => {
        test("delegation on belhalf of granter", async () => {
            try {
                let tx = await granteeWallet.createAndSignTx({
                    msgs: [new MsgExecAuthorized(
                        granteeAddr,
                        [new MsgDelegate(
                            granterAddr,
                            val2Addr,
                            Coin.fromString("1000000uluna"),
                        )]
                    )],
                    chainID: "test-2",
                });
                let result = await LCD.chain2.tx.broadcastSync(tx, "test-2");
                await blockInclusion();

                console.log(result);
            }
            catch (e) {
                console.log(e)
            }
        });
    })
});