package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/csv"
	"fmt"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254/fp"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	edwards "github.com/consensys/gnark-crypto/ecc/bn254/twistededwards"
	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/std/algebra/twistededwards"
	eddsa2 "github.com/consensys/gnark/std/signature/eddsa"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/status-im/keycard-go/hexutils"
	"math/big"
	"math/rand"
	"node/pkg/zkOracle"
	"os"
	"strconv"
	"time"
)

const aggregatorIndex = 1

type AggregationCircuitAnalyzer struct {
	runs       int
	dst        string
	ethClient  *ethclient.Client
	contract   *zkOracle.ZKOracleContract
	privateKey *ecdsa.PrivateKey
	csvWriter  *csv.Writer
}

func (a *AggregationCircuitAnalyzer) Analyze() {
	var circuit zkOracle.AggregationCircuit

	_r1cs, err := frontend.Compile(ecc.BN254, r1cs.NewBuilder, &circuit)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	pkFile, err := os.Open("./build/pk")
	if err != nil {
		panic(err)
	}

	pk := groth16.NewProvingKey(ecc.BN254)
	_, err = pk.ReadFrom(pkFile)
	if err != nil {
		panic(err)
	}

	vkFile, err := os.Open("./build/vk")
	if err != nil {
		panic(err)
	}

	vk := groth16.NewVerifyingKey(ecc.BN254)
	_, err = vk.ReadFrom(vkFile)
	if err != nil {
		panic(err)
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

	a.csvWriter = csv.NewWriter(f)

	headerRow := []string{
		"accounts", "registerCosts", "getBlockByNumberCosts", "submitBlockCosts", "replaceCosts", "exitCosts", "withdrawCosts", "provingTime", "memoryUsage",
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

		registerCosts, err := a.RegisterAccounts(accounts)
		if err != nil {
			panic(err)
		}

		getBlockByNumberCosts, err := a.RequestBlockByNumber()
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

		pw, err := w.Public()
		if err != nil {
			fmt.Printf("%v", err)
			return
		}

		err = groth16.Verify(p, vk, pw)
		if err != nil {
			panic(err)
		}

		submitBlockCosts, err := a.SubmitBlock(p, assignment)
		if err != nil {
			panic(err)
		}

		replaceCosts, err := a.ReplaceAccount(0, state)
		if err != nil {
			panic(err)
		}

		exitCosts, err := a.ExitAccounts(accounts, state)
		if err != nil {
			panic(err)
		}

		withdrawCosts, err := a.WithDrawAccounts(accounts, state)
		if err != nil {
			panic(err)
		}

		data = append(data,
			[]string{
				strconv.Itoa(zkOracle.NumAccounts),
				strconv.Itoa(int(registerCosts)),
				strconv.Itoa(int(getBlockByNumberCosts)),
				strconv.Itoa(int(submitBlockCosts)),
				strconv.Itoa(int(replaceCosts)),
				strconv.Itoa(int(exitCosts)),
				strconv.Itoa(int(withdrawCosts)),
				strconv.Itoa(int(provingTime.Milliseconds())),
				strconv.Itoa(int(<-memoryMeasurement)),
			})
	}

	err = a.csvWriter.WriteAll(data)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	a.csvWriter.Flush()
}

func (a *AggregationCircuitAnalyzer) RegisterAccounts(accounts []*zkOracle.Account) (uint64, error) {
	chainID, err := a.ethClient.ChainID(context.Background())
	if err != nil {
		panic(err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(a.privateKey, chainID)
	if err != nil {
		panic(err)
	}
	auth.GasPrice = big.NewInt(20000000000)

	var averageCost uint64
	for _, account := range accounts {
		tx, err := a.contract.Register(auth, *PublicKeyToZKOraclePublicKey(account.PublicKey), "localhost:25565")
		if err != nil {
			return 0, fmt.Errorf("register: %v", err)
		}
		receipt, err := bind.WaitMined(context.Background(), a.ethClient, tx)
		if err != nil {
			return 0, fmt.Errorf("wait mined: %v", err)
		}
		averageCost += receipt.CumulativeGasUsed
	}
	averageCost = averageCost / uint64(len(accounts))
	return averageCost, nil
}

func (a *AggregationCircuitAnalyzer) RequestBlockByNumber() (uint64, error) {
	chainID, err := a.ethClient.ChainID(context.Background())
	if err != nil {
		panic(err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(a.privateKey, chainID)
	if err != nil {
		panic(err)
	}
	auth.GasPrice = big.NewInt(20000000000)
	auth.Value = big.NewInt(1000000000000000)

	tx, err := a.contract.GetBlockByNumber(auth, big.NewInt(42))
	if err != nil {
		return 0, fmt.Errorf("register: %v", err)
	}

	receipt, err := bind.WaitMined(context.Background(), a.ethClient, tx)
	if err != nil {
		return 0, fmt.Errorf("wait mined: %v", err)
	}

	return receipt.CumulativeGasUsed, nil
}

func (a *AggregationCircuitAnalyzer) ReplaceAccount(i uint64, state *zkOracle.State) (uint64, error) {
	chainID, err := a.ethClient.ChainID(context.Background())
	if err != nil {
		panic(err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(a.privateKey, chainID)
	if err != nil {
		panic(err)
	}
	auth.GasPrice = big.NewInt(20000000000)
	auth.Value = big.NewInt(200000000000)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	sk, err := eddsa.GenerateKey(r)
	if err != nil {
		panic(err)
	}
	newAccount := &zkOracle.Account{
		Index:     big.NewInt(int64(i)),
		PublicKey: &sk.PublicKey,
		Balance:   big.NewInt(200000000000),
	}

	account, err := state.ReadAccount(i)
	if err != nil {
		panic(err)
	}

	_, path, helper, err := state.MerkleProofTest(account.Index.Uint64())
	if err != nil {
		panic(err)
	}

	tx, err := a.contract.Replace(auth, *PublicKeyToZKOraclePublicKey(newAccount.PublicKey), *AccountToZKOracleAccount(&account), path[:], helper[:])
	if err != nil {
		return 0, fmt.Errorf("register: %v", err)
	}

	receipt, err := bind.WaitMined(context.Background(), a.ethClient, tx)
	if err != nil {
		return 0, fmt.Errorf("wait mined: %v", err)
	}

	err = state.WriteAccount(*newAccount)
	if err != nil {
		return 0, fmt.Errorf("wait mined: %v", err)
	}

	return receipt.CumulativeGasUsed, nil
}

func (a *AggregationCircuitAnalyzer) ExitAccounts(accounts []*zkOracle.Account, state *zkOracle.State) (uint64, error) {
	chainID, err := a.ethClient.ChainID(context.Background())
	if err != nil {
		panic(err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(a.privateKey, chainID)
	if err != nil {
		panic(err)
	}
	auth.GasPrice = big.NewInt(20000000000)

	var averageCost uint64
	for i := 0; i < len(accounts); i++ {

		account, err := state.ReadAccount(uint64(i))
		if err != nil {
			panic(err)
		}

		_, path, helper, err := state.MerkleProofTest(account.Index.Uint64())
		if err != nil {
			panic(err)
		}

		tx, err := a.contract.Exit(auth, *AccountToZKOracleAccount(&account), path[:], helper[:])
		if err != nil {
			return 0, fmt.Errorf("register: %v", err)
		}
		receipt, err := bind.WaitMined(context.Background(), a.ethClient, tx)
		if err != nil {
			return 0, fmt.Errorf("wait mined: %v", err)
		}
		averageCost += receipt.CumulativeGasUsed
	}
	averageCost = averageCost / uint64(len(accounts))
	return averageCost, nil
}

func (a *AggregationCircuitAnalyzer) WithDrawAccounts(accounts []*zkOracle.Account, state *zkOracle.State) (uint64, error) {
	chainID, err := a.ethClient.ChainID(context.Background())
	if err != nil {
		panic(err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(a.privateKey, chainID)
	if err != nil {
		panic(err)
	}
	auth.GasPrice = big.NewInt(20000000000)

	var averageCost uint64
	for i := 0; i < len(accounts); i++ {

		account, err := state.ReadAccount(uint64(i))
		if err != nil {
			panic(err)
		}

		_, path, helper, err := state.MerkleProofTest(account.Index.Uint64())
		if err != nil {
			panic(err)
		}

		tx, err := a.contract.Withdraw(auth, *AccountToZKOracleAccount(&account), path[:], helper[:])
		if err != nil {
			return 0, fmt.Errorf("register: %v", err)
		}
		receipt, err := bind.WaitMined(context.Background(), a.ethClient, tx)
		if err != nil {
			return 0, fmt.Errorf("wait mined: %v", err)
		}

		account.Balance = big.NewInt(0)
		err = state.WriteAccount(account)
		if err != nil {
			panic(err)
		}

		averageCost += receipt.CumulativeGasUsed
	}
	averageCost = averageCost / uint64(len(accounts))
	return averageCost, nil
}

func (a *AggregationCircuitAnalyzer) SubmitBlock(p groth16.Proof, assignment *zkOracle.AggregationCircuit) (uint64, error) {
	proof, err := zkOracle.ProofToEthereumProof(p)
	if err != nil {
		panic(err)
	}

	chainID, err := a.ethClient.ChainID(context.Background())
	if err != nil {
		panic(err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(a.privateKey, chainID)
	if err != nil {
		panic(err)
	}
	auth.GasPrice = big.NewInt(20000000000)

	var blockHash [32]byte
	copy(blockHash[2:], hexutils.HexToBytes("4e20b625e3020de61240b36ab7dba952e662c03214206559c03b004f08f3")[:30])

	postStateRoot := assignment.PostStateRoot.([]byte)
	validatorBits := assignment.ValidatorBits.(*big.Int)
	postSeedX := assignment.Aggregator.PostSeed.X.(*big.Int)
	postSeedY := assignment.Aggregator.PostSeed.Y.(*big.Int)

	tx, err := a.contract.SubmitBlock(
		auth,
		new(big.Int).SetUint64(aggregatorIndex),
		big.NewInt(0),
		validatorBits,
		blockHash,
		new(big.Int).SetBytes(postStateRoot),
		postSeedX,
		postSeedY,
		proof.A,
		proof.B,
		proof.C,
	)
	if err != nil {
		panic(err)
	}

	receipt, err := bind.WaitMined(context.Background(), a.ethClient, tx)
	if err != nil {
		panic(err)
	}
	return receipt.CumulativeGasUsed, nil
}

func (a *AggregationCircuitAnalyzer) AssignVariables(state *zkOracle.State, privateKeys []*eddsa.PrivateKey) (*zkOracle.AggregationCircuit, error) {

	preStateRoot, err := state.Root()
	if err != nil {
		return nil, fmt.Errorf("pre state root: %w", err)
	}

	aggregatorConstraints, err := a.AssignAggregatorConstraints(state, privateKeys[aggregatorIndex])
	if err != nil {
		return nil, fmt.Errorf("assign aggregator constraints: %w", err)
	}

	validatorConstraints, validatorBits, err := a.AssignValidatorConstraints(state, privateKeys)
	if err != nil {
		return nil, fmt.Errorf("assign validator constraints: %w", err)
	}

	postStateRoot, err := state.Root()
	if err != nil {
		return nil, fmt.Errorf("pre state root: %w", err)
	}

	return &zkOracle.AggregationCircuit{
		PreStateRoot:  preStateRoot,
		PostStateRoot: postStateRoot,
		BlockHash:     hexutils.HexToBytes("4e20b625e3020de61240b36ab7dba952e662c03214206559c03b004f08f3"),
		Request:       big.NewInt(0),
		ValidatorBits: validatorBits,
		Aggregator:    *aggregatorConstraints,
		Validators:    validatorConstraints,
	}, nil
}

func (a *AggregationCircuitAnalyzer) AssignAggregatorConstraints(state *zkOracle.State, privateKey *eddsa.PrivateKey) (*zkOracle.AggregatorConstraints, error) {
	var assignment zkOracle.AggregationCircuit

	merkleRoot, proof, helper, err := state.MerkleProof(aggregatorIndex)
	if err != nil {
		return nil, fmt.Errorf("merkle proof: %w", err)
	}

	assignment.PreStateRoot = merkleRoot
	assignment.BlockHash = hexutils.HexToBytes("4e20b625e3020de61240b36ab7dba952e662c03214206559c03b004f08f3")
	assignment.Request = big.NewInt(0)

	testX := new(big.Int)
	testY := new(big.Int)

	testX, _ = testX.SetString("5491184307399689246197683245202605692069525215510636283504164930708453453685", 10)
	testY, _ = testY.SetString("2576048849028791939551994783150968389338965397796293068226051430557680319904", 10)

	x := fr.NewElement(0)
	y := fr.NewElement(1)

	preSeed := &edwards.PointAffine{
		X: *x.SetBigInt(testX),
		Y: *y.SetBigInt(testY),
	}

	sk := big.NewInt(0).SetBytes(privateKey.Bytes()[fp.Bytes : 2*fp.Bytes])
	order, _ := new(big.Int).SetString("2736030358979909402780800718157159386076813972158567259200215660948447373041", 10)

	sk.Mod(sk, order)

	preSeedX := new(big.Int)
	preSeedY := new(big.Int)

	preSeed.X.ToBigIntRegular(preSeedX)
	preSeed.Y.ToBigIntRegular(preSeedY)

	var postSeed edwards.PointAffine
	postSeed.ScalarMul(preSeed, sk)

	postSeedX := new(big.Int)
	postSeedY := new(big.Int)

	postSeed.X.ToBigIntRegular(postSeedX)
	postSeed.Y.ToBigIntRegular(postSeedY)

	account, err := state.ReadAccount(aggregatorIndex)
	if err != nil {
		return nil, fmt.Errorf("read account: %w", err)
	}
	account.Balance.Add(account.Balance, big.NewInt(zkOracle.AggregatorReward))
	err = state.WriteAccount(account)
	if err != nil {
		return nil, fmt.Errorf("write account: %w", err)
	}

	return &zkOracle.AggregatorConstraints{
		Index:             aggregatorIndex,
		PreSeed:           twistededwards.Point{X: preSeedX, Y: preSeedY},
		PostSeed:          twistededwards.Point{X: postSeedX, Y: postSeedY},
		SecretKey:         sk,
		Balance:           big.NewInt(0),
		MerkleProof:       proof,
		MerkleProofHelper: helper,
	}, nil
}

func (a *AggregationCircuitAnalyzer) AssignValidatorConstraints(state *zkOracle.State, privateKeys []*eddsa.PrivateKey) ([zkOracle.NumAccounts]zkOracle.ValidatorConstraints, *big.Int, error) {
	var validatorConstraints [zkOracle.NumAccounts]zkOracle.ValidatorConstraints
	validatorBits := big.NewInt(0)
	for i, privateKey := range privateKeys {

		var pub eddsa2.PublicKey
		var sig eddsa2.Signature
		result := hexutils.HexToBytes("4e20b625e3020de61240b36ab7dba952e662c03214206559c03b004f08f3")

		pub.Assign(ecc.BN254, privateKey.PublicKey.Bytes())

		vote := &zkOracle.Vote{
			Index:     uint64(i),
			Request:   big.NewInt(0),
			BlockHash: common.HexToHash("4e20b625e3020de61240b36ab7dba952e662c03214206559c03b004f08f3"),
		}

		hasher := mimc.NewMiMC()
		hasher.Write(vote.Serialize())
		msg := hasher.Sum(nil)

		sigBin, err := privateKey.Sign(msg, mimc.NewMiMC())
		if err != nil {
			return validatorConstraints, validatorBits, fmt.Errorf("sign: %w", err)
		}
		sig.Assign(ecc.BN254, sigBin)

		_, proof, helper, err := state.MerkleProof(uint64(i))
		if err != nil {
			return validatorConstraints, validatorBits, fmt.Errorf("merkle proof: %w", err)
		}

		account, err := state.ReadAccount(uint64(i))
		if err != nil {
			return validatorConstraints, validatorBits, fmt.Errorf("read account: %w", err)
		}

		validatorBit := new(big.Int)
		validatorBit.Exp(big.NewInt(2), account.Index, nil)

		validatorBits = validatorBits.Add(validatorBits, validatorBit)

		validatorConstraints[i] = zkOracle.ValidatorConstraints{
			Index:             account.Index,
			PublicKey:         pub,
			Balance:           new(big.Int).Set(account.Balance), //passed by reference
			MerkleProof:       proof,
			MerkleProofHelper: helper,
			Signature:         sig,
			BlockHash:         result,
		}

		account.Balance.Add(account.Balance, big.NewInt(zkOracle.ValidatorReward))
		err = state.WriteAccount(account)
		if err != nil {
			return validatorConstraints, validatorBits, fmt.Errorf("write account: %w", err)
		}
	}

	return validatorConstraints, validatorBits, nil
}
