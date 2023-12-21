import { LCDClient } from "@terra-money/feather.js";

export class LCDClients {
    public chain1 = new LCDClient({
        "test-1": {
            lcd: "http://localhost:1316",
            chainID: "test-1",
            gasPrices: "0.15uluna",
            gasAdjustment: 1.5,
            prefix: "terra"
        }
    })
    public chain2 = new LCDClient({
        "test-2": {
            lcd: "http://localhost:1317",
            chainID: "test-2",
            gasPrices: "0.15uluna",
            gasAdjustment: 1.5,
            prefix: "terra"
        }
    })

    static create() {
        return new LCDClients();
    }

    private constructor() { }

    async blockInclusionChain1() {
        let res = await this.chain1.tendermint.blockInfo("test-1")
        let height = res.block.header.height;
        let currentHeight = res.block.header.height;

        for await (const _ of new Array(10)) {
            await interval();
            let res = await this.chain1.tendermint.blockInfo("test-1")
            currentHeight = res.block.header.height;

            if (height != currentHeight) return Promise.resolve();
        }
    }

    async blockInclusionChain2() {
        let res = await this.chain2.tendermint.blockInfo("test-2")
        let height = res.block.header.height;
        let currentHeight = res.block.header.height;

        for await (const _ of new Array(10)) {
            await interval();
            let res = await this.chain2.tendermint.blockInfo("test-2")
            currentHeight = res.block.header.height;

            if (height != currentHeight) return Promise.resolve();
        }
    }
}

function interval(): Promise<void>{
    return new Promise((resolve) => {
        setTimeout(() => {
          resolve();
        }, 400);
      });
}
