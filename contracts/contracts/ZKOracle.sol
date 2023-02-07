// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import "hardhat/console.sol";
import "./MerkleTree.sol";

contract ZKOracle {
    MerkleTree merkleTree;

    struct Account {
        uint256 index;
        PublicKey pubKey;
        uint256 balance;
    }

    struct PublicKey {
        uint256 x;
        uint256 y;
    }

    uint256 public exitDelay = 604800;

    mapping(uint256 => address) accounts;
    mapping(address => uint256) exitTimes;

    event Registered(address indexed sender);
    event Replaced(address indexed sender, address indexed replaced);
    event Exiting(address indexed sender);
    event Withdrawn(address indexed sender);

    constructor(address merkleTreeAddress) {
        merkleTree = MerkleTree(merkleTreeAddress);
    }

    function register(PublicKey memory publicKey) public payable {
        Account memory account = Account(
            merkleTree.getNextLeafIndex(),
            publicKey,
            msg.value
        );
        accounts[account.index] = msg.sender;
        uint256 accountHash = hashAccount(account);

        uint[] memory input = new uint[](1);
        input[0] = accountHash;
        uint256 h = merkleTree.hash(input);

        merkleTree.insert(h);

        emit Registered(msg.sender);
    }

    function replace(
        PublicKey memory publicKey,
        Account memory toReplace,
        uint256[] memory path,
        uint256[] memory helper
    ) public payable {
        require(msg.value > toReplace.balance, "value too low");
        require(merkleTree.verify(path, helper), "account not included");
        require(
            path[0] == hashAccount(toReplace),
            "leaf does not match account"
        );

        Account memory replaced = Account(
            toReplace.index,
            publicKey,
            msg.value
        );

        merkleTree.update(hashAccount(replaced), path, helper);
        payable(accounts[toReplace.index]).transfer(toReplace.balance);
        emit Replaced(msg.sender, accounts[toReplace.index]);
        accounts[toReplace.index] = msg.sender;
    }

    function exit(
        Account memory account,
        uint256[] memory path,
        uint256[] memory helper
    ) public {
        require(accounts[account.index] == msg.sender, "wrong sender address");
        require(path[0] == hashAccount(account), "leaf does not match account");
        require(merkleTree.verify(path, helper), "account not included");

        exitTimes[msg.sender] = block.timestamp + exitDelay;
        emit Exiting(msg.sender);
    }

    function withdraw(
        Account memory account,
        uint256[] memory path,
        uint256[] memory helper
    ) public {
        require(block.number < exitTimes[msg.sender], "time not passed");
        require(accounts[account.index] == msg.sender, "wrong sender address");
        require(path[0] == hashAccount(account), "leaf does not match account");
        require(merkleTree.verify(path, helper), "account not included");

        payable(msg.sender).transfer(account.balance);
        delete accounts[account.index];

        Account memory empty = Account(account.index, account.pubKey, 0);
        merkleTree.update(hashAccount(empty), path, helper);
        emit Withdrawn(msg.sender);
    }

    function hashAccount(Account memory account) public view returns (uint256) {
        uint[] memory input = new uint[](4);
        input[0] = account.index;
        input[1] = account.pubKey.x;
        input[2] = account.pubKey.y;
        input[3] = account.balance;
        return merkleTree.hash(input);
    }

    function getExitTime(address addr) public view returns (uint256) {
        return exitTimes[addr];
    }
}
