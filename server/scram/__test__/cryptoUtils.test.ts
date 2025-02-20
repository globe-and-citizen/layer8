import {
  saltAndHashPassword,
  clientKeyHMAC,
  serverKeyHMAC,
  storedKeySHA256,
  SignatureHMAC,
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

  test("saltAndHashPassword should return a hashed password", () => {
    const hashedPassword = saltAndHashPassword(password, salt, iterationCount);
    expect(hashedPassword).toBeDefined();
    expect(hashedPassword.length).toBe(40); // SHA-1 produces a 40-character hex string
  });

  // test saltAndHashedPassword with special characters
  test("saltAndHashPassword should return a hashed password with special characters", () => {
    const hashedPassword = saltAndHashPassword(
      "password!@#$",
      salt,
      iterationCount
    );
    expect(hashedPassword).toBeDefined();
    expect(hashedPassword.length).toBe(40);
  });

  // test saltAndHashedPassword with a very long password
  test("saltAndHashPassword should return a hashed password with a very long password", () => {
    const hashedPassword = saltAndHashPassword(
      "ThisIsaVeryLongStringPasswordWithNumbers1234567890AndSpecialCharacters!@#$%^&*()",
      salt,
      iterationCount
    );
    expect(hashedPassword).toBeDefined();
    expect(hashedPassword.length).toBe(40);
  });

  test("clientKeyHMAC should return a HMAC-SHA256 hash", () => {
    const saltedPassword = saltAndHashPassword(password, salt, iterationCount);
    const hmac = clientKeyHMAC(saltedPassword);
    expect(hmac).toBeDefined();
    expect(hmac.length).toBe(64); // SHA-256 produces a 64-character hex string
  });

  test("serverKeyHMAC should return a HMAC-SHA256 hash", () => {
    const saltedPassword = saltAndHashPassword(password, salt, iterationCount);
    const hmac = serverKeyHMAC(saltedPassword);
    expect(hmac).toBeDefined();
    expect(hmac.length).toBe(64); // SHA-256 produces a 64-character hex string
  });

  test("storedKeySHA256 should return a SHA-256 hash", () => {
    const hash = storedKeySHA256(clientKey);
    expect(hash).toBeDefined();
    expect(hash.length).toBe(64); // SHA-256 produces a 64-character hex string
  });

  test("SignatureHMAC should return a HMAC-SHA256 hash", () => {
    const signature = SignatureHMAC(authMessage, clientKey);
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
