package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/consensys/gnark-crypto/accumulator/merkletree"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254/fp"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
	"github.com/consensys/gnark/backend"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/std/accumulator/merkle"
	"github.com/consensys/gnark/std/algebra/twistededwards"
	eddsa2 "github.com/consensys/gnark/std/signature/eddsa"
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
		x := sk.PublicKey.A.X.Bytes()
		y := sk.PublicKey.A.Y.Bytes()
		fmt.Printf("Pub-X: %v\n", big.NewInt(0).SetBytes(x[:]))
		fmt.Printf("Pub-Y: %v\n", big.NewInt(0).SetBytes(y[:]))

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
		result := hexutils.HexToBytes("8a37bed7896a37e676fe5498e7fc14da08897b13147f7181190253c9841e09bb")

		pub.Assign(ecc.BN254, privateKey.PublicKey.Bytes())

		sigBin, err := privateKey.Sign(result, mimc.NewMiMC())
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
		fmt.Printf("Account: %v\n", account)

		votes[i] = zkOracle.ValidatorConstraints{
			Index:             account.Index,
			PublicKey:         pub,
			Balance:           new(big.Int).Set(account.Balance), //passed by reference
			MerkleProof:       proof,
			MerkleProofHelper: helper,
			Signature:         sig,
			BlockHash:         result,
		}

		account.Balance.Add(account.Balance, big.NewInt(5))
		err = state.WriteAccount(account)
		if err != nil {
			return votes, fmt.Errorf("write account: %w", err)
		}
	}

	return votes, nil
}

func MerkleProofs(state *zkOracle.State) (frontend.Variable, [nbAccounts][depth]frontend.Variable, [nbAccounts][depth - 1]frontend.Variable, error) {
	hFunc := mimc.NewMiMC()
	var merkleProofs [nbAccounts][depth]frontend.Variable
	var merkleHelpers [nbAccounts][depth - 1]frontend.Variable
	var merkleRoot frontend.Variable

	for i := 0; i < nbAccounts; i++ {

		var stateBuf bytes.Buffer
		_, err := stateBuf.Write(state.HashData())
		if err != nil {
			return merkleRoot, merkleProofs, merkleHelpers, fmt.Errorf("%v", err)
		}
		root, proof, numLeaves, _ := merkletree.BuildReaderProof(&stateBuf, hFunc, hFunc.Size(), uint64(i))
		proofHelper := merkle.GenerateProofHelper(proof, uint64(i), numLeaves)

		if !merkletree.VerifyProof(hFunc, root, proof, uint64(i), numLeaves) {
			return merkleRoot, merkleProofs, merkleHelpers, errors.New("invalid merkle proof")
		}

		p := make([]*big.Int, len(proof))
		for i, node := range proof {
			p[i] = big.NewInt(0).SetBytes(node)
		}
		fmt.Printf("Proof: %v\n", p)
		fmt.Printf("Helper: %v\n", proofHelper)
		fmt.Printf("Root: %v\n", big.NewInt(0).SetBytes(root))

		var path [depth]frontend.Variable
		var helper [depth - 1]frontend.Variable

		for i := 0; i < len(proof); i++ {
			path[i] = proof[i]
		}

		for i := 0; i < len(proofHelper); i++ {
			helper[i] = proofHelper[i]
		}
		merkleProofs[i] = path
		merkleHelpers[i] = helper
		if i == 0 {
			merkleRoot = root
		}
	}
	return merkleRoot, merkleProofs, merkleHelpers, nil
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

	var circuit zkOracle.Circuit

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

	blockHash := hexutils.HexToBytes("8a37bed7896a37e676fe5498e7fc14da08897b13147f7181190253c9841e09bb")

	var assignment zkOracle.Circuit

	merkleRoot, proof, helper, err := state.MerkleProof(0)
	if err != nil {
		fmt.Printf("merkle proof: %w", err)
		return
	}

	assignment.Root = merkleRoot
	assignment.BlockHash = blockHash

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
	account.Balance = big.NewInt(50)
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

	_ = groth16.Verify(p, vk, pw)

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
