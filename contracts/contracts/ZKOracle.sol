// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

// import "hardhat/console.sol";
import "./MerkleTree.sol";

contract ZKOracle {
    MerkleTree merkleTree;

    event Registered(address sender);

    constructor(address merkleTreeAddress) {
        merkleTree = MerkleTree(merkleTreeAddress);
    }

    function register(uint256 pubKeyX, uint256 pubKeyY) public payable {
        uint[] memory input = new uint[](4);
        input[0] = merkleTree.getNextLeafIndex();
        input[1] = pubKeyX;
        input[2] = pubKeyY;
        input[3] = msg.value;

        uint256 leaf = merkleTree.hash(input);
        merkleTree.insert(leaf);

        emit Registered(msg.sender);
    }
}
