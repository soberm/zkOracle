package zkOracle

import (
	"encoding/binary"
	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

const (
	voteSize = 40
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
	binary.BigEndian.PutUint64(b[:8], v.index)
	copy(b[8:common.HashLength], v.blockHash.Bytes())
	return b[:]
}
