import { writeContract } from 'https://esm.sh/@wagmi/core@2.11.6';
import { WAGMI_CONFIG } from "./web3modal.js";
import { abi } from "./abi.js";
import { polygon } from 'https://esm.sh/@wagmi/core@2.11.6/chains'

const payWithCryptoContractAddress = "0x7Ff8Cd44330964c95f5A8C442E75cfa50fc8D450";

export const payBill = async (amount)=> {
    try {
        const result = await writeContract(WAGMI_CONFIG, {
            abi,
            address: payWithCryptoContractAddress,
            functionName: "pay",
            value: amount,
            chainId: polygon.id,
        });

        console.log(result);

        return result;
    } catch (error) {
        return error;
    }
}
