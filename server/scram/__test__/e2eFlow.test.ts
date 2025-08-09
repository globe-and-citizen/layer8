import fetchMock from 'jest-fetch-mock';
import { jest, describe, test, expect, beforeEach } from '@jest/globals';
import {
    keysHMAC,
    signatureHMAC,
    hexStringToBytes,
    bytesToHexString,
    xorBytes,
} from "../cryptoUtils";

// Mock the global fetch
global.fetch = fetchMock as any;

// Mock the mnemonic module
const mnemonic = {
    generateBip39Mnemonic: jest.fn(() => 'mock-mnemonic'),
    getPrivateAndPublicKeys: jest.fn(() => ({
        publicKey: new Uint8Array([1, 2, 3]),
    })),
};

// Global variable to store keys
let registerUserKeys: { storedKey: string; serverKey: string; clientKey: string } | null = null;
let loginUserKeys: { storedKey: string; serverKey: string; clientKey: string } | null = null;
let resetPasswordKeys: { storedKey: string; serverKey: string; clientKey: string } | null = null;

// Define the registerUser function and its dependencies
const registerUser = async () => {
    const registerUsername = 'testuser';
    const registerPassword = 'testpass';
    const registerFirstName = 'Test';
    const registerLastName = 'User';
    const registerDisplayName = 'Test User';
    const registerCountry = 'Canada';
    const showToastMessage = jest.fn();
    const modalWindowActive = { value: false };

    try {
        const keyPair = mnemonic.getPrivateAndPublicKeys();

        const responseOne = await fetch('[[ .ProxyURL ]]/api/v1/register-user-precheck', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                username: registerUsername,
            }),
        });

        const registerPrecheckResponseBody = await responseOne.json();
        if (responseOne.status !== 201) {
            showToastMessage('Something went wrong!', 'error');
            return;
        }

        const { data } = keysHMAC(
            registerPassword,
            registerPrecheckResponseBody.data.salt,
            registerPrecheckResponseBody.data.iterationCount
        );

        // Save the keys to the global variable
        registerUserKeys = {
            storedKey: data.storedKey,
            serverKey: data.serverKey,
            clientKey: data.clientKey
        };

        const resp = await fetch('[[ .ProxyURL ]]/api/v1/register-user', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                username: registerUsername,
                first_name: registerFirstName,
                last_name: registerLastName,
                display_name: registerDisplayName,
                country: registerCountry,
                public_key: Array.from(keyPair.publicKey),
                stored_key: data.storedKey,
                server_key: data.serverKey,
            }),
        });

        const registerResponseBody = await resp.json();
        if (resp.status === 201) {
            modalWindowActive.value = true;
        } else if (registerResponseBody.message) {
            showToastMessage(registerResponseBody.message, 'error');
        } else {
            showToastMessage('Something went wrong!', 'error');
        }
    } catch (error) {
        console.error(error);
        showToastMessage('Registration failed!', 'error');
    }
};

// Define the loginUser function and its dependencies
const loginUser = async (loginUsername: string, loginPassword: string, cNonce: string) => {
    // const loginUsername = 'testuser';
    // const loginPassword = 'testpass';
    // const cNonce = "42ZHCYjtfdaykloIJFvxOLlYZ0MGMLtqwKKDnfUjfC0=";
    const showToastMessage = jest.fn();

    try {
        const responseOne = await fetch('[[ .ProxyURL ]]/api/v1/login-user-precheck', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                username: loginUsername,
                c_nonce: cNonce
            }),
        });

        const loginPrecheckResponseBody = await responseOne.json();
        if (responseOne.status !== 200) {
            showToastMessage('Something went wrong!', 'error');
            return;
        }

        const { data } = keysHMAC(
            loginPassword,
            loginPrecheckResponseBody.data.salt,
            loginPrecheckResponseBody.data.iter_count
        );

        // Save the keys to the global variable
        loginUserKeys = {
            storedKey: data.storedKey,
            serverKey: data.serverKey,
            clientKey: data.clientKey
        };

        const clientKeyBytes = hexStringToBytes(data.clientKey);

        const authMessage = `[n=${loginUsername},r=${cNonce},s=${loginPrecheckResponseBody.data.salt},i=${loginPrecheckResponseBody.data.iter_count},r=${loginPrecheckResponseBody.data.nonce}]`;

        // Signature HMAC
        const clientSignature = signatureHMAC(authMessage, data.storedKey);

        const clientSignatureBytes = hexStringToBytes(clientSignature);

        const clientProofBytes = xorBytes(
            clientKeyBytes,
            clientSignatureBytes
        );

        const clientProof = bytesToHexString(clientProofBytes);

        const loginUserResponse = await fetch(
            "[[ .ProxyURL ]]/api/v1/login-user",
            {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({
                    username: loginUsername,
                    nonce: loginPrecheckResponseBody.data.nonce,
                    c_nonce: cNonce,
                    client_proof: clientProof,
                }),
            }
        );

        const loginUserResponseJSON = await loginUserResponse.json();

        if (loginUserResponseJSON.data.server_signature) {
            const serverSignatureCheck = signatureHMAC(authMessage, data.serverKey);

            if (serverSignatureCheck === loginUserResponseJSON.data.server_signature) {
                showToastMessage("Login successful!", "success");
            }
        } else if (loginUserResponseJSON.message) {
            showToastMessage(loginUserResponseJSON.message, "error");
        } else {
            showToastMessage("Login failed, please try again later", "error");
        }
    } catch (error) {
        console.error(error);
        showToastMessage('Registration failed!', 'error');
    }
};

const resetPassword = async () => {
    const loginUsername = 'testuser';
    const newPassword = 'newpass';
    const signature = "test signature";
    const alert = jest.fn();

    try {
        const responseOne = await fetch(
            "[[ .ProxyURL ]]/api/v1/reset-password-precheck",
            {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({
                    username: loginUsername,
                }),
            },
        );

        const resetPasswordPrecheckResponseBody = await responseOne.json();
        if (responseOne.status !== 200) {
            alert("Error: " + resetPasswordPrecheckResponseBody.message);
            return;
        }

        const { data } = keysHMAC(
            newPassword,
            resetPasswordPrecheckResponseBody.data.salt,
            resetPasswordPrecheckResponseBody.data.iterationCount
        );


        // Save the keys to the global variable
        resetPasswordKeys = {
            storedKey: data.storedKey,
            serverKey: data.serverKey,
            clientKey: data.clientKey
        }

        const responseTwo = await fetch("[[ .ProxyURL ]]/api/v1/reset-password", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                username: loginUsername,
                signature: Array.from(signature),
                stored_key: data.storedKey,
                server_key: data.serverKey,
            }),
        });

        const resetPasswordResponseBody = await responseTwo.json();

        if (resetPasswordResponseBody.is_success === true) {
            alert(resetPasswordResponseBody.message);
        } else {
            alert("Error: " + resetPasswordResponseBody.message);
        }
    } catch (error) {
        console.error(error);
        alert("Error happened");
    }
}

describe('Perform a Register and Login flow and match their keys', () => {
    test('call the registerUser function, extract keys, and save them globally', async () => {
        const showToastMessage = jest.fn();
        const modalWindowActive = { value: false };
        const mockResponsePrecheck = { data: { salt: 'mock-salt', iterationCount: 4096 } };
        const mockResponseRegister = { message: 'User registered successfully' };

        fetchMock.mockResponses(
            [JSON.stringify(mockResponsePrecheck), { status: 201 }],
            [JSON.stringify(mockResponseRegister), { status: 201 }]
        );

        await registerUser.call({
            showToastMessage,
            modalWindowActive,
            registerUsername: 'testuser',
            registerPassword: 'testpass',
            registerFirstName: 'Test',
            registerLastName: 'User',
            registerDisplayName: 'Test User',
            registerCountry: 'USA',
        });

        expect(fetchMock).toHaveBeenCalledTimes(2);
        expect(showToastMessage).not.toHaveBeenCalled();
        // Verify that the keys were saved globally
        expect(registerUserKeys).not.toBeNull();
        expect(registerUserKeys?.storedKey).toBeDefined();
        expect(registerUserKeys?.serverKey).toBeDefined();
        expect(registerUserKeys?.clientKey).toBeDefined();
    });

    test('call the loginUser function, and match their keys', async () => {
        const showToastMessage = jest.fn();
        const modalWindowActive = { value: false };
        const mockResponsePrecheck = { data: { salt: 'mock-salt', iter_count: 4096 } };
        const mockResponseLogin = { data: { server_signature: 'mock-server-signature', token: 'mock-token' } };

        fetchMock.mockResponses(
            [JSON.stringify(mockResponsePrecheck), { status: 200 }],
            [JSON.stringify(mockResponseLogin), { status: 200 }]
        );

        await loginUser('testuser', 'testpass', "42ZHCYjtfdaykloIJFvxOLlYZ0MGMLtqwKKDnfUjfC0=");

        expect(fetchMock).toHaveBeenCalledTimes(4);
        expect(showToastMessage).not.toHaveBeenCalled();
        // Verify that the keys were saved globally
        expect(loginUserKeys).not.toBeNull();
        expect(loginUserKeys?.storedKey).toBeDefined();
        expect(loginUserKeys?.serverKey).toBeDefined();
        expect(loginUserKeys?.clientKey).toBeDefined();
        // Check if the keys match
        expect(loginUserKeys?.storedKey).toEqual(registerUserKeys?.storedKey);
        expect(loginUserKeys?.serverKey).toEqual(registerUserKeys?.serverKey);
        expect(loginUserKeys?.clientKey).toEqual(registerUserKeys?.clientKey);
    });

    test('call the resetPassword function, and match their keys', async () => {
        const alert = jest.fn();
        const modalWindowActive = { value: false };
        const mockResponsePrecheck = { data: { salt: 'mock-salt', iterationCount: 4096 } };
        const mockResponseResetPasword = { is_success: true, message: 'Password reset successfully' };

        fetchMock.mockResponses(
            [JSON.stringify(mockResponsePrecheck), { status: 200 }],
            [JSON.stringify(mockResponseResetPasword), { status: 200 }]
        );

        await resetPassword.call({
            alert,
            modalWindowActive,
            loginUsername: 'testuser',
            mnemonicSentence: 'test mnemonic',
            newPassword: 'newpass',
            repeatedNewPassword: 'newpass'
        });

        expect(fetchMock).toHaveBeenCalledTimes(6);
        expect(alert).not.toHaveBeenCalled();
        // Verify that the keys were saved globally
        expect(resetPasswordKeys).not.toBeNull();
        expect(resetPasswordKeys?.storedKey).toBeDefined();
        expect(resetPasswordKeys?.serverKey).toBeDefined();
        expect(resetPasswordKeys?.clientKey).toBeDefined();
        // Make sure that key no longer matches
        expect(resetPasswordKeys?.storedKey).not.toEqual(registerUserKeys?.storedKey);
        expect(resetPasswordKeys?.serverKey).not.toEqual(registerUserKeys?.serverKey);
        expect(resetPasswordKeys?.clientKey).not.toEqual(registerUserKeys?.clientKey);
    });

    test('call the loginUser function again to match their keys with the new ones', async () => {
        const showToastMessage = jest.fn();
        const modalWindowActive = { value: false };
        const mockResponsePrecheck = { data: { salt: 'mock-salt', iter_count: 4096 } };
        const mockResponseLogin = { data: { server_signature: 'mock-server-signature', token: 'mock-token' } };

        fetchMock.mockResponses(
            [JSON.stringify(mockResponsePrecheck), { status: 200 }],
            [JSON.stringify(mockResponseLogin), { status: 200 }]
        );

        await loginUser('testuser', 'newpass', "42ZHCYjtfdaykloIJFvxOLlYZ0MGMLtqwKKDnfUjfC0=");

        expect(fetchMock).toHaveBeenCalledTimes(8);
        expect(showToastMessage).not.toHaveBeenCalled();
        // Verify that the keys were saved globally
        expect(loginUserKeys).not.toBeNull();
        expect(loginUserKeys?.storedKey).toBeDefined();
        expect(loginUserKeys?.serverKey).toBeDefined();
        expect(loginUserKeys?.clientKey).toBeDefined();
        // Check if the keys match
        expect(loginUserKeys?.storedKey).toEqual(resetPasswordKeys?.storedKey);
        expect(loginUserKeys?.serverKey).toEqual(resetPasswordKeys?.serverKey);
        expect(loginUserKeys?.clientKey).toEqual(resetPasswordKeys?.clientKey);
    });
});