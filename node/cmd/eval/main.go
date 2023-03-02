package main

import (
	"encoding/csv"
	"flag"
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
	"github.com/ethereum/go-ethereum/common"
	"github.com/status-im/keycard-go/hexutils"
	"math/big"
	"math/rand"
	"node/pkg/zkOracle"
	"os"
	"runtime"
	"strconv"
	"time"
)

const (
	nbAccounts = 4
)

func main() {

	runs := flag.Int("r", 10, "filename of the config file")
	csvFileName := flag.String("d", "./data.csv", "filename of the csv file to store the data")
	flag.Parse()

	var circuit zkOracle.AggregationCircuit

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

	f, err := os.Create(*csvFileName)
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

	for i := 0; i < *runs; i++ {
		privateKeys, err := GeneratePrivateKeys(nbAccounts)
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

		assignment, err := AssignVariables(state, privateKeys)
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

func GeneratePrivateKeys(number int) ([]*eddsa.PrivateKey, error) {
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

func CreateAccounts(privateKeys []*eddsa.PrivateKey) ([]*zkOracle.Account, error) {
	accounts := make([]*zkOracle.Account, nbAccounts)
	for i, privateKey := range privateKeys {
		accounts[i] = &zkOracle.Account{
			Index:     big.NewInt(int64(i)),
			PublicKey: &privateKey.PublicKey,
			Balance:   big.NewInt(0),
		}
	}
	return accounts, nil
}

func AssignVariables(state *zkOracle.State, privateKeys []*eddsa.PrivateKey) (*zkOracle.AggregationCircuit, error) {
	var assignment zkOracle.AggregationCircuit

	preStateRoot, err := state.Root()
	if err != nil {
		return nil, fmt.Errorf("pre state root: %w", err)
	}

	aggregatorConstraints, err := AssignAggregatorConstraints(state, privateKeys[2])
	if err != nil {
		return nil, fmt.Errorf("assign aggregator constraints: %w", err)
	}

	assignment.Aggregator = *aggregatorConstraints

	validatorConstraints, validatorBits, err := AssignValidatorConstraints(state, privateKeys)
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
		BlockHash:     hexutils.HexToBytes("fc404e20b625e3020de61240b36ab7dba952e662c03214206559c03b004f08f3"),
		Request:       big.NewInt(0),
		ValidatorBits: validatorBits,
		Aggregator:    *aggregatorConstraints,
		Validators:    validatorConstraints,
	}, nil
}

func AssignAggregatorConstraints(state *zkOracle.State, privateKey *eddsa.PrivateKey) (*zkOracle.AggregatorConstraints, error) {
	var assignment zkOracle.AggregationCircuit

	merkleRoot, proof, helper, err := state.MerkleProof(2)
	if err != nil {
		return nil, fmt.Errorf("merkle proof: %w", err)
	}

	assignment.PreStateRoot = merkleRoot
	assignment.BlockHash = hexutils.HexToBytes("fc404e20b625e3020de61240b36ab7dba952e662c03214206559c03b004f08f3")
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

	account, err := state.ReadAccount(2)
	if err != nil {
		return nil, fmt.Errorf("read account: %w", err)
	}
	account.Balance.Add(account.Balance, big.NewInt(zkOracle.AggregatorReward))
	err = state.WriteAccount(account)
	if err != nil {
		return nil, fmt.Errorf("write account: %w", err)
	}

	return &zkOracle.AggregatorConstraints{
		Index:             2,
		PreSeed:           twistededwards.Point{X: preSeedX, Y: preSeedY},
		PostSeed:          twistededwards.Point{X: postSeedX, Y: postSeedY},
		SecretKey:         sk,
		Balance:           big.NewInt(0),
		MerkleProof:       proof,
		MerkleProofHelper: helper,
	}, nil
}

func AssignValidatorConstraints(state *zkOracle.State, privateKeys []*eddsa.PrivateKey) ([nbAccounts]zkOracle.ValidatorConstraints, *big.Int, error) {
	var validatorConstraints [nbAccounts]zkOracle.ValidatorConstraints
	validatorBits := big.NewInt(0)
	for i, privateKey := range privateKeys {

		var pub eddsa2.PublicKey
		var sig eddsa2.Signature
		result := hexutils.HexToBytes("fc404e20b625e3020de61240b36ab7dba952e662c03214206559c03b004f08f3")

		pub.Assign(ecc.BN254, privateKey.PublicKey.Bytes())

		vote := &zkOracle.Vote{
			Index:     uint64(i),
			Request:   big.NewInt(0),
			BlockHash: common.HexToHash("fc404e20b625e3020de61240b36ab7dba952e662c03214206559c03b004f08f3"),
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

func MemUsage(stop chan struct{}) uint64 {
	var m runtime.MemStats
	var memory uint64
loop:
	for {
		select {
		case <-stop: // triggered when the stop channel is closed
			break loop // exit
		default:
			runtime.ReadMemStats(&m)
			current := bToMb(m.Sys)
			if memory < current {
				memory = current
			}
			time.Sleep(time.Millisecond)
		}
	}
	return memory
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
