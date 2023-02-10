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
)

type Circuit struct {
	Root       frontend.Variable `gnark:",public"`
	BlockHash  frontend.Variable `gnark:",public"`
	Aggregator AggregatorConstraints
	Validators [nbAccounts]ValidatorConstraints
}

type AggregatorConstraints struct {
	Index             frontend.Variable    `gnark:",public"`
	Seed              twistededwards.Point `gnark:",public"`
	SecretKey         frontend.Variable
	Balance           frontend.Variable
	MerkleProof       [depth]frontend.Variable
	MerkleProofHelper [depth - 1]frontend.Variable
}

type ValidatorConstraints struct {
	Index             frontend.Variable
	PublicKey         eddsa.PublicKey
	Balance           frontend.Variable
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

	hFunc, err := mimc.NewMiMC(api)
	if err != nil {
		return fmt.Errorf("mimc: %w", err)
	}

	//Compute next Seed
	c.Aggregator.Seed = curve.ScalarMul(c.Aggregator.Seed, c.Aggregator.SecretKey)

	// Compute aggregator public key
	base := curve.Params().Base
	g := twistededwards.Point{X: base[0], Y: base[1]}
	pubKey := curve.ScalarMul(g, c.Aggregator.SecretKey)
	curve.AssertIsOnCurve(pubKey)

	// Verify that the public key from the Merkle proof matches the computed public key of the aggregator
	hFunc.Write(c.Aggregator.Index)
	hFunc.Write(pubKey.X)
	hFunc.Write(pubKey.Y)
	hFunc.Write(c.Aggregator.Balance)
	api.AssertIsEqual(hFunc.Sum(), c.Aggregator.MerkleProof[0])
	api.Println(hFunc.Sum())
	hFunc.Reset()

	api.Println(c.Root)
	// Check aggregator included
	merkle.VerifyProof(api, hFunc, c.Root, c.Aggregator.MerkleProof[:], c.Aggregator.MerkleProofHelper[:])

	hFunc.Reset()

	updatedBalance := api.Add(c.Aggregator.Balance, 50)
	api.Println(updatedBalance)

	//Reward the aggregator
	hFunc.Write(c.Aggregator.Index)
	hFunc.Write(pubKey.X)
	hFunc.Write(pubKey.Y)
	hFunc.Write(updatedBalance)
	c.Aggregator.MerkleProof[0] = hFunc.Sum()

	//Compute new intermediate root
	hFunc.Reset()
	intermediateRoot := ComputeRootFromPath(api, hFunc, c.Aggregator.MerkleProof[:], c.Aggregator.MerkleProofHelper[:])
	api.Println(intermediateRoot)

	for _, validator := range c.Validators {
		hFunc.Reset()

		//Verify that the account matches the leaf
		hFunc.Write(validator.Index)
		hFunc.Write(validator.PublicKey.A.X)
		hFunc.Write(validator.PublicKey.A.Y)
		hFunc.Write(validator.Balance)
		api.AssertIsEqual(hFunc.Sum(), validator.MerkleProof[0])

		hFunc.Reset()

		//Check validator included
		merkle.VerifyProof(api, hFunc, intermediateRoot, validator.MerkleProof[:], validator.MerkleProofHelper[:])

		if err := eddsa.Verify(curve, validator.Signature, validator.BlockHash, validator.PublicKey, &hFunc); err != nil {
			return fmt.Errorf("verify eddsa: %w", err)
		}

		api.AssertIsEqual(c.BlockHash, validator.BlockHash)

		//Reward the validator
		/*		hFunc.Write(validator.Index)
				hFunc.Write(validator.PublicKey.A.X)
				hFunc.Write(validator.PublicKey.A.Y)
				hFunc.Write(api.Add(validator.Balance, 5))
				validator.MerkleProof[0] = hFunc.Sum()*/

		//Compute new intermediate root
		//hFunc.Reset()
		//intermediateRoot = ComputeRootFromPath(api, hFunc, c.Aggregator.MerkleProof[:], c.Aggregator.MerkleProofHelper[:])
	}

	return nil
}
