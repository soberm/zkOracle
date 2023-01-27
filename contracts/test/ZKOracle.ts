import { ethers } from "hardhat";
import { loadFixture } from "@nomicfoundation/hardhat-network-helpers";
import { expect } from "chai";
describe("ZKOracle", function () {
  async function deploy() {
    const MiMC = await ethers.getContractFactory("MiMC");
    const miMC = await MiMC.deploy();

    const MerkleTree = await ethers.getContractFactory("MerkleTree", {
      libraries: {
        MiMC: miMC.address,
      },
    });

    const merkleTree = await MerkleTree.deploy(3);

    const ZKOracle = await ethers.getContractFactory("ZKOracle");

    const zkOracle = await ZKOracle.deploy(merkleTree.address);

    return { zkOracle };
  }

  it("register", async function () {
    const { zkOracle } = await loadFixture(deploy);
    const tx = await zkOracle.register(1, 2);
    const receipt = await tx.wait();
    console.log(receipt.cumulativeGasUsed);
  });
});
