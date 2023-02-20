package zkOracle

import (
	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
	"math/big"
)

const accountSize = 128

type Account struct {
	Index     *big.Int
	PublicKey *eddsa.PublicKey
	Balance   *big.Int
}

// Serialize Index ∥ pubkeyX ∥ pubkeyY ∥ balance
func (a *Account) Serialize() []byte {
	var b [accountSize]byte

	copy(b[:32], padOrTrim(a.Index.Bytes(), 32))

	var buf [32]byte
	buf = a.PublicKey.A.X.Bytes()
	copy(b[32:], buf[:])
	buf = a.PublicKey.A.Y.Bytes()
	copy(b[64:], buf[:])

	copy(b[96:], padOrTrim(a.Balance.Bytes(), 32))

	return b[:]
}

func (a *Account) Deserialize(data []byte) {

	a.Index = big.NewInt(0).SetBytes(data[:32])

	a.PublicKey = new(eddsa.PublicKey)

	a.PublicKey.A.X.SetZero()
	a.PublicKey.A.Y.SetOne()

	a.PublicKey.A.X.SetBytes(data[32:64])
	a.PublicKey.A.Y.SetBytes(data[64:96])

	a.Balance = big.NewInt(0).SetBytes(data[96:accountSize])
}
