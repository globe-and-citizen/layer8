import { HardhatUserConfig } from "hardhat/config";
import "@nomicfoundation/hardhat-toolbox";

const config: HardhatUserConfig = {
  solidity: "0.8.24",
  defaultNetwork: "sepolia",
  networks: {
    sepolia: {
      url: "https://eth-sepolia.g.alchemy.com/v2/pDkGIXxTZkhVJIrW01sJe4tVOO2oHtcj",
      accounts: [
        "0x20844dd9f85f374a3102609f7dacd53d75b26e41f937188b7bbdb9b524707cae"
      ]
    }
  }
};

export default config;