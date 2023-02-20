package main

import (
	"bytes"
	"fmt"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254/fp"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
	"github.com/consensys/gnark/backend"
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
)

const (
	nbAccounts = 4
	depth      = 3
)

func GenerateAccounts() ([]*eddsa.PrivateKey, []*zkOracle.Account, error) {
	privateKeys := make([]*eddsa.PrivateKey, nbAccounts)
	accounts := make([]*zkOracle.Account, nbAccounts)
	for i := 0; i < nbAccounts; i++ {
		r := rand.New(rand.NewSource(int64(i)))
		sk, err := eddsa.GenerateKey(r)
		if err != nil {
			return nil, nil, fmt.Errorf("generate key: %w", err)
		}
		/*		x := sk.PublicKey.A.X.Bytes()
				y := sk.PublicKey.A.Y.Bytes()
				fmt.Printf("Pub-X: %v\n", big.NewInt(0).SetBytes(x[:]))
				fmt.Printf("Pub-Y: %v\n", big.NewInt(0).SetBytes(y[:]))*/

		accounts[i] = &zkOracle.Account{
			big.NewInt(int64(i)),
			&sk.PublicKey,
			big.NewInt(0),
		}
		privateKeys[i] = sk
	}
	return privateKeys, accounts, nil
}

func GenerateVotes(privateKeys []*eddsa.PrivateKey, state *zkOracle.State) ([nbAccounts]zkOracle.ValidatorConstraints, error) {
	var votes [nbAccounts]zkOracle.ValidatorConstraints
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
			return votes, fmt.Errorf("sign: %w", err)
		}
		sig.Assign(ecc.BN254, sigBin)

		_, proof, helper, err := state.MerkleProof(uint64(i))
		if err != nil {
			return votes, fmt.Errorf("merkle proof: %w", err)
		}

		account, err := state.ReadAccount(uint64(i))
		if err != nil {
			return votes, fmt.Errorf("read account: %w", err)
		}

		votes[i] = zkOracle.ValidatorConstraints{
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
			return votes, fmt.Errorf("write account: %w", err)
		}
	}

	return votes, nil
}

func main() {

	privateKeys, accounts, err := GenerateAccounts()
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	state, err := zkOracle.NewState(mimc.NewMiMC(), accounts)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	var circuit zkOracle.AggregationCircuit

	// compile a circuit
	_r1cs, err := frontend.Compile(ecc.BN254, r1cs.NewBuilder, &circuit)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	// setup
	pk, vk, err := groth16.Setup(_r1cs)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	//f, _ := os.Create("./Verifier.sol")
	//vk.ExportSolidity(f)

	blockHash := hexutils.HexToBytes("fc404e20b625e3020de61240b36ab7dba952e662c03214206559c03b004f08f3")

	var assignment zkOracle.AggregationCircuit

	merkleRoot, proof, helper, err := state.MerkleProof(0)
	if err != nil {
		fmt.Printf("merkle proof: %w", err)
		return
	}

	assignment.PreStateRoot = merkleRoot
	assignment.BlockHash = blockHash
	assignment.Request = big.NewInt(0)

	const fpSize = fp.Bytes
	assignment.Aggregator = zkOracle.AggregatorConstraints{
		Index:             0,
		Seed:              twistededwards.Point{X: 0, Y: 1},
		SecretKey:         privateKeys[0].Bytes()[fpSize : 2*fpSize],
		Balance:           big.NewInt(0),
		MerkleProof:       proof,
		MerkleProofHelper: helper,
	}

	account, err := state.ReadAccount(0)
	if err != nil {
		fmt.Printf("read account: %w", err)
		return
	}
	account.Balance.Add(account.Balance, big.NewInt(zkOracle.AggregatorReward))
	err = state.WriteAccount(account)
	if err != nil {
		fmt.Printf("write account: %w", err)
		return
	}

	votes, err := GenerateVotes(privateKeys, state)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	assignment.Validators = votes

	root, err := state.Root()
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	assignment.PostStateRoot = root

	/*	for i := 0; i < nbAccounts; i++ {
		a, _ := state.ReadAccount(uint64(i))
		fmt.Printf("Account: %v\n", a)
	}*/

	w, err := frontend.NewWitness(&assignment, ecc.BN254)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	p, err := groth16.Prove(_r1cs, pk, w, backend.IgnoreSolverError())
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	pw, err := w.Public()
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	err = groth16.Verify(p, vk, pw)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	var (
		a [2]*big.Int
		b [2][2]*big.Int
		c [2]*big.Int
		//input [1]*big.Int
	)

	// get proof bytes
	var buf bytes.Buffer
	p.WriteRawTo(&buf)
	proofBytes := buf.Bytes()

	// proof.Ar, proof.Bs, proof.Krs
	a[0] = new(big.Int).SetBytes(proofBytes[fpSize*0 : fpSize*1])
	a[1] = new(big.Int).SetBytes(proofBytes[fpSize*1 : fpSize*2])
	b[0][0] = new(big.Int).SetBytes(proofBytes[fpSize*2 : fpSize*3])
	b[0][1] = new(big.Int).SetBytes(proofBytes[fpSize*3 : fpSize*4])
	b[1][0] = new(big.Int).SetBytes(proofBytes[fpSize*4 : fpSize*5])
	b[1][1] = new(big.Int).SetBytes(proofBytes[fpSize*5 : fpSize*6])
	c[0] = new(big.Int).SetBytes(proofBytes[fpSize*6 : fpSize*7])
	c[1] = new(big.Int).SetBytes(proofBytes[fpSize*7 : fpSize*8])

	fmt.Printf("A: [%v,%v]\n", a[0], a[1])
	fmt.Printf("B: [[%v,%v],[%v,%v]]\n", b[0][0], b[0][1], b[1][0], b[1][1])
	fmt.Printf("C: [%v,%v]\n", c[0], c[1])
	fmt.Printf("[%v,%v], [[%v,%v],[%v,%v]], [%v,%v]", a[0], a[1], b[0][0], b[0][1], b[1][0], b[1][1], c[0], c[1])
}
