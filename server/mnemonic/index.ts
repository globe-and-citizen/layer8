import * as bip39 from "@scure/bip39";
import { wordlist } from '@scure/bip39/wordlists/english';

import * as bip32 from "@scure/bip32"

export function generateBip39Mnemonic(): string {
    return bip39.generateMnemonic(wordlist, 128)
}

export function isMnemonicValid(mnemonic: string): boolean {
    return bip39.validateMnemonic(mnemonic, wordlist)
}

export function getBinarySeed(mnemonic: string): Uint8Array {
    return bip39.mnemonicToSeedSync(mnemonic)
}

export function getPrivateAndPublicKeys(binarySeed: Uint8Array) {
    // generate a master key
    const hdKey = bip32.HDKey.fromMasterSeed(binarySeed)
    return {
        "privateKey": hdKey.privateKey,
        "publicKey": hdKey.publicKey
    }
}
