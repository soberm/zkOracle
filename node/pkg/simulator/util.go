package simulator

import (
	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
	"math/big"
	"node/pkg/zkOracle"
)

func publicKeyToZKOraclePublicKey(publicKey *eddsa.PublicKey) *zkOracle.ZKOraclePublicKey {
	return &zkOracle.ZKOraclePublicKey{
		X: publicKey.A.X.ToBigIntRegular(big.NewInt(0)),
		Y: publicKey.A.Y.ToBigIntRegular(big.NewInt(0)),
	}
}

func accountToZKOracleAccount(account *zkOracle.Account) *zkOracle.ZKOracleAccount {
	return &zkOracle.ZKOracleAccount{
		Index:   account.Index,
		PubKey:  *publicKeyToZKOraclePublicKey(account.PublicKey),
		Balance: account.Balance,
	}
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
