package main

import (
	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
	"math/big"
	"node/pkg/zkOracle"
)

func PublicKeyToZKOraclePublicKey(publicKey *eddsa.PublicKey) *zkOracle.ZKOraclePublicKey {
	return &zkOracle.ZKOraclePublicKey{
		X: publicKey.A.X.ToBigIntRegular(big.NewInt(0)),
		Y: publicKey.A.Y.ToBigIntRegular(big.NewInt(0)),
	}
}

func AccountToZKOracleAccount(account *zkOracle.Account) *zkOracle.ZKOracleAccount {
	return &zkOracle.ZKOracleAccount{
		Index:   account.Index,
		PubKey:  *PublicKeyToZKOraclePublicKey(account.PublicKey),
		Balance: account.Balance,
	}
}
