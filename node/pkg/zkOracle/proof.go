package zkOracle

import (
	"bytes"
	"fmt"
	"github.com/consensys/gnark-crypto/ecc/bn254/fp"
	"github.com/consensys/gnark/backend/groth16"
	"math/big"
)

type EthereumProof struct {
	a [2]*big.Int
	b [2][2]*big.Int
	c [2]*big.Int
}

func ProofToEthereumProof(p groth16.Proof) (*EthereumProof, error) {

	var proof EthereumProof

	var buf bytes.Buffer
	_, err := p.WriteRawTo(&buf)
	if err != nil {
		return nil, fmt.Errorf("write raw proof to: %w", err)
	}
	proofBytes := buf.Bytes()

	proof.a[0] = new(big.Int).SetBytes(proofBytes[fp.Bytes*0 : fp.Bytes*1])
	proof.a[1] = new(big.Int).SetBytes(proofBytes[fp.Bytes*1 : fp.Bytes*2])
	proof.b[0][0] = new(big.Int).SetBytes(proofBytes[fp.Bytes*2 : fp.Bytes*3])
	proof.b[0][1] = new(big.Int).SetBytes(proofBytes[fp.Bytes*3 : fp.Bytes*4])
	proof.b[1][0] = new(big.Int).SetBytes(proofBytes[fp.Bytes*4 : fp.Bytes*5])
	proof.b[1][1] = new(big.Int).SetBytes(proofBytes[fp.Bytes*5 : fp.Bytes*6])
	proof.c[0] = new(big.Int).SetBytes(proofBytes[fp.Bytes*6 : fp.Bytes*7])
	proof.c[1] = new(big.Int).SetBytes(proofBytes[fp.Bytes*7 : fp.Bytes*8])

	return &proof, nil
}
