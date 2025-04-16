const {
    time,
    loadFixture,
  } = require("@nomicfoundation/hardhat-toolbox/network-helpers");
  const { expect } = require("chai");
const { ethers } = require("hardhat");
  
describe("PayWithCrypto", function() {
    const clientId = "clientID";
    const paymentAmount = 50000;

    async function deployPayWithCryptoFixture() {
        const [owner, layer8Account] = await ethers.getSigners();

        const PayWithCrypto = await ethers.getContractFactory("PayWithCrypto");
        const payWithCrypto = await PayWithCrypto.deploy(layer8Account.address);

        return {payWithCrypto, owner, layer8Account};
    }

    describe("Getters and Setters", function() {
        it("should receive the right wallet address", async function() {
            const { payWithCrypto, layer8Account } = await loadFixture(deployPayWithCryptoFixture);

            expect(await payWithCrypto.layer8WalletAddress()).to.equal(layer8Account.address);
        });

        it("assert layer8 wallet address is set correctly", async function () {
            const [owner, _, account2] = await ethers.getSigners();
            const { payWithCrypto } = await loadFixture(deployPayWithCryptoFixture);

            await expect(payWithCrypto.setLayer8WalletAddress(account2.address)).not.to.be.reverted;

            expect(await payWithCrypto.layer8WalletAddress()).to.equal(account2.address);
        });

        it("assert that setLayer8WalletAddress is reverted when called by a non-owner", async function() {
            const [_, account1, account2] = await ethers.getSigners();

            const { payWithCrypto, } = await loadFixture(deployPayWithCryptoFixture);

            await expect(payWithCrypto.connect(account2).setLayer8WalletAddress(account1.address)).to.be.reverted;
        });
    });

    describe("Pay functionality works correctly", function() {
        it("funds are transferred to the layer8 wallet address successfully", async function() {
            const [owner, _, account2] = await ethers.getSigners();

            const { payWithCrypto, layer8Account } = await loadFixture(deployPayWithCryptoFixture);

            const result = payWithCrypto.connect(account2).pay(clientId, { value: paymentAmount });

            await expect(result).not.to.be.reverted;
            await expect(result).to.emit(payWithCrypto, "TrafficPaid").withArgs(clientId, account2.address, paymentAmount);
            await expect(result).to.changeEtherBalance(layer8Account, paymentAmount);
        });
    });
});
