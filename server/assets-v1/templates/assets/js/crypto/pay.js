import { writeContract } from 'https://esm.sh/@wagmi/core@2.11.6';
import { WAGMI_CONFIG } from "./web3modal.js";
import { abi } from "./abi.js";
import { polygon } from 'https://esm.sh/@wagmi/core@2.11.6/chains'

export const payBill = async (payWithCryptoContractAddress, clientId, amount)=> {
    try {
        console.log(payWithCryptoContractAddress)

        const result = await writeContract(WAGMI_CONFIG, {
            abi,
            address: payWithCryptoContractAddress,
            functionName: "pay",
            args: [ clientId ],
            value: amount,
            chainId: polygon.id,
        });

        console.log(result);

        return result;
    } catch (error) {
        return error;
    }
}
