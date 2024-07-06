import { abi } from "./abi.js";
import { readContract, writeContract } from 'https://esm.sh/@wagmi/core@2.11.6'
import { WAGMI_CONFIG } from './web3modal.js'

export const getClientById = async (smartContractAddress, id) => {
    const clientData = await readContract(WAGMI_CONFIG, {
        abi,
        address: smartContractAddress,
        functionName: 'getClientById',
        args: [id],
    });

    return {
        id: clientData.clientId,
        unpaidBill: clientData.unpaidBill
    }
}

export const payBill = async (smartContractAddress, id, amount) => {
    await writeContract(WAGMI_CONFIG, {
        abi,
        address: smartContractAddress,
        functionName: 'payBill',
        args: [id],
        value: amount,
    })
}