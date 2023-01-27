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

  it("verify", async function () {
    const { merkleTree } = await loadFixture(deploy);
    for (let i = 1; i < 9; i++) {
      await merkleTree.insert(await merkleTree.hash([i]));
    }

    let result = await merkleTree.verify(
      [
        "2",
        "2317231519361726365316928999413088996422906353456011146656543284985006859395",
        "6757385243957864171620915041150172487332935094868613572887924196788898827793",
        "18480134192469085030266242755750140883208719495719185708340400278682450074026",
      ],
      [0, 1, 1]
    );

    expect(result).to.be.true;
  });
});
