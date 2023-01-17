#!/bin/sh

cd "$(dirname "$0")" || exit 1

if ! which abigen >/dev/null; then
  echo "error: abigen not installed" >&2
  exit 1
fi

abigen --abi ../../contracts/artifacts/abi/contracts/ZKOracle.sol/ZKOracle.json --pkg zkOracle --type ZKOracleContract --out ../pkg/zkOracle/contract.go