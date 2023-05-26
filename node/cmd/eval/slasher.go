package main

import (
	"encoding/csv"
	"fmt"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	eddsa2 "github.com/consensys/gnark/std/signature/eddsa"
	"github.com/ethereum/go-ethereum/common"
	"github.com/status-im/keycard-go/hexutils"
	"math/big"
	"node/pkg/zkOracle"
	"os"
	"strconv"
	"time"
)

type SlashingCircuitAnalyzer struct {
	runs int
	dst  string
}

func (a *SlashingCircuitAnalyzer) Analyze() {
	var circuit zkOracle.SlashingCircuit

	_r1cs, err := frontend.Compile(ecc.BN254, r1cs.NewBuilder, &circuit)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	pk, vk, err := groth16.Setup(_r1cs)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	f, err := os.Create(a.dst)
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}(f)

	if err != nil {
		panic(err)
	}

	writer := csv.NewWriter(f)

	headerRow := []string{
		"ID", "provingTime", "memoryUsage",
	}

	data := [][]string{
		headerRow,
	}

	for i := 0; i < a.runs; i++ {
		privateKeys, err := GeneratePrivateKeys(zkOracle.NumAccounts)
		if err != nil {
			panic(err)
		}

		accounts, err := CreateAccounts(privateKeys)
		if err != nil {
			panic(err)
		}

		state, err := zkOracle.NewState(mimc.NewMiMC(), accounts)
		if err != nil {
			panic(err)
		}

		assignment, err := a.AssignVariables(state, privateKeys)
		if err != nil {
			panic(err)
		}

		w, err := frontend.NewWitness(assignment, ecc.BN254)
		if err != nil {
			panic(err)
		}

		stop := make(chan struct{})
		memoryMeasurement := make(chan uint64)
		go func() {
			memoryMeasurement <- MemUsage(stop)
		}()

		start := time.Now()
		p, err := groth16.Prove(_r1cs, pk, w)
		if err != nil {
			fmt.Printf("%v", err)
			return
		}

		provingTime := time.Since(start)
		stop <- struct{}{}

		data = append(data, []string{strconv.Itoa(i), strconv.Itoa(int(provingTime.Milliseconds())), strconv.Itoa(int(<-memoryMeasurement))})

		pw, err := w.Public()
		if err != nil {
			fmt.Printf("%v", err)
			return
		}

		err = groth16.Verify(p, vk, pw)
		if err != nil {
			panic(err)
		}
	}

	err = writer.WriteAll(data)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	writer.Flush()
}

func (a *SlashingCircuitAnalyzer) AssignVariables(state *zkOracle.State, privateKeys []*eddsa.PrivateKey) (*zkOracle.SlashingCircuit, error) {

	preStateRoot, err := state.Root()
	if err != nil {
		return nil, fmt.Errorf("pre state root: %w", err)
	}

	validatorConstraints, balance, err := a.AssignSlashedValidatorConstraints(state, privateKeys[0])
	if err != nil {
		return nil, fmt.Errorf("assign validator constraints: %w", err)
	}

	slasherConstraints, err := a.AssignSlasherConstraints(state, privateKeys[2], balance)
	if err != nil {
		return nil, fmt.Errorf("assign validator constraints: %w", err)
	}

	postStateRoot, err := state.Root()
	if err != nil {
		return nil, fmt.Errorf("pre state root: %w", err)
	}

	return &zkOracle.SlashingCircuit{
		PreStateRoot:  preStateRoot,
		PostStateRoot: postStateRoot,
		BlockHash:     hexutils.HexToBytes("fc404e20b625e3020de61240b36ab7dba952e662c03214206559c03b004f08f3"),
		Request:       big.NewInt(0),
		Validator:     validatorConstraints,
		Slasher:       slasherConstraints,
	}, nil
}

func (a *SlashingCircuitAnalyzer) AssignSlasherConstraints(state *zkOracle.State, privateKey *eddsa.PrivateKey, reward *big.Int) (zkOracle.SlasherConstraints, error) {
	var slasherConstraints zkOracle.SlasherConstraints

	var pub eddsa2.PublicKey
	pub.Assign(ecc.BN254, privateKey.PublicKey.Bytes())

	_, proof, helper, err := state.MerkleProof(2)
	if err != nil {
		return slasherConstraints, fmt.Errorf("merkle proof: %w", err)
	}

	account, err := state.ReadAccount(2)
	if err != nil {
		return slasherConstraints, fmt.Errorf("read account: %w", err)
	}
	account.Balance.Add(account.Balance, reward)
	err = state.WriteAccount(account)
	if err != nil {
		return slasherConstraints, fmt.Errorf("write account: %w", err)
	}

	return zkOracle.SlasherConstraints{
		Index:             2,
		PublicKey:         pub,
		Balance:           big.NewInt(0),
		MerkleProof:       proof,
		MerkleProofHelper: helper,
	}, nil
}

func (a *SlashingCircuitAnalyzer) AssignSlashedValidatorConstraints(state *zkOracle.State, privateKey *eddsa.PrivateKey) (zkOracle.SlashedValidatorConstraints, *big.Int, error) {
	var validatorConstraints zkOracle.SlashedValidatorConstraints

	var pub eddsa2.PublicKey
	var sig eddsa2.Signature
	result := hexutils.HexToBytes("ab404e20b625e3020de61240b36ab7dba952e662c03214206559c03b004f08f3")

	pub.Assign(ecc.BN254, privateKey.PublicKey.Bytes())

	vote := &zkOracle.Vote{
		Index:     uint64(0),
		Request:   big.NewInt(0),
		BlockHash: common.HexToHash("ab404e20b625e3020de61240b36ab7dba952e662c03214206559c03b004f08f3"),
	}

	hasher := mimc.NewMiMC()
	hasher.Write(vote.Serialize())
	msg := hasher.Sum(nil)

	sigBin, err := privateKey.Sign(msg, mimc.NewMiMC())
	if err != nil {
		return validatorConstraints, nil, fmt.Errorf("sign: %w", err)
	}
	sig.Assign(ecc.BN254, sigBin)

	_, proof, helper, err := state.MerkleProof(0)
	if err != nil {
		return validatorConstraints, nil, fmt.Errorf("merkle proof: %w", err)
	}

	account, err := state.ReadAccount(0)
	if err != nil {
		return validatorConstraints, nil, fmt.Errorf("read account: %w", err)
	}

	balance := new(big.Int).Set(account.Balance)
	validatorConstraints = zkOracle.SlashedValidatorConstraints{
		Index:             account.Index,
		PublicKey:         pub,
		Balance:           balance, //passed by reference
		MerkleProof:       proof,
		MerkleProofHelper: helper,
		Signature:         sig,
		BlockHash:         result,
	}

	account.Balance.Sub(account.Balance, account.Balance)
	err = state.WriteAccount(account)
	if err != nil {
		return validatorConstraints, nil, fmt.Errorf("write account: %w", err)
	}

	return validatorConstraints, balance, nil
}
