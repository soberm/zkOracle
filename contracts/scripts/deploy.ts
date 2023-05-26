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

  const AggregationVerifier = await ethers.getContractFactory(
    "AggregationVerifier"
  );

  const aggregationVerifier = await AggregationVerifier.deploy();

  const SlashingVerifier = await ethers.getContractFactory("SlashingVerifier");

  const slashingVerifier = await SlashingVerifier.deploy();

  const ZKOracle = await ethers.getContractFactory(
    "contracts/ZKOracle.sol:ZKOracle"
  );

  const zkOracle = await ZKOracle.deploy(
    merkleTree.address,
    aggregationVerifier.address,
    slashingVerifier.address,
    "5491184307399689246197683245202605692069525215510636283504164930708453453685",
    "2576048849028791939551994783150968389338965397796293068226051430557680319904"
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
