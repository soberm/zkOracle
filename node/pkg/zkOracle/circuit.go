package zkOracle

import (
	"fmt"
	twisteded "github.com/consensys/gnark-crypto/ecc/twistededwards"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/accumulator/merkle"
	"github.com/consensys/gnark/std/algebra/twistededwards"
	"github.com/consensys/gnark/std/hash/mimc"
	"github.com/consensys/gnark/std/signature/eddsa"
)

const (
	batchSize = 4
	depth     = 3
)

type Circuit struct {
	//Aggregator AggregatorConstraints
	Root  frontend.Variable
	Votes [batchSize]VoteConstraints
}

type AggregatorConstraints struct {
	PublicKey eddsa.PublicKey
}

type VoteConstraints struct {
	PublicKey         eddsa.PublicKey
	MerkleProof       [depth]frontend.Variable
	MerkleProofHelper [depth - 1]frontend.Variable
	Signature         eddsa.Signature
	Result            frontend.Variable
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

		hFunc.Write(vote.PublicKey.A.X)
		hFunc.Write(vote.PublicKey.A.Y)
		api.AssertIsEqual(hFunc.Sum(), vote.MerkleProof[0])

		hFunc.Reset()

		merkle.VerifyProof(api, hFunc, c.Root, vote.MerkleProof[:], vote.MerkleProofHelper[:])

		if err := eddsa.Verify(curve, vote.Signature, vote.Result, vote.PublicKey, &hFunc); err != nil {
			return fmt.Errorf("verify eddsa: %w", err)
		}
	}

	return nil
}
