package zkOracle

import (
	"fmt"
	twisteded "github.com/consensys/gnark-crypto/ecc/twistededwards"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/algebra/twistededwards"
	"github.com/consensys/gnark/std/hash/mimc"
	"github.com/consensys/gnark/std/signature/eddsa"
)

const (
	batchSize = 3
	depth     = 2
)

type Circuit struct {
	//Aggregator AggregatorConstraints
	Votes [batchSize]VoteConstraints
}

type AggregatorConstraints struct {
	PublicKey eddsa.PublicKey
}

type VoteConstraints struct {
	PublicKey eddsa.PublicKey
	//MerkleProofValidator [batchSize][depth]frontend.Variable
	Signature eddsa.Signature
	Result    frontend.Variable
}

func (c *Circuit) Define(api frontend.API) error {
	curve, err := twistededwards.NewEdCurve(api, twisteded.BN254)
	if err != nil {
		return fmt.Errorf("edwards curve: %w", err)
	}

	for _, vote := range c.Votes {
		hFunc, err := mimc.NewMiMC(api)
		if err != nil {
			return fmt.Errorf("mimc: %w", err)
		}

		if err := eddsa.Verify(curve, vote.Signature, vote.Result, vote.PublicKey, &hFunc); err != nil {
			return fmt.Errorf("verify eddsa: %w", err)
		}
	}

	return nil
}
