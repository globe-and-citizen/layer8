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

const Username = "test_user";
const Password = "test_password_123+!$&";
const salt = "1234567890abcdef1234567890abcdef";
const iterationCount = 4096;
const cNonce = "42ZHCYjtfdaykloIJFvxOLlYZ0MGMLtqwKKDnfUjfC0=";
const nonce =
  "42ZHCYjtfdaykloIJFvxOLlYZ0MGMLtqwKKDnfUjfC0=hFeC5hCyjefM0/6enYoahN+IRrqRr/VAi14ueGKFSKI=";
const hasedPassword = "8bfd7458a68206ee3841a820f43942f60f9a9ec2";
const clientKey =
  "41bf61eb5b6f29c517176770ef609cecac75559d84c292a1f0b162fa219f360f";
const serverKey =
  "1c0a646578ebe5b2f1ac20e88307a4b3c2845457c28c9fda3ce9654b05803669";
const storedKey =
  "a13a4e5eecf4c6ad49560c3a640a8d541a66fca5c16c90b9d38aa19b552c3257";
const clientKeyBytes = [
  65, 191, 97, 235, 91, 111, 41, 197, 23, 23, 103, 112, 239, 96, 156, 236, 172,
  117, 85, 157, 132, 194, 146, 161, 240, 177, 98, 250, 33, 159, 54, 15,
];
const clientSignatureBytes = [
  138, 103, 130, 183, 112, 217, 112, 226, 185, 68, 81, 225, 13, 97, 138, 53, 81,
  202, 108, 154, 160, 243, 7, 246, 215, 202, 80, 124, 171, 201, 231, 217,
];
const clientProofBytes = [
  203, 216, 227, 92, 43, 182, 89, 39, 174, 83, 54, 145, 226, 1, 22, 217, 253,
  191, 57, 7, 36, 49, 149, 87, 39, 123, 50, 134, 138, 86, 209, 214,
];
const authMessage = `[n=${Username},r=${cNonce},s=${salt},i=${iterationCount},r=${nonce}]`;
const clientSignature =
  "8a6782b770d970e2b94451e10d618a3551ca6c9aa0f307f6d7ca507cabc9e7d9";
const clientProof =
  "cbd8e35c2bb65927ae533691e20116d9fdbf390724319557277b32868a56d1d6";

describe("E2E Tests for Register Flow", () => {
  test("saltAndHashPassword should return a hashed password", () => {
    const generatedHashedPassword = saltAndHashPassword(
      Password,
      salt,
      iterationCount
    );
    expect(generatedHashedPassword).toBeDefined();
    expect(generatedHashedPassword.length).toBe(40); // SHA-1 produces a 40-character hex string
    // Expect hasedPassword to be equal to `8bfd7458a68206ee3841a820f43942f60f9a9ec2` as per the test
    expect(generatedHashedPassword).toBe(hasedPassword);
  });

  test("clientKeyHMAC should return a HMAC-SHA256 hash", () => {
    const generatedClientKey = clientKeyHMAC(hasedPassword);
    expect(generatedClientKey).toBeDefined();
    expect(generatedClientKey.length).toBe(64); // SHA-256 produces a 64-character hex string
    // Expect clientKey to be equal to `41bf61eb5b6f29c517176770ef609cecac75559d84c292a1f0b162fa219f360f` as per the test
    expect(generatedClientKey).toBe(clientKey);
  });

  test("serverKeyHMAC should return a HMAC-SHA256 hash", () => {
    const generatedServerKey = serverKeyHMAC(hasedPassword);
    expect(generatedServerKey).toBeDefined();
    expect(generatedServerKey.length).toBe(64); // SHA-256 produces a 64-character hex string
    // Expect serverKey to be equal to `1c0a646578ebe5b2f1ac20e88307a4b3c2845457c28c9fda3ce9654b05803669` as per the test
    expect(generatedServerKey).toBe(serverKey);
  });

  test("storedKeySHA256 should return a SHA-256 hash", () => {
    const generatedStoredKey = storedKeySHA256(clientKey);
    expect(generatedStoredKey).toBeDefined();
    expect(generatedStoredKey.length).toBe(64); // SHA-256 produces a 64-character hex string
    // Expect storedKey to be equal to `a13a4e5eecf4c6ad49560c3a640a8d541a66fca5c16c90b9d38aa19b552c3257` as per the test
    expect(generatedStoredKey).toBe(storedKey);
  });
});

describe("E2E Tests for Login Flow", () => {
  test("saltAndHashPassword should return a hashed password", () => {
    const generatedHashedPassword = saltAndHashPassword(
      Password,
      salt,
      iterationCount
    );
    expect(generatedHashedPassword).toBeDefined();
    expect(generatedHashedPassword.length).toBe(40); // SHA-1 produces a 40-character hex string
    // Expect hasedPassword to be equal to `8bfd7458a68206ee3841a820f43942f60f9a9ec2` as per the test
    expect(generatedHashedPassword).toBe(hasedPassword);
  });

  test("clientKeyHMAC should return a HMAC-SHA256 hash", () => {
    const generatedClientKey = clientKeyHMAC(hasedPassword);
    expect(generatedClientKey).toBeDefined();
    expect(generatedClientKey.length).toBe(64); // SHA-256 produces a 64-character hex string
    // Expect clientKey to be equal to `41bf61eb5b6f29c517176770ef609cecac75559d84c292a1f0b162fa219f360f` as per the test
    expect(generatedClientKey).toBe(clientKey);
  });

  test("hexStringToBytes should convert clientKey to a byte array", () => {
    const generatedClientKeyBytes = hexStringToBytes(clientKey);
    expect(generatedClientKeyBytes).toBeDefined();
    expect(generatedClientKeyBytes.length).toBe(32); // It should be a 32-byte array
    // It should match the clientKeyBytes array as per the test
    expect(generatedClientKeyBytes).toStrictEqual(clientKeyBytes);
  });

  test("storedKeySHA256 should return a SHA-256 hash", () => {
    const generatedStoredKey = storedKeySHA256(clientKey);
    expect(generatedStoredKey).toBeDefined();
    expect(generatedStoredKey.length).toBe(64); // SHA-256 produces a 64-character hex string
    // Expect storedKey to be equal to `a13a4e5eecf4c6ad49560c3a640a8d541a66fca5c16c90b9d38aa19b552c3257` as per the test
    expect(generatedStoredKey).toBe(storedKey);
  });

  test("serverKeyHMAC should return a HMAC-SHA256 hash", () => {
    const generatedServerKey = serverKeyHMAC(hasedPassword);
    expect(generatedServerKey).toBeDefined();
    expect(generatedServerKey.length).toBe(64); // SHA-256 produces a 64-character hex string
    // Expect serverKey to be equal to `1c0a646578ebe5b2f1ac20e88307a4b3c2845457c28c9fda3ce9654b05803669` as per the test
    expect(generatedServerKey).toBe(serverKey);
  });

  test("SignatureHMAC should return a HMAC-SHA256 hash", () => {
    const generatedClientSignature = SignatureHMAC(authMessage, storedKey);
    expect(generatedClientSignature).toBeDefined();
    expect(generatedClientSignature.length).toBe(64); // SHA-256 produces a 64-character hex string
    // Expect generatedClientSignature to be equal to `8a6782b770d970e2b94451e10d618a3551ca6c9aa0f307f6d7ca507cabc9e7d9` as per the test
    expect(generatedClientSignature).toBe(clientSignature);
  });

  test("hexStringToBytes should convert clientSignature to a byte array", () => {
    const generatedClientSignatureBytes = hexStringToBytes(clientSignature);
    expect(generatedClientSignatureBytes).toBeDefined();
    expect(generatedClientSignatureBytes.length).toBe(32); // It should be a 32-byte array
    // It should match the clientSignatureBytes array as per the test
    expect(generatedClientSignatureBytes).toStrictEqual(clientSignatureBytes);
  });

  test("xorBytes should perform a bitwise XOR operation on clientKeyBytes and clientSignatureBytes", () => {
    const generatedClientProofBytes = xorBytes(
      clientKeyBytes,
      clientSignatureBytes
    );
    expect(generatedClientProofBytes).toBeDefined();
    expect(generatedClientProofBytes.length).toBe(32); // It should be a 32-byte array
    // It should match the clientProofBytes array as per the test
    expect(generatedClientProofBytes).toStrictEqual(clientProofBytes);
  });

  test("bytesToHexString should convert clientProofBytes to a hex string", () => {
    const generatedClientProof = bytesToHexString(clientProofBytes);
    expect(generatedClientProof).toBeDefined();
    expect(generatedClientProof.length).toBe(64); // It should be a 64-character hex string
    // Expect generatedClientProof to be equal to `cbd8e35c2bb65927ae533691e20116d9fdbf390724319557277b32868a56d1d6` as per the test
    expect(generatedClientProof).toBe(clientProof);
  });
});
