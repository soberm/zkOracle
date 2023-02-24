package zkOracle

import (
	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

const (
	voteSize = 96
)

type Vote struct {
	Index     uint64
	Request   *big.Int
	BlockHash common.Hash
	Sender    eddsa.PublicKey
	Signature eddsa.Signature
}

func (v *Vote) Serialize() []byte {

	var b [voteSize]byte

	copy(b[:32], PadOrTrim(big.NewInt(int64(v.Index)).Bytes(), 32))
	copy(b[32:64], PadOrTrim(v.Request.Bytes(), 32))
	copy(b[64:voteSize], v.BlockHash.Bytes())

	return b[:]
}
