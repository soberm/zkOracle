package zkOracle

import (
	"encoding/binary"
	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
	"github.com/consensys/gnark/frontend"
	"math/big"
)

const (
	voteSize = 8
)

type Vote struct {
	index     uint64
	request   *big.Int
	blockHash frontend.Variable
	sender    eddsa.PublicKey
	signature eddsa.Signature
}

func (v *Vote) Serialize() []byte {
	var b [voteSize]byte

	binary.BigEndian.PutUint64(b[:voteSize], v.index)
	return b[:]
}
