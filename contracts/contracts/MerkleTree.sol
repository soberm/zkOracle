// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import "./MiMC.sol";
import "hardhat/console.sol";

contract MerkleTree {
    uint256 ZERO_VALUE =
        4555114089170143013007615382799372902997177870479602349537353593038812875418;

    uint256 levels;
    uint256 public root;
    mapping(uint256 => uint256) public filledSubtrees;
    uint256 public nextIndex;

    event Inserted(uint256 index);

    constructor(uint256 _levels) {
        require(_levels > 0, "_levels should be greater than zero");
        require(_levels < 8, "_levels should be less than 8");

        levels = _levels;

        for (uint32 i = 0; i < levels; i++) {
            filledSubtrees[i] = zeros(i);
        }

        root = zeros(levels - 1);
    }

    function hash(uint256[] memory data) public pure returns (uint256) {
        return MiMC.Hash(data);
    }

    function insert(uint256 leaf) public {
        require(nextIndex != 2 ** levels, "tree is full");

        uint256 left;
        uint256 right;
        uint256 leafIndex = nextIndex;
        uint256 currentIndex = nextIndex;
        uint256 currentHash = leaf;

        for (uint i = 0; i < levels; i++) {
            if (currentIndex % 2 == 0) {
                left = currentHash;
                right = zeros(i);
                filledSubtrees[i] = currentHash;
            } else {
                left = filledSubtrees[i];
                right = currentHash;
            }
            uint[] memory input = new uint[](2);
            input[0] = left;
            input[1] = right;
            currentHash = MiMC.Hash(input);
            currentIndex /= 2;
        }

        root = currentHash;

        nextIndex += 1;
        emit Inserted(leafIndex);
    }

    function getRoot() public view returns (uint256) {
        return root;
    }

    function getLevels() public view returns (uint256) {
        return levels;
    }

    function zeros(uint256 i) public pure returns (uint256) {
        if (i == 0)
            return
                4555114089170143013007615382799372902997177870479602349537353593038812875418;
        else if (i == 1)
            return
                9614759978327623946452646332910600180945773348102064399025967221305784063943;
        else if (i == 2)
            return
                15762506290347708512348905356059207185046946323941989403490412292643733744343;
        else if (i == 3)
            return
                2078761282949659850987695139228042067769933906673781403014209677812047702550;
        else if (i == 4)
            return
                20395412135670005113982952294980644860334762516174975965396550838918688642133;
        else if (i == 5)
            return
                17560454953585356688949527489694249319830065182192048221544285096802159445652;
        else if (i == 6)
            return
                20019762671335178393512154978075455201849332419823879662510519485824706883752;
        else if (i == 7)
            return
                12065157948427223398688325372361960271507319753018415581972466307863230644170;
        else revert("index out of bounds");
    }
}
