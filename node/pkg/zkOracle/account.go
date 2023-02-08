package zkOracle

import (
	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
	"math/big"
)

const accountSize = 128

type Account struct {
	Index     *big.Int
	SecretKey *eddsa.PrivateKey
	Balance   *big.Int
}

// Serialize index ∥ pubkeyX ∥ pubkeyY ∥ balance
func (a *Account) Serialize() []byte {
	var b [accountSize]byte

	copy(b[:32], padOrTrim(a.Index.Bytes(), 32))

	publicKey := a.SecretKey.PublicKey

	var buf [32]byte
	buf = publicKey.A.X.Bytes()
	copy(b[32:], buf[:])
	buf = publicKey.A.Y.Bytes()
	copy(b[64:], buf[:])

	copy(b[96:], padOrTrim(a.Balance.Bytes(), 32))

	return b[:]
}
