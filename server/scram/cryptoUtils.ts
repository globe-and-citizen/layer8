"use strict";
import * as CryptoJS from "crypto-js";

export function saltAndHashPassword(password: string, salt: string, iterationCount: number): string {
    return CryptoJS.PBKDF2(
        password,
        CryptoJS.enc.Hex.parse(salt),
        {
            keySize: 160 / 32,
            iterations: iterationCount,
            hasher: CryptoJS.algo.SHA1,
        }
    ).toString(CryptoJS.enc.Hex);
}

export function clientKeyHMAC(saltedPassword: string): string {
    return CryptoJS.HmacSHA256(
        CryptoJS.enc.Hex.parse(saltedPassword),
        CryptoJS.enc.Utf8.parse("Client Key")
    ).toString(CryptoJS.enc.Hex);
}

export function serverKeyHMAC(saltedPassword: string): string {
    return CryptoJS.HmacSHA256(
        CryptoJS.enc.Hex.parse(saltedPassword),
        CryptoJS.enc.Utf8.parse("Server Key")
    ).toString(CryptoJS.enc.Hex);
}

export function storedKeySHA256(clientKey: string): string {
    return CryptoJS.SHA256(
        CryptoJS.enc.Hex.parse(clientKey)
    ).toString(CryptoJS.enc.Hex);
}

// Can be used for both, clientSignature and serverSignature
export function SignatureHMAC(authMessage: string, key: string): string {
    return CryptoJS.HmacSHA256(
        CryptoJS.enc.Utf8.parse(authMessage),
        CryptoJS.enc.Hex.parse(key)
    ).toString(CryptoJS.enc.Hex);
}

export function hexStringToBytes(hex: string): number[] {
    const bytes: number[] = [];
    for (let i = 0; i < hex.length; i += 2) {
        bytes.push(parseInt(hex.substr(i, 2), 16));
    }
    return bytes;
}

export function bytesToHexString(bytes: number[]): string {
    return bytes.map((byte) => byte.toString(16).padStart(2, "0")).join("");
}

export function xorBytes(bytesA: number[], bytesB: number[]): number[] {
    return bytesA.map((byte, index) => byte ^ bytesB[index]);
}