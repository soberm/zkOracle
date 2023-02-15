package zkOracle

import (
	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

const (
	voteSize = 32
)

type Vote struct {
	index     uint64
	request   *big.Int
	blockHash common.Hash
	sender    eddsa.PublicKey
	signature eddsa.Signature
}

func (v *Vote) Serialize() []byte {
	var b [voteSize]byte

	copy(b[:common.HashLength], v.blockHash.Bytes())
	return b[:]
}
