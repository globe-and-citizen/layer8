import { sepolia, mainnet } from 'https://esm.sh/@wagmi/core@2.11.6/chains'
import { createWeb3Modal, defaultWagmiConfig } from 'https://esm.sh/@web3modal/wagmi?bundle'

const WALLET_CONNECT_PROJECT_ID = "a2b392649ab49d473ca531f98ed09ae0"

const WALLET_CONNECT_METADATA = {
    name: 'Layer8',
    description: 'Decentralized Reverse Proxy Platform',
    url: 'https://layer8proxy.net',
    icons: ['https://avatars.githubusercontent.com/u/37784886']
}

const SELECTED_CHAINS = [sepolia, mainnet]

export const WAGMI_CONFIG = defaultWagmiConfig({
    chains: SELECTED_CHAINS,
    projectId: WALLET_CONNECT_PROJECT_ID,
    metadata: WALLET_CONNECT_METADATA,
})

const web3ModalConfig = {
    wagmiConfig: WAGMI_CONFIG,
    projectId: WALLET_CONNECT_PROJECT_ID,
    enableOnramp: true // Optional - false as default
}

export const setupWeb3Modal = () => {
    createWeb3Modal(web3ModalConfig)
}