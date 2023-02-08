package zkOracle

import (
	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
	"github.com/consensys/gnark/frontend"
	"math/big"
)

type Vote struct {
	index     uint64
	request   *big.Int
	blockHash frontend.Variable
	sender    eddsa.PublicKey
	signature eddsa.Signature
}
