const { buildModule } = require("@nomicfoundation/hardhat-ignition/modules");

module.exports = buildModule("Layer8WalletModule", (m) => {
    const layer8Wallet = m.contract("Layer8Wallet", []);

    return { layer8Wallet };
});
