import { HardhatUserConfig } from "hardhat/config";
import "@nomicfoundation/hardhat-toolbox";
import "hardhat-abi-exporter";
//import "@dlsl/hardhat-gobind"

const config: HardhatUserConfig = {
  solidity: "0.8.17",
  abiExporter: {
    path: "./artifacts/abi",
    runOnCompile: true,
    clear: true,
    flat: false,
    only: ["ZKOracle"],
    spacing: 2,
    pretty: false,
  },
};

export default config;
