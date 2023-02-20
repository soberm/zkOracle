package zkOracle

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254/fp"
	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
	"github.com/consensys/gnark/backend"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/algebra/twistededwards"
	eddsa2 "github.com/consensys/gnark/std/signature/eddsa"
	"math/big"
)

type Aggregator struct {
	index            uint64
	privateKey       *eddsa.PrivateKey
	state            *State
	votePool         *VotePool
	constraintSystem frontend.CompiledConstraintSystem
	provingKey       groth16.ProvingKey
	verifyingKey     groth16.VerifyingKey
}

func NewAggregator(index uint64, privateKey *eddsa.PrivateKey, state *State, votePool *VotePool, constraintSystem frontend.CompiledConstraintSystem, provingKey groth16.ProvingKey, verifyingKey groth16.VerifyingKey) *Aggregator {
	return &Aggregator{
		index:            index,
		privateKey:       privateKey,
		state:            state,
		votePool:         votePool,
		constraintSystem: constraintSystem,
		provingKey:       provingKey,
		verifyingKey:     verifyingKey,
	}
}

func (a *Aggregator) Aggregate(ctx context.Context) error {
	for requestNumber := range a.votePool.sink {
		logger.Info().Uint64("requestNumber", requestNumber).Msg("aggregate votes")
		votes, err := a.votePool.getVotes(requestNumber)
		if err != nil {
			return fmt.Errorf("get votes: %w", err)
		}
		err = a.ProcessVotes(votes)
		if err != nil {
			return fmt.Errorf("process votes: %w", err)
		}
	}

	return nil
}

func (a *Aggregator) ProcessVotes(votes []*Vote) error {

	preStateRoot, aggregatorProof, aggregatorHelper, err := a.state.MerkleProof(a.index)
	if err != nil {
		return fmt.Errorf("aggregator merkle proof: %w", err)
	}

	aggregatorAccount, err := a.state.ReadAccount(a.index)
	if err != nil {
		return fmt.Errorf("read aggregator account: %w", err)
	}

	logger.Info().Str("preStateRoot", hex.EncodeToString(preStateRoot)).Msg("Start processing")

	aggregatorConstraints := AggregatorConstraints{
		Index:             a.index,
		Seed:              twistededwards.Point{X: 0, Y: 1},
		SecretKey:         a.privateKey.Bytes()[fp.Bytes : 2*fp.Bytes],
		Balance:           new(big.Int).Set(aggregatorAccount.Balance),
		MerkleProof:       aggregatorProof,
		MerkleProofHelper: aggregatorHelper,
	}

	logger.Info().
		Uint64("Index", aggregatorAccount.Index.Uint64()).
		Uint64("balance", new(big.Int).Set(aggregatorAccount.Balance).Uint64()).
		Msg("aggregator constraints")

	//fmt.Printf("MerkleProof: %v\n", aggregatorProof)

	aggregatorAccount.Balance.Add(aggregatorAccount.Balance, big.NewInt(AggregatorReward))
	err = a.state.WriteAccount(aggregatorAccount)
	if err != nil {
		return fmt.Errorf("write account: %w", err)
	}

	var validatorConstraints [nbAccounts]ValidatorConstraints
	for i, vote := range votes {
		validatorAccount, err := a.state.ReadAccount(vote.Index)
		if err != nil {
			return fmt.Errorf("read aggregator account: %w", err)
		}

		var publicKey eddsa2.PublicKey
		var signature eddsa2.Signature

		publicKey.Assign(ecc.BN254, validatorAccount.PublicKey.Bytes())
		signature.Assign(ecc.BN254, vote.Signature.Bytes())

		logger.Info().Str("PublicKey", hex.EncodeToString(validatorAccount.PublicKey.Bytes())).
			Str("Signature", hex.EncodeToString(vote.Signature.Bytes())).
			Msg("Test")

		_, proof, helper, err := a.state.MerkleProof(validatorAccount.Index.Uint64())
		if err != nil {
			return fmt.Errorf("validator merkle proof: %w", err)
		}

		validatorConstraints[i] = ValidatorConstraints{
			Index:             validatorAccount.Index,
			PublicKey:         publicKey,
			Balance:           new(big.Int).Set(validatorAccount.Balance), //passed by reference
			MerkleProof:       proof,
			MerkleProofHelper: helper,
			Signature:         signature,
			BlockHash:         vote.BlockHash.Bytes(),
		}

		logger.Info().
			Uint64("Index", validatorAccount.Index.Uint64()).
			Uint64("balance", new(big.Int).Set(validatorAccount.Balance).Uint64()).
			Msg("validator constraints")

		//fmt.Printf("MerkleProof: %v\n", proof)

		validatorAccount.Balance.Add(validatorAccount.Balance, big.NewInt(ValidatorReward))
		err = a.state.WriteAccount(validatorAccount)
		if err != nil {
			return fmt.Errorf("write account: %w", err)
		}
	}

	postStateRoot, err := a.state.Root()
	if err != nil {
		return fmt.Errorf("state root: %w", err)
	}

	logger.Info().Str("PostStateRoot", hex.EncodeToString(postStateRoot)).Msg("post state root")

	assignment := AggregationCircuit{
		PreStateRoot:  preStateRoot,
		PostStateRoot: postStateRoot,
		BlockHash:     votes[0].BlockHash.Bytes(),
		Request:       votes[0].Request,
		Aggregator:    aggregatorConstraints,
		Validators:    validatorConstraints,
	}

	witness, err := frontend.NewWitness(&assignment, ecc.BN254)
	if err != nil {
		return fmt.Errorf("create witness: %w", err)
	}

	p, err := groth16.Prove(a.constraintSystem, a.provingKey, witness, backend.IgnoreSolverError())
	if err != nil {
		return fmt.Errorf("prove: wv", err)
	}

	pw, err := witness.Public()
	if err != nil {
		return fmt.Errorf("public witness: %w", err)
	}

	err = groth16.Verify(p, a.verifyingKey, pw)
	if err != nil {
		return fmt.Errorf("verify proof: %w", err)
	}

	return nil
}
