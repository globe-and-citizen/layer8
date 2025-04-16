import { polygon } from 'https://esm.sh/@wagmi/core@2.11.6/chains'
import { createWeb3Modal, defaultWagmiConfig } from 'https://esm.sh/@web3modal/wagmi?bundle'

const WALLET_CONNECT_PROJECT_ID = "339b2396ae9c42ea2ffd29ab386c8a9e"

const WALLET_CONNECT_METADATA = {
    name: 'Layer8',
    description: 'Reverse Proxy Implementing the Layer8 protocol',
    url: 'https://layer8proxy.net',
    icons: ['https://avatars.githubusercontent.com/u/37784886']
}

const CHAINS = [ polygon ]

export const WAGMI_CONFIG = defaultWagmiConfig({
    chains: CHAINS,
    projectId: WALLET_CONNECT_PROJECT_ID,
    metadata: WALLET_CONNECT_METADATA,
})

export const setupWeb3Modal = () => {
    const walletConfig = {
        wagmiConfig: WAGMI_CONFIG,
        projectId: WALLET_CONNECT_PROJECT_ID,
        enableOnramp: true,
    }
    createWeb3Modal(walletConfig)
}
