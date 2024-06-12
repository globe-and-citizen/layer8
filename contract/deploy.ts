import { ethers } from "hardhat";

async function main() {
    const PayAsYouGo = await ethers.getContractFactory("PayAsYouGo");
    const payAsYouGo = await PayAsYouGo.deploy("0x0dc02b96b0960400bcF02A0Fcb87aB0BE80A2264");
    console.log("Contract Deployed to Address:", await payAsYouGo.getAddress());
  }

  main()
    .then(() => process.exit(0))
    .catch(error => {
      console.error(error);
      process.exit(1);
    });