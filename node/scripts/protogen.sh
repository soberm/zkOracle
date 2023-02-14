#!/bin/sh

cd "$(dirname "$0")" || exit 1

if ! which protoc >/dev/null; then
  echo "error: protoc not installed" >&2
  exit 1
fi

protoc -I ../api/proto/zkOracle --go_out=../pkg/zkOracle --go-grpc_out=../pkg/zkOracle ../api/proto/zkOracle/zkOracle.proto