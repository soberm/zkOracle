import { ethers } from "hardhat";

async function main() {
  const MiMC = await ethers.getContractFactory("MiMC");
  const miMC = await MiMC.deploy();

  const MerkleTree = await ethers.getContractFactory("MerkleTree", {
    libraries: {
      MiMC: miMC.address,
    },
  });

  const merkleTree = await MerkleTree.deploy(2);

  const Verifier = await ethers.getContractFactory("Verifier");

  const verifier = await Verifier.deploy();

  const ZKOracle = await ethers.getContractFactory(
    "contracts/ZKOracle.sol:ZKOracle"
  );

  const zkOracle = await ZKOracle.deploy(
    merkleTree.address,
    verifier.address,
    0,
    1
  );

  await zkOracle.deployed();

  console.log(`zkOracle deployed to ${zkOracle.address}`);
}

// We recommend this pattern to be able to use async/await everywhere
// and properly handle errors.
main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
