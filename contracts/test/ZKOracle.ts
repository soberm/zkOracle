import { ethers } from "hardhat";
import { loadFixture, time } from "@nomicfoundation/hardhat-network-helpers";
import { expect } from "chai";

describe("ZKOracle", function () {
  async function setup() {
    const [owner, addr1, addr2] = await ethers.getSigners();

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

    return { zkOracle, owner, addr1, addr2 };
  }

  it("register", async function () {
    const { zkOracle } = await loadFixture(setup);
    const tx = await zkOracle.register(1, 2);
    const receipt = await tx.wait();
    console.log(receipt.cumulativeGasUsed);
  });

  it("exit should emit event", async function () {
    const { zkOracle, owner } = await loadFixture(setup);
    await expect(zkOracle.exit())
      .to.emit(zkOracle, "Exiting")
      .withArgs(owner.address);
  });

  it("exit should set time", async function () {
    const { zkOracle, owner } = await loadFixture(setup);
    const tx = await zkOracle.exit();
    const receipt = await tx.wait();

    const block = await ethers.provider.getBlock(receipt.blockNumber);
    const timestamp = ethers.BigNumber.from(block.timestamp);
    const exitDelay = await zkOracle.exitDelay();

    const expected = timestamp.add(exitDelay);
    const actual = await zkOracle.getExitTime(owner.address);

    expect(actual).to.equal(expected);
  });

  it("withdraw should revert when time not passed", async function () {
    const { zkOracle } = await loadFixture(setup);
    expect(zkOracle.withdraw(null, [], [])).to.be.revertedWith(
      "time not passed"
    );
  });

  it("withdraw should revert when index wrong", async function () {
    const { zkOracle } = await loadFixture(setup);
    const account = { index: 1 };
    expect(zkOracle.withdraw(account, [], [])).to.be.revertedWith(
      "wrong index"
    );
  });
});
