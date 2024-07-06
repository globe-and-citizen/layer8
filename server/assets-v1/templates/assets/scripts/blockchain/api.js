import { abi } from "./abi.js";
import { readContract, writeContract } from 'https://esm.sh/@wagmi/core@2.11.6'
import { WAGMI_CONFIG } from './web3modal.js'

const SMART_CONTRACT_ADDRESS = "0xEf0129b8493d98596B2E00C964DB616B24CccdA5"

export const getClientById = async (id) => {
    const clientData = await readContract(WAGMI_CONFIG, {
        abi,
        address: SMART_CONTRACT_ADDRESS,
        functionName: 'getClientById',
        args: [id],
    });

    return {
        id: clientData.clientId,
        unpaidBill: clientData.unpaidBill
    }
}

export const payBill = async (id, amount) => {
    await writeContract(WAGMI_CONFIG, {
        abi,
        address: SMART_CONTRACT_ADDRESS,
        functionName: 'payBill',
        args: [id],
        value: amount,
    })
}