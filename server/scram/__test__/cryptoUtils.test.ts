import {
  keysHMAC,
  signatureHMAC,
  hexStringToBytes,
  bytesToHexString,
  xorBytes,
} from "../cryptoUtils";

import { describe, test, expect } from "@jest/globals";
describe("Crypto Functions", () => {
  const password = "password";
  const salt = "1234567890abcdef1234567890abcdef";
  const iterationCount = 4096;
  const authMessage = "authMessage";
  const clientKey = "clientKey";

  // keysHMAC should return three keys
  test("keysHMAC should return three keys", () => {
    const { data } = keysHMAC(
      password,
      salt,
      iterationCount
    );
    expect(data.clientKey).toBeDefined();
    expect(data.serverKey).toBeDefined();
    expect(data.storedKey).toBeDefined();
    expect(data.clientKey.length).toBe(64); // SHA-256 produces a 64-character hex string
    expect(data.serverKey.length).toBe(64); // SHA-256 produces a 64-character hex string
    expect(data.storedKey.length).toBe(64); // SHA-256 produces a 64-character hex string
    expect(data.clientKey).toBe(
      "1d282febf2a3aa49c13c172fcf7dbb47fd1cc868332bf1d4edeb326f3c53d415"
    );
    expect(data.serverKey).toBe(
      "006cd21a24ef54c13dcece0dfa52de8d43871f24d3a7848bb0a136eed6ddeece"
    );
    expect(data.storedKey).toBe(
      "d8cde98fb85f1e12796adec01247a3a0fd39088e75b933a81cc6204fc1b1736a"
    );
  });

  test("SignatureHMAC should return a HMAC-SHA256 hash", () => {
    const signature = signatureHMAC(authMessage, clientKey);
    expect(signature).toBeDefined();
    expect(signature.length).toBe(64); // SHA-256 produces a 64-character hex string
  });

  test("hexStringToBytes should convert a hex string to a byte array", () => {
    const hex = "4a6f686e446f65";
    const bytes = hexStringToBytes(hex);
    expect(bytes).toEqual([74, 111, 104, 110, 68, 111, 101]);
  });

  test("bytesToHexString should convert a byte array to a hex string", () => {
    const bytes = [74, 111, 104, 110, 68, 111, 101];
    const hex = bytesToHexString(bytes);
    expect(hex).toBe("4a6f686e446f65");
  });

  test("xorBytes should perform a bitwise XOR operation on two byte arrays", () => {
    const bytesA = [0x0f, 0xf0, 0x55];
    const bytesB = [0xf0, 0x0f, 0xaa];
    const result = xorBytes(bytesA, bytesB);
    expect(result).toEqual([0xff, 0xff, 0xff]);
  });
});
