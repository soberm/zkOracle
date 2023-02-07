import { ethers } from "hardhat";
import { loadFixture } from "@nomicfoundation/hardhat-network-helpers";
import { expect } from "chai";

describe("MerkleTree", function () {
  async function deployFixture() {
    const MiMC = await ethers.getContractFactory("MiMC");
    const miMC = await MiMC.deploy();

    const MerkleTree = await ethers.getContractFactory("MerkleTree", {
      libraries: {
        MiMC: miMC.address,
      },
    });

    const merkleTree = await MerkleTree.deploy(2);
    return { merkleTree };
  }

  async function fullTreeFixture() {
    const { merkleTree } = await deployFixture();
    await merkleTree.insert(
      "5653921414509890009709956482085915841678544141157758128975480650296143636189"
    );
    await merkleTree.insert(
      "6187350521517272486117148237635192271041670665937219625917563897233388432910"
    );
    await merkleTree.insert(
      "11788837250340196669906759998716062330026986825826762343780050738774786628008"
    );
    await merkleTree.insert(
      "5930779311236679553824249069471014709811380042820293899590311964685757495005"
    );
    return { merkleTree };
  }

  it("insert should compute correct root", async function () {
    const { merkleTree } = await loadFixture(fullTreeFixture);
    expect(await merkleTree.getRoot()).to.equal(
      "14669474610586036795355717835678114088252099439454524063142345888814163114094"
    );
  });

  it("update should update the root with given path", async function () {
    const { merkleTree } = await loadFixture(fullTreeFixture);

    await merkleTree.update(
      "5930779311236679553824249069471014709811380042820293899590311964685757495005",
      [
        "6205836767976675972003616131281122597329838876688673966264947687549471156130",
        "6187350521517272486117148237635192271041670665937219625917563897233388432910",
        "6584420187527354519485242243152059973161633040054605148949951423401777995392",
      ],
      [1, 1]
    );

    let result = await merkleTree.verify(
      [
        "5930779311236679553824249069471014709811380042820293899590311964685757495005",
        "6187350521517272486117148237635192271041670665937219625917563897233388432910",
        "6584420187527354519485242243152059973161633040054605148949951423401777995392",
      ],
      [1, 1]
    );

    expect(result).to.be.true;
  });

  it("update should revert if tree not full", async function () {
    const { merkleTree } = await loadFixture(deployFixture);

    let tx = merkleTree.update(
      "5930779311236679553824249069471014709811380042820293899590311964685757495005",
      [
        "6205836767976675972003616131281122597329838876688673966264947687549471156130",
        "6187350521517272486117148237635192271041670665937219625917563897233388432910",
        "6584420187527354519485242243152059973161633040054605148949951423401777995392",
      ],
      [1, 1]
    );

    await expect(tx).to.be.revertedWith("tree not full");
  });

  it("update should revert if proof is invalid", async function () {
    const { merkleTree } = await loadFixture(fullTreeFixture);

    let tx = merkleTree.update(
      "5930779311236679553824249069471014709811380042820293899590311964685757495005",
      [
        "6205836767976675972003616131281122597329838876688673966264947687549471156130",
        "6287350521517272486117148237635192271041670665937219625917563897233388432910",
        "6584420187527354519485242243152059973161633040054605148949951423401777995392",
      ],
      [1, 1]
    );

    expect(tx).to.be.revertedWith("leaf to update not included");
  });

  it("verify should return true if data is included", async function () {
    const { merkleTree } = await loadFixture(fullTreeFixture);

    let result = await merkleTree.verify(
      [
        "6205836767976675972003616131281122597329838876688673966264947687549471156130",
        "6187350521517272486117148237635192271041670665937219625917563897233388432910",
        "6584420187527354519485242243152059973161633040054605148949951423401777995392",
      ],
      [1, 1]
    );

    expect(result).to.be.true;
  });

  it("computeRootFromPath should compute root", async function () {
    const { merkleTree } = await loadFixture(deployFixture);

    let result = await merkleTree.computeRootFromPath(
      [
        "6205836767976675972003616131281122597329838876688673966264947687549471156130",
        "6187350521517272486117148237635192271041670665937219625917563897233388432910",
        "6584420187527354519485242243152059973161633040054605148949951423401777995392",
      ],
      [1, 1]
    );

    expect(result).to.equal(
      "14669474610586036795355717835678114088252099439454524063142345888814163114094"
    );
  });
});
