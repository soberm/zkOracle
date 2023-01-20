package zkOracle

import (
	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
	"math/big"
)

const accountSize = 96

type Account struct {
	Index     *big.Int
	SecretKey *eddsa.PrivateKey
}

func (a *Account) Serialize() []byte {
	var b [accountSize]byte

	copy(b[:32], padOrTrim(a.Index.Bytes(), 32))

	publicKey := a.SecretKey.PublicKey

	var buf [32]byte
	buf = publicKey.A.X.Bytes()
	copy(b[32:], buf[:])
	buf = publicKey.A.Y.Bytes()
	copy(b[64:], buf[:])

	return b[:]
}
