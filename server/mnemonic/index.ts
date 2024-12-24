import * as bip39 from "@scure/bip39";
import { wordlist } from '@scure/bip39/wordlists/english';

import * as bip32 from "@scure/bip32"

import { secp256k1 } from '@noble/curves/secp256k1';
import { keccak_256 } from "@noble/hashes/sha3"

export function generateBip39Mnemonic(): string {
    return bip39.generateMnemonic(wordlist, 128)
}

export function isValid(mnemonic: string): boolean {
    return bip39.validateMnemonic(mnemonic, wordlist)
}

function getBinarySeed(mnemonic: string): Uint8Array {
    return bip39.mnemonicToSeedSync(mnemonic)
}

export interface KeyPair {
    privateKey: Uint8Array;
    publicKey: Uint8Array
}

export function getPrivateAndPublicKeys(mnemonic: string): KeyPair {
    const binarySeed = getBinarySeed(mnemonic)
    // generate a master key
    const hdKey = bip32.HDKey.fromMasterSeed(binarySeed)
    return {
        privateKey: hdKey.privateKey,
        publicKey: hdKey.publicKey
    }
}

export function sign(privateKey: Uint8Array, message: string): Uint8Array {
    const msgHash = keccak_256(message)
    return secp256k1.sign(msgHash, privateKey).toCompactRawBytes()
}
