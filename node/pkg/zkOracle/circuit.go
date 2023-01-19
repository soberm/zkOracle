package zkOracle

import (
	"fmt"
	edwards "github.com/consensys/gnark-crypto/ecc/twistededwards"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/accumulator/merkle"
	"github.com/consensys/gnark/std/algebra/twistededwards"
	"github.com/consensys/gnark/std/hash/mimc"
	"github.com/consensys/gnark/std/signature/eddsa"
)

const (
	nbAccounts = 4
	depth      = 3
	threshold  = 2
)

type Circuit struct {
	Aggregator AggregatorConstraints
	Root       frontend.Variable `gnark:",public"`
	BlockHash  frontend.Variable `gnark:",public"`
	//Seed      twistededwards.Point `gnark:",public"`
	//Epoch frontend.Variable `gnark:",public"`
	Votes [nbAccounts]VoteConstraints
}

type AggregatorConstraints struct {
	//	PublicKey         eddsa.PublicKey
	SecretKey frontend.Variable
	//	MerkleProof       [depth]frontend.Variable
	//	MerkleProofHelper [depth - 1]frontend.Variable
}

type VoteConstraints struct {
	PublicKey         eddsa.PublicKey
	MerkleProof       [depth]frontend.Variable
	MerkleProofHelper [depth - 1]frontend.Variable
	Signature         eddsa.Signature
	BlockHash         frontend.Variable
}

func (c *Circuit) Define(api frontend.API) error {
	curve, err := twistededwards.NewEdCurve(api, edwards.BN254)
	if err != nil {
		return fmt.Errorf("edwards curve: %w", err)
	}

	checkOwnership(api, curve, c.Aggregator.SecretKey, c.Votes[0].PublicKey)

	count := frontend.Variable(0)
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

		if err := eddsa.Verify(curve, vote.Signature, vote.BlockHash, vote.PublicKey, &hFunc); err != nil {
			return fmt.Errorf("verify eddsa: %w", err)
		}

		count = api.Select(api.Cmp(c.BlockHash, vote.BlockHash), count, api.Add(count, 1))
	}

	api.AssertIsEqual(api.Cmp(count, threshold), 1)

	return nil
}
