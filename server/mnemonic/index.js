"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.generateBip39Mnemonic = generateBip39Mnemonic;
exports.isMnemonicValid = isMnemonicValid;
exports.getBinarySeed = getBinarySeed;
exports.getPrivateAndPublicKeys = getPrivateAndPublicKeys;
var bip39 = require("@scure/bip39");
var english_1 = require("@scure/bip39/wordlists/english");
var bip32 = require("@scure/bip32");
function generateBip39Mnemonic() {
    return bip39.generateMnemonic(english_1.wordlist, 128);
}
function isMnemonicValid(mnemonic) {
    return bip39.validateMnemonic(mnemonic, english_1.wordlist);
}
function getBinarySeed(mnemonic) {
    return bip39.mnemonicToSeedSync(mnemonic);
}
function getPrivateAndPublicKeys(binarySeed) {
    // generate a master key
    var hdKey = bip32.HDKey.fromMasterSeed(binarySeed);
    return {
        "privateKey": hdKey.privateKey,
        "publicKey": hdKey.publicKey
    };
}
