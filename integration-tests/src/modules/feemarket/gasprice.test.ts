// import { Params as Params_pb } from "@terra-money/terra.proto/feemarket/feemarket/v1/params"
import BigNumber from 'bignumber.js';
import { getLCDClient } from "../../helpers";

describe("Feemarket Module dynamic fees (https://github.com/terra-money/feemarket/tree/v0.0.1-alpha.2-terra.0) ", () => {
    // Prepare environment clients, accounts and wallets
    const LCD = getLCDClient();
    test('Check gas price increases and decrease back to 0.0015', async () => {
        try {
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
});
