# Cross-Chain Oracle Using an Off-Chain Aggregation Mechanism Based on zk-SNARKs

This project contains the source code for the paper "Cross-Blockchain Communication Using Oracles With an
Off-Chain Aggregation Mechanism Based on zk-SNARKs." We provide a prototypical implementation of the smart contracts and the Oracle node.

## Prerequisites

You need to have the following software installed:

* [Golang](https://golang.org/doc/install) (version 1.19)
* [Node.js](https://nodejs.org/) (version >= 19.4.0)
* [Hardhat](https://hardhat.org/) (version >= 2.11.1)
* [Solidity](https://docs.soliditylang.org/en/latest/installing-solidity.html) (^0.8.0)

## Installation

### Constraint System Setup

1. Change into the following directory: `cd node/cmd/compiler`
2. Install all dependencies: `go mod download`
3. Adapt the parameters of the circuits in the respective src files
4. Build the constraint system setup: `go build -o compiler`
5. Run the constraint system setup: `./compiler -b ./build`
6. Update the verifier contracts using the generated contracts

### Smart Contracts

1. Change into the contract directory: `cd contracts/`
2. Install all dependencies: `npm install`
3. Compile contracts: `hardhat compile`
4. Update the parameters in ./scripts/deploy.ts
5. Deploy contracts: `hardhat run --network <your_network> ./scripts/deploy.ts`

### Node

1. Change into the operator directory: `cd node/cmd/zkOracle`
2. Install all dependencies: `go mod download`
3. Build the node: `go build -o node`
4. Run the node: `./node -c ./configs/config.json`

### Evaluation

For the evaluation of the prototype, we also provide a simulated version that allows measuring the gas consumption, proving time and memory usage in node/cmd/simulator

## Contributing

This is a research prototype. We welcome anyone to contribute. File a bug report or submit feature requests through the issue tracker. If you want to contribute, feel free to submit a pull request.

## Acknowledgement

The financial support by the Austrian Federal Ministry for Digital and Economic Affairs, the National Foundation for Research, Technology and Development as well as the Christian Doppler Research Association is gratefully acknowledged.

## Licence

This project is licensed under the [MIT License](LICENSE).
