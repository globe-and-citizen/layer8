const { buildModule } = require("@nomicfoundation/hardhat-ignition/modules");
const { vars } = require("hardhat/config");

const LAYER8_WALLET_ADDRESS = vars.get("LAYER8_WALLET_ADDRESS");

module.exports = buildModule("PayWithCryptoModule", (m) => {
    const payWithCrypto = m.contract("PayWithCrypto", [LAYER8_WALLET_ADDRESS]);

    return { payWithCrypto };
});