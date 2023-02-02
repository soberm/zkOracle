// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

// import "hardhat/console.sol";
import "./MerkleTree.sol";

contract ZKOracle {
    MerkleTree merkleTree;

    struct Account {
        uint256 index;
        uint256 pubKeyX;
        uint256 pubKeyY;
        uint256 balance;
    }

    uint256 public exitDelay = 604800;

    mapping(uint256 => address) accounts;
    mapping(address => uint256) exitTimes;

    event Registered(address indexed sender);
    event Exiting(address indexed sender);

    constructor(address merkleTreeAddress) {
        merkleTree = MerkleTree(merkleTreeAddress);
    }

    function register(uint256 pubKeyX, uint256 pubKeyY) public payable {
        Account memory account = Account(
            merkleTree.getNextLeafIndex(),
            pubKeyX,
            pubKeyY,
            msg.value
        );
        accounts[account.index] = msg.sender;
        uint256 accountHash = hashAccount(account);

        merkleTree.insert(accountHash);

        emit Registered(msg.sender);
    }

    function replace(
        uint256 pubKeyX,
        uint256 pubKeyY,
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

        merkleTree.update(0, path, helper);
        payable(accounts[toReplace.index]).transfer(toReplace.balance);
    }

    function exit() public {
        exitTimes[msg.sender] = block.timestamp + exitDelay;
        emit Exiting(msg.sender);
    }

    function withdraw(
        Account memory account,
        uint256[] memory path,
        uint256[] memory helper
    ) public {
        require(block.number < exitTimes[msg.sender], "time not passed");
        require(accounts[account.index] == msg.sender, "wrong index");
        require(merkleTree.verify(path, helper), "account not included");
        require(path[0] == hashAccount(account), "leaf does not match account");

        payable(msg.sender).transfer(account.balance);
        delete accounts[account.index];
        //TODO: Set the accounts balance to 0
        //TODO: Verify Merkle proof and recompute root
    }

    function hashAccount(Account memory account) public view returns (uint256) {
        uint[] memory input = new uint[](4);
        input[0] = account.index;
        input[1] = account.pubKeyX;
        input[2] = account.pubKeyY;
        input[3] = account.balance;
        return merkleTree.hash(input);
    }

    function getExitTime(address addr) public view returns (uint256) {
        return exitTimes[addr];
    }
}
