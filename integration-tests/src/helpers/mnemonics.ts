import { MnemonicKey } from "@terra-money/feather.js";

export function getMnemonics() {
    let val1 = new MnemonicKey({
        mnemonic: "clock post desk civil pottery foster expand merit dash seminar song memory figure uniform spice circle try happy obvious trash crime hybrid hood cushion",
    });
    let val2 = new MnemonicKey({
        mnemonic: "angry twist harsh drastic left brass behave host shove marriage fall update business leg direct reward object ugly security warm tuna model broccoli choice",
    })
    let mnemonic1 = new MnemonicKey({
        mnemonic: "banner spread envelope side kite person disagree path silver will brother under couch edit food venture squirrel civil budget number acquire point work mass",
    })
    let mnemonic11 = new MnemonicKey({
        mnemonic: "vacuum burst ordinary enact leaf rabbit gather lend left chase park action dish danger green jeans lucky dish mesh language collect acquire waste load",
    })
    let mnemonic2 = new MnemonicKey({
        mnemonic: "veteran try aware erosion drink dance decade comic dawn museum release episode original list ability owner size tuition surface ceiling depth seminar capable only",
    })
    let mnemonic22 = new MnemonicKey({
        mnemonic: "open attitude harsh casino rent attitude midnight debris describe spare cancel crisp olive ride elite gallery leaf buffalo sheriff filter rotate path begin soldier",
    })
    let rly1 = new MnemonicKey({
        mnemonic: "alley afraid soup fall idea toss can goose become valve initial strong forward bright dish figure check leopard decide warfare hub unusual join cart",
    })
    let rly2 = new MnemonicKey({
        mnemonic: "record gift you once hip style during joke field prize dust unique length more pencil transfer quit train device arrive energy sort steak upset",
    })

    return {
        val1,
        val2,
        mnemonic1,
        mnemonic11,
        mnemonic2,
        mnemonic22,
        rly1,
        rly2,
    }
}