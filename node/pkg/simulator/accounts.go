package simulator

import (
	"fmt"
	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
	"math/big"
	"math/rand"
	"node/pkg/zkOracle"
	"time"
)

func generatePrivateKeys(number int) ([]*eddsa.PrivateKey, error) {
	privateKeys := make([]*eddsa.PrivateKey, number)
	for i := 0; i < number; i++ {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		sk, err := eddsa.GenerateKey(r)
		if err != nil {
			return nil, fmt.Errorf("eddsa generate key: %w", err)
		}
		privateKeys[i] = sk
	}
	return privateKeys, nil
}

func createAccounts(privateKeys []*eddsa.PrivateKey) ([]*zkOracle.Account, error) {
	accounts := make([]*zkOracle.Account, len(privateKeys))
	for i, privateKey := range privateKeys {
		accounts[i] = &zkOracle.Account{
			Index:     big.NewInt(int64(i)),
			PublicKey: &privateKey.PublicKey,
			Balance:   big.NewInt(0),
		}
	}
	return accounts, nil
}
