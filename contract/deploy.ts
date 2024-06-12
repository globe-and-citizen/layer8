import { ethers } from "hardhat";

async function main() {
    const Langle = await ethers.getContractFactory("Langle");
    const langle = await Langle.deploy();
    console.log("Contract Deployed to Address:", await langle.getAddress());
  }

  main()
    .then(() => process.exit(0))
    .catch(error => {
      console.error(error);
      process.exit(1);
    });