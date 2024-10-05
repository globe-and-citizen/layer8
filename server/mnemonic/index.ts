import * as bip39 from "@scure/bip39";
import { wordlist } from '@scure/bip39/wordlists/english';

import * as bip32 from "@scure/bip32"

export function generateMnemonic(): string {
    return bip39.generateMnemonic(wordlist, 128)
}

export function validateMnemonic(mnemonic: string): boolean {
    return bip39.validateMnemonic(mnemonic, wordlist)
}

export function getBinarySeed(mnemonic: string): Uint8Array {
    return bip39.mnemonicToSeedSync(mnemonic)
}

export function getPrivateAndPublicKeys(binarySeed: Uint8Array) {
    // generate a master key
    return bip32.HDKey.fromMasterSeed(binarySeed)
}

function test() {
    const mnemonic = generateMnemonic()
    console.log(mnemonic)

    const isValid = validateMnemonic(mnemonic)
    console.log(isValid)

    const seed = getBinarySeed(mnemonic)
    console.log(seed)

    const hdKey = getPrivateAndPublicKeys(seed)

    console.log("Private key:")
    console.log(hdKey.privateKey)

    console.log("Public key:")
    console.log(hdKey.publicKey)
}
