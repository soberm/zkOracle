package zkOracle

import (
	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

const (
// voteSize = 40
)

type Vote struct {
	Index     uint64
	Request   *big.Int
	BlockHash common.Hash
	Sender    eddsa.PublicKey
	Signature eddsa.Signature
}

func (v *Vote) Serialize() []byte {
	var b [96]byte
	//binary.BigEndian.PutUint64(b[:8], v.Index)
	copy(b[:32], padOrTrim(big.NewInt(int64(v.Index)).Bytes(), 32))
	copy(b[32:64], padOrTrim(v.Request.Bytes(), 32))
	copy(b[64:96], v.BlockHash.Bytes())
	return b[:]
}
