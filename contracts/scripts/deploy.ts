import { ethers } from "hardhat";

async function main() {
  const zkOracleContractFactory = await ethers.getContractFactory("ZKOracle");
  const zkOracle = await zkOracleContractFactory.deploy();

  await zkOracle.deployed();

  console.log(`zkOracle deployed to ${zkOracle.address}`);
}

// We recommend this pattern to be able to use async/await everywhere
// and properly handle errors.
main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
