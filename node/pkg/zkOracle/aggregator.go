package zkOracle

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254/fp"
	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
	"github.com/consensys/gnark/backend"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/algebra/twistededwards"
	eddsa2 "github.com/consensys/gnark/std/signature/eddsa"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
)

type Proof struct {
	a [2]*big.Int
	b [2][2]*big.Int
	c [2]*big.Int
}

type Aggregator struct {
	index            uint64
	privateKey       *eddsa.PrivateKey
	state            *State
	votePool         *VotePool
	constraintSystem frontend.CompiledConstraintSystem
	provingKey       groth16.ProvingKey
	verifyingKey     groth16.VerifyingKey
	contract         *ZKOracleContract
	chainID          *big.Int
	ecdsaPrivateKey  *ecdsa.PrivateKey
	ethClient        *ethclient.Client
}

func NewAggregator(
	index uint64,
	privateKey *eddsa.PrivateKey,
	state *State,
	votePool *VotePool,
	constraintSystem frontend.CompiledConstraintSystem,
	provingKey groth16.ProvingKey,
	verifyingKey groth16.VerifyingKey,
	contract *ZKOracleContract,
	chainID *big.Int,
	ecdsaPrivateKey *ecdsa.PrivateKey,
	ethClient *ethclient.Client,
) *Aggregator {
	return &Aggregator{
		index:            index,
		privateKey:       privateKey,
		state:            state,
		votePool:         votePool,
		constraintSystem: constraintSystem,
		provingKey:       provingKey,
		verifyingKey:     verifyingKey,
		contract:         contract,
		chainID:          chainID,
		ecdsaPrivateKey:  ecdsaPrivateKey,
		ethClient:        ethClient,
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

	aggregatorConstraints := AggregatorConstraints{
		Index:             a.index,
		Seed:              twistededwards.Point{X: 0, Y: 1},
		SecretKey:         a.privateKey.Bytes()[fp.Bytes : 2*fp.Bytes],
		Balance:           new(big.Int).Set(aggregatorAccount.Balance),
		MerkleProof:       aggregatorProof,
		MerkleProofHelper: aggregatorHelper,
	}

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

	var proof Proof

	// get proof bytes
	var buf bytes.Buffer
	_, err = p.WriteRawTo(&buf)
	if err != nil {
		return fmt.Errorf("write proof: %w", err)
	}
	proofBytes := buf.Bytes()

	proof.a[0] = new(big.Int).SetBytes(proofBytes[fp.Bytes*0 : fp.Bytes*1])
	proof.a[1] = new(big.Int).SetBytes(proofBytes[fp.Bytes*1 : fp.Bytes*2])
	proof.b[0][0] = new(big.Int).SetBytes(proofBytes[fp.Bytes*2 : fp.Bytes*3])
	proof.b[0][1] = new(big.Int).SetBytes(proofBytes[fp.Bytes*3 : fp.Bytes*4])
	proof.b[1][0] = new(big.Int).SetBytes(proofBytes[fp.Bytes*4 : fp.Bytes*5])
	proof.b[1][1] = new(big.Int).SetBytes(proofBytes[fp.Bytes*5 : fp.Bytes*6])
	proof.c[0] = new(big.Int).SetBytes(proofBytes[fp.Bytes*6 : fp.Bytes*7])
	proof.c[1] = new(big.Int).SetBytes(proofBytes[fp.Bytes*7 : fp.Bytes*8])

	auth, err := bind.NewKeyedTransactorWithChainID(a.ecdsaPrivateKey, a.chainID)
	if err != nil {
		return fmt.Errorf("new transactor: %w", err)
	}
	auth.GasPrice = big.NewInt(20000000000)

	var blockHash [32]byte
	copy(blockHash[:], votes[0].BlockHash.Bytes()[:32])

	tx, err := a.contract.SubmitBlock(
		auth,
		big.NewInt(0).SetUint64(a.index),
		votes[0].Request,
		blockHash,
		big.NewInt(0).SetBytes(postStateRoot),
		proof.a,
		proof.b,
		proof.c,
	)
	if err != nil {
		return fmt.Errorf("submit block: %w", err)
	}

	_, err = bind.WaitMined(context.Background(), a.ethClient, tx)
	if err != nil {
		return fmt.Errorf("wait submit block: %w", err)
	}

	return nil
}
