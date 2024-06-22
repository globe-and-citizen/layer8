import { ethers } from "hardhat";
import { expect } from "chai";
import { PayAsYouGo } from "../typechain-types";
import { SignerWithAddress } from "@nomicfoundation/hardhat-ethers/signers";

describe("PayAsYouGo Contract", () => {
  let contract: PayAsYouGo;
  let owner: SignerWithAddress;
  let user1: SignerWithAddress;

  beforeEach(async () => {
    [owner, user1] = await ethers.getSigners();

    const PayAsYouGoFactory = await ethers.getContractFactory("PayAsYouGo", owner);
    contract = await PayAsYouGoFactory.deploy(owner.address) as PayAsYouGo;
    await contract.waitForDeployment();
  });


  it("Should create a new contract and emit an event", async () => {
    const clientId = "client123";
    const rate = 10;
    await contract.newContract(rate, clientId)

    const contractId = await contract.contractIds(0);
    const retrievedContract = await contract.getContractById(contractId);
    expect(retrievedContract.clientId).to.equal(clientId);
    expect(retrievedContract.rate).to.equal(rate);
  });

  it("Should allow owner to add bills to a contract", async () => {
    await contract.newContract(10, "client123");

    const contractId = await contract.contractIds(0);
    await contract.addBillToContract(contractId, 5, 1686592800); 

    const retrievedContract = await contract.getContractById(contractId);
    expect(retrievedContract.unpaidBill).to.equal(50); 
  });

  it("Should allow users to pay bills and transfer funds", async () => {
    await contract.newContract(10, "client123");
    const contractId = await contract.contractIds(0);
    await contract.addBillToContract(contractId, 5, 1686592800); 

    const initialOwnerBalance = await owner.provider.getBalance(owner.address);

    await expect(contract.connect(user1).payBill(contractId, { value: 50 }))
      .to.changeEtherBalances([user1, owner], [-50, 50]);

    const finalOwnerBalance = await owner.provider.getBalance(owner.address);
    expect(finalOwnerBalance).to.be.gt(initialOwnerBalance); 

    const retrievedContract = await contract.getContractById(contractId);
    expect(retrievedContract.unpaidBill).to.equal(0); 
  });

});