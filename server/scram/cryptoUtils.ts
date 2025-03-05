"use strict";
import * as CryptoJS from "crypto-js";

// clientAndServerKeyHMAC will return three keys
export function keysHMAC(
  password: string,
  salt: string,
  iterationCount: number
): string[] {
  const hashedPassword = CryptoJS.PBKDF2(
    password,
    CryptoJS.enc.Hex.parse(salt),
    {
      keySize: 160 / 32,
      iterations: iterationCount,
      hasher: CryptoJS.algo.SHA1,
    }
  ).toString(CryptoJS.enc.Hex);

  const clientKey = CryptoJS.HmacSHA256(
    CryptoJS.enc.Hex.parse(hashedPassword),
    CryptoJS.enc.Utf8.parse("Client Key")
  ).toString(CryptoJS.enc.Hex);

  const serverKey = CryptoJS.HmacSHA256(
    CryptoJS.enc.Hex.parse(hashedPassword),
    CryptoJS.enc.Utf8.parse("Server Key")
  ).toString(CryptoJS.enc.Hex);

  const storedKey = CryptoJS.SHA256(CryptoJS.enc.Hex.parse(clientKey)).toString(
    CryptoJS.enc.Hex
  );

  return [storedKey, serverKey, clientKey];
}

// Can be used for both, clientSignature and serverSignature
export function signatureHMAC(authMessage: string, key: string): string {
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
