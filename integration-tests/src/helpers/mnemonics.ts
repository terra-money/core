import { MnemonicKey } from "@terra-money/feather.js";

export function getMnemonics() {
    // Chain test-1
    let val1 = new MnemonicKey({
        mnemonic: "clock post desk civil pottery foster expand merit dash seminar song memory figure uniform spice circle try happy obvious trash crime hybrid hood cushion",
    });
    let rly1 = new MnemonicKey({
        mnemonic: "alley afraid soup fall idea toss can goose become valve initial strong forward bright dish figure check leopard decide warfare hub unusual join cart",
    })

    // Chain test-2
    let val2 = new MnemonicKey({
        mnemonic: "angry twist harsh drastic left brass behave host shove marriage fall update business leg direct reward object ugly security warm tuna model broccoli choice",
    })
    let rly2 = new MnemonicKey({
        mnemonic: "record gift you once hip style during joke field prize dust unique length more pencil transfer quit train device arrive energy sort steak upset",
    })

    // Funded wallets available in both chains
    let allianceMnemonic = new MnemonicKey({
        mnemonic: "broken title little open demand ladder mimic keen execute word couple door relief rule pulp demand believe cactus swing fluid tired what crop purse"
    })
    let pobMnemonic = new MnemonicKey({
        mnemonic: "banner spread envelope side kite person disagree path silver will brother under couch edit food venture squirrel civil budget number acquire point work mass"
    })
    let pobMnemonic1 = new MnemonicKey({
        mnemonic: "veteran try aware erosion drink dance decade comic dawn museum release episode original list ability owner size tuition surface ceiling depth seminar capable only"
    })
    let feeshareMnemonic = new MnemonicKey({
        mnemonic: "same heavy travel border destroy catalog music manual love festival exile resist always gas off coffee crystal provide random harvest sea cloud child field"
    })
    let genesisVesting = new MnemonicKey({
        mnemonic: "vacuum burst ordinary enact leaf rabbit gather lend left chase park action dish danger green jeans lucky dish mesh language collect acquire waste load"
    })
    let genesisVesting1 = new MnemonicKey({
        mnemonic: "open attitude harsh casino rent attitude midnight debris describe spare cancel crisp olive ride elite gallery leaf buffalo sheriff filter rotate path begin soldier"
    })
    let icaMnemonic = new MnemonicKey({
        mnemonic: "unit question bulk desk slush answer share bird earth brave book wing special gorilla ozone release permit mercy luxury version advice impact unfair drama"
    })
    let tokenFactoryMnemonic = new MnemonicKey({
        mnemonic: "year aim panel oyster sunny faint dress skin describe chair guilt possible venue pottery inflict mass debate poverty multiply pulse ability purse situate inmate"
    })
    let ibcHooksMnemonic = new MnemonicKey({
        mnemonic: "leave side blue panel curve ancient suspect slide seminar neutral doctor boring only curious spell surround remind obtain slogan hire giant soccer crunch system"
    })
    let wasmContracts = new MnemonicKey({
        mnemonic: "degree under tray object thought mercy mushroom captain bus work faint basic twice cube noble man ripple close flush bunker dish spare hungry arm"
    })
    let mnemonic2 = new MnemonicKey({
        mnemonic: "range struggle season mesh antenna delay sell light yard path risk curve brain nut cabin injury dilemma fun comfort crumble today transfer bring draft"
    })
    let mnemonic3 = new MnemonicKey({
        mnemonic: "giraffe trim element wheel cannon nothing enrich shiver upon output iron recall already fix appear produce fix behind scissors artefact excite tennis into side"
    })
    let mnemonic4 = new MnemonicKey({
        mnemonic: "run turn cup combine sad toast roof already melt chimney arctic save avocado theory bracket cherry cotton fee once favorite swarm ignore dream element"
    })
    let mnemonic5 = new MnemonicKey({
        mnemonic: "script key fold coyote cage squirrel prevent pole auction slide vintage shoot mirror erosion equip goose capable critic test space sketch monkey eight candy"
    })
    let mnemonic6 = new MnemonicKey({
        mnemonic: "work clap clarify edit explain exact depth ramp law hard feel beauty stumble occur prevent crush distance purpose scrap current describe skirt panther skirt"
    })

    return {
        val1,
        rly1,
        val2,
        rly2,
        allianceMnemonic,
        feeshareMnemonic,
        pobMnemonic,
        pobMnemonic1,
        genesisVesting,
        genesisVesting1,
        icaMnemonic,
        tokenFactoryMnemonic,
        ibcHooksMnemonic,
        wasmContracts,
        mnemonic2,
        mnemonic3,
        mnemonic4,
        mnemonic5,
        mnemonic6,
    }
}