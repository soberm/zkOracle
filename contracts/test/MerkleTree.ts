import { ethers } from "hardhat";
import { loadFixture } from "@nomicfoundation/hardhat-network-helpers";
import { expect } from "chai";

describe("MerkleTree", function () {
  async function deploy() {
    const MiMC = await ethers.getContractFactory("MiMC");
    const miMC = await MiMC.deploy();

    const MerkleTree = await ethers.getContractFactory("MerkleTree", {
      libraries: {
        MiMC: miMC.address,
      },
    });

    const merkleTree = await MerkleTree.deploy(3);
    return { merkleTree };
  }

  it("insert", async function () {
    const { merkleTree } = await loadFixture(deploy);
    for (let i = 1; i < 9; i++) {
      let tx = await merkleTree.insert(await merkleTree.hash([i]));
      let receipt = await tx.wait();
    }

    console.log(await merkleTree.getRoot());
  });
});
