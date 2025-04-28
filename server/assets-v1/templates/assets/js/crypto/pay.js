import {waitForTransactionReceipt, writeContract} from 'https://esm.sh/@wagmi/core@2.11.6';
import {abi} from "./abi.js";
import {polygon} from 'https://esm.sh/@wagmi/core@2.11.6/chains'

export const payBill = async (WAGMI_CONFIG, payWithCryptoContractAddress, clientId, amount)=> {
    try {
        return await writeContract(WAGMI_CONFIG, {
            abi,
            address: payWithCryptoContractAddress,
            functionName: "pay",
            args: [clientId],
            value: amount,
            chainId: polygon.id,
        });
    } catch (error) {
        return error;
    }
}

export const awaitTransactionConfirmation = async (WAGMI_CONFIG, transactionHash) => {
    const receipt = await waitForTransactionReceipt(WAGMI_CONFIG, {
        hash: transactionHash,
        confirmations: 1,
        chainId: polygon.id,
    });

    return receipt.status;
}
