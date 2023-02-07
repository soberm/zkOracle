import { ethers } from "hardhat";
import {
  loadFixture,
  time,
  mine,
} from "@nomicfoundation/hardhat-network-helpers";
import { expect } from "chai";

describe("ZKOracle", function () {
  async function deployFixture() {
    const [owner, addr1, addr2] = await ethers.getSigners();

    const MiMC = await ethers.getContractFactory("MiMC");
    const miMC = await MiMC.deploy();

    const MerkleTree = await ethers.getContractFactory("MerkleTree", {
      libraries: {
        MiMC: miMC.address,
      },
    });

    const merkleTree = await MerkleTree.deploy(2);

    const ZKOracle = await ethers.getContractFactory("ZKOracle");

    const zkOracle = await ZKOracle.deploy(merkleTree.address);

    return { zkOracle, owner, addr1, addr2 };
  }

  async function fullTreeFixture() {
    const { zkOracle, owner, addr1, addr2 } = await deployFixture();
    await zkOracle.register(
      {
        x: "7794373982259243195870592346785104092432649697832080133780782253104282782817",
        y: "16580021058669382711579094818964719809751621462486576197985799831318116474539",
      },
      { value: ethers.utils.parseEther("0") }
    );
    await zkOracle.register(
      {
        x: "19259775561661129033490267958867540712323727687827132538598271216435741353390",
        y: "17569980102909044676256001479640411087334772294425842357847622915410385256152",
      },
      { value: ethers.utils.parseEther("0") }
    );
    await zkOracle.register(
      {
        x: "2601225367854716029338107863118094577720373831631316901128802328263635799774",
        y: "20145565086840628487646378555659304143430966406145099422862446733209566497019",
      },
      { value: ethers.utils.parseEther("0") }
    );
    await zkOracle.register(
      {
        x: "8937264771331091228241633643830445497309816371951727702041886381904154965568",
        y: "11149480284076626963282053785587761751146600894461318988418480014122904346043",
      },
      { value: ethers.utils.parseEther("0") }
    );
    return { zkOracle, owner, addr1, addr2 };
  }

  it("replace should revert if value too low", async function () {
    const { zkOracle, owner } = await loadFixture(fullTreeFixture);
    let tx = zkOracle.replace(
      {
        x: "12087715944556734152926489964688099768311811690770181328500068944869829899594",
        y: "9368839392207299142928271420427271432427205218213129795542264289384888483164",
      },
      {
        index: 0,
        pubKeyX:
          "7794373982259243195870592346785104092432649697832080133780782253104282782817",
        pubKeyY:
          "16580021058669382711579094818964719809751621462486576197985799831318116474539",
        balance: ethers.utils.parseEther("1.0"),
      },
      [
        "6205836767976675972003616131281122597329838876688673966264947687549471156130",
        "6187350521517272486117148237635192271041670665937219625917563897233388432910",
        "6584420187527354519485242243152059973161633040054605148949951423401777995392",
      ],
      [1, 1]
    );
    await expect(tx).to.be.revertedWith("value too low");
  });

  it("replace should revert if account not included", async function () {
    const { zkOracle, owner } = await loadFixture(fullTreeFixture);
    let tx = zkOracle.replace(
      {
        x: "12087715944556734152926489964688099768311811690770181328500068944869829899594",
        y: "9368839392207299142928271420427271432427205218213129795542264289384888483164",
      },
      {
        index: 0,
        pubKeyX:
          "7794373982259243195870592346785104092432649697832080133780782253104282782817",
        pubKeyY:
          "16580021058669382711579094818964719809751621462486576197985799831318116474539",
        balance: ethers.utils.parseEther("0"),
      },
      [
        "5205836767976675972003616131281122597329838876688673966264947687549471156130",
        "6187350521517272486117148237635192271041670665937219625917563897233388432910",
        "6584420187527354519485242243152059973161633040054605148949951423401777995392",
      ],
      [1, 1],
      { value: ethers.utils.parseEther("1.0") }
    );
    await expect(tx).to.be.revertedWith("account not included");
  });

  it("replace should revert if leaf does not match account", async function () {
    const { zkOracle, owner } = await loadFixture(fullTreeFixture);
    let tx = zkOracle.replace(
      {
        x: "12087715944556734152926489964688099768311811690770181328500068944869829899594",
        y: "9368839392207299142928271420427271432427205218213129795542264289384888483164",
      },
      {
        index: 1,
        pubKeyX:
          "7794373982259243195870592346785104092432649697832080133780782253104282782817",
        pubKeyY:
          "16580021058669382711579094818964719809751621462486576197985799831318116474539",
        balance: ethers.utils.parseEther("0"),
      },
      [
        "6205836767976675972003616131281122597329838876688673966264947687549471156130",
        "6187350521517272486117148237635192271041670665937219625917563897233388432910",
        "6584420187527354519485242243152059973161633040054605148949951423401777995392",
      ],
      [1, 1],
      { value: ethers.utils.parseEther("1.0") }
    );
    await expect(tx).to.be.revertedWith("leaf does not match account");
  });

  it("replace should replace the account", async function () {
    const { zkOracle, owner, addr1 } = await loadFixture(fullTreeFixture);
    let tx = zkOracle.connect(addr1).replace(
      {
        x: "12087715944556734152926489964688099768311811690770181328500068944869829899594",
        y: "9368839392207299142928271420427271432427205218213129795542264289384888483164",
      },
      {
        index: 0,
        pubKeyX:
          "7794373982259243195870592346785104092432649697832080133780782253104282782817",
        pubKeyY:
          "16580021058669382711579094818964719809751621462486576197985799831318116474539",
        balance: ethers.utils.parseEther("0"),
      },
      [
        "6205836767976675972003616131281122597329838876688673966264947687549471156130",
        "6187350521517272486117148237635192271041670665937219625917563897233388432910",
        "6584420187527354519485242243152059973161633040054605148949951423401777995392",
      ],
      [1, 1],
      { value: ethers.utils.parseEther("1.0") }
    );
    await expect(tx)
      .to.emit(zkOracle, "Replaced")
      .withArgs(addr1.address, owner.address);
  });

  it("exit should emit event", async function () {
    const { zkOracle, owner } = await loadFixture(fullTreeFixture);
    await expect(
      zkOracle.exit(
        {
          index: 0,
          pubKeyX:
            "7794373982259243195870592346785104092432649697832080133780782253104282782817",
          pubKeyY:
            "16580021058669382711579094818964719809751621462486576197985799831318116474539",
          balance: ethers.utils.parseEther("0"),
        },
        [
          "6205836767976675972003616131281122597329838876688673966264947687549471156130",
          "6187350521517272486117148237635192271041670665937219625917563897233388432910",
          "6584420187527354519485242243152059973161633040054605148949951423401777995392",
        ],
        [1, 1]
      )
    )
      .to.emit(zkOracle, "Exiting")
      .withArgs(owner.address);
  });

  it("exit should revert with wrong sender address", async function () {
    const { zkOracle, owner, addr1 } = await loadFixture(fullTreeFixture);
    let tx = zkOracle.connect(addr1).exit(
      {
        index: 0,
        pubKeyX:
          "7794373982259243195870592346785104092432649697832080133780782253104282782817",
        pubKeyY:
          "16580021058669382711579094818964719809751621462486576197985799831318116474539",
        balance: ethers.utils.parseEther("0"),
      },
      [
        "6205836767976675972003616131281122597329838876688673966264947687549471156130",
        "6187350521517272486117148237635192271041670665937219625917563897233388432910",
        "6584420187527354519485242243152059973161633040054605148949951423401777995392",
      ],
      [1, 1]
    );

    await expect(tx).to.be.revertedWith("wrong sender address");
  });

  it("exit should revert if leaf does not match account", async function () {
    const { zkOracle, owner, addr1 } = await loadFixture(fullTreeFixture);
    let tx = zkOracle.exit(
      {
        index: 1,
        pubKeyX:
          "7794373982259243195870592346785104092432649697832080133780782253104282782817",
        pubKeyY:
          "16580021058669382711579094818964719809751621462486576197985799831318116474539",
        balance: ethers.utils.parseEther("0"),
      },
      [
        "6205836767976675972003616131281122597329838876688673966264947687549471156130",
        "6187350521517272486117148237635192271041670665937219625917563897233388432910",
        "6584420187527354519485242243152059973161633040054605148949951423401777995392",
      ],
      [1, 1]
    );

    await expect(tx).to.be.revertedWith("leaf does not match account");
  });

  it("exit should revert if account not included", async function () {
    const { zkOracle, owner, addr1 } = await loadFixture(fullTreeFixture);
    let tx = zkOracle.exit(
      {
        index: 0,
        pubKeyX:
          "7794373982259243195870592346785104092432649697832080133780782253104282782817",
        pubKeyY:
          "16580021058669382711579094818964719809751621462486576197985799831318116474539",
        balance: ethers.utils.parseEther("0"),
      },
      [
        "6205836767976675972003616131281122597329838876688673966264947687549471156130",
        "5187350521517272486117148237635192271041670665937219625917563897233388432910",
        "6584420187527354519485242243152059973161633040054605148949951423401777995392",
      ],
      [1, 1]
    );

    await expect(tx).to.be.revertedWith("account not included");
  });

  it("exit should set time", async function () {
    const { zkOracle, owner } = await loadFixture(fullTreeFixture);
    const tx = await zkOracle.exit(
      {
        index: 0,
        pubKeyX:
          "7794373982259243195870592346785104092432649697832080133780782253104282782817",
        pubKeyY:
          "16580021058669382711579094818964719809751621462486576197985799831318116474539",
        balance: ethers.utils.parseEther("0"),
      },
      [
        "6205836767976675972003616131281122597329838876688673966264947687549471156130",
        "6187350521517272486117148237635192271041670665937219625917563897233388432910",
        "6584420187527354519485242243152059973161633040054605148949951423401777995392",
      ],
      [1, 1]
    );
    const receipt = await tx.wait();

    const block = await ethers.provider.getBlock(receipt.blockNumber);
    const timestamp = ethers.BigNumber.from(block.timestamp);
    const exitDelay = await zkOracle.exitDelay();

    const expected = timestamp.add(exitDelay);
    const actual = await zkOracle.getExitTime(owner.address);

    expect(actual).to.equal(expected);
  });

  it("withdraw should revert when time not passed", async function () {
    const { zkOracle } = await loadFixture(deployFixture);
    let tx = zkOracle.withdraw(
      {
        index: 0,
        pubKeyX:
          "7794373982259243195870592346785104092432649697832080133780782253104282782817",
        pubKeyY:
          "16580021058669382711579094818964719809751621462486576197985799831318116474539",
        balance: ethers.utils.parseEther("0"),
      },
      [],
      []
    );
    await expect(tx).to.be.revertedWith("time not passed");
  });

  it("withdraw should revert if leaf does not match account", async function () {
    const { zkOracle, owner, addr1 } = await loadFixture(fullTreeFixture);
    let account = {
      index: 0,
      pubKeyX:
        "7794373982259243195870592346785104092432649697832080133780782253104282782817",
      pubKeyY:
        "16580021058669382711579094818964719809751621462486576197985799831318116474539",
      balance: ethers.utils.parseEther("0"),
    };

    await zkOracle.exit(
      account,
      [
        "6205836767976675972003616131281122597329838876688673966264947687549471156130",
        "6187350521517272486117148237635192271041670665937219625917563897233388432910",
        "6584420187527354519485242243152059973161633040054605148949951423401777995392",
      ],
      [1, 1]
    );
    let exitDelay = await zkOracle.exitDelay();
    exitDelay.add(1);

    await mine(exitDelay);

    account.index = 1;
    let tx = zkOracle.withdraw(
      account,
      [
        "6205836767976675972003616131281122597329838876688673966264947687549471156130",
        "6187350521517272486117148237635192271041670665937219625917563897233388432910",
        "6584420187527354519485242243152059973161633040054605148949951423401777995392",
      ],
      [1, 1]
    );
    await expect(tx).to.be.revertedWith("leaf does not match account");
  });
});
