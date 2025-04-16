require("@nomicfoundation/hardhat-toolbox");

const { vars } = require("hardhat/config");

const ALCHEMY_API_KEY = vars.get("ALCHEMY_API_KEY");

const SEPOLIA_PRIVATE_KEY = vars.get("SEPOLIA_PRIVATE_KEY");

const POLYGON_PRIVATE_KEY = vars.get("POLYGON_PRIVATE_KEY");

const POLYGONSCAN_API_KEY = vars.get("POLYGONSCAN_API_KEY");

module.exports = {
  solidity: "0.8.28",
  networks: {
    sepolia: {
      url: `https://eth-sepolia.g.alchemy.com/v2/${ALCHEMY_API_KEY}`,
      accounts: [SEPOLIA_PRIVATE_KEY]
    },
    matic: {
      url: `https://polygon-mainnet.g.alchemy.com/v2/${ALCHEMY_API_KEY}`,
      accounts: [POLYGON_PRIVATE_KEY]
    }
  },
  etherscan: {
    apiKey: POLYGONSCAN_API_KEY,
  }
};
