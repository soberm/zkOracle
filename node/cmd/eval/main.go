package main

import (
	"context"
	"encoding/csv"
	"flag"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"math/rand"
	"node/pkg/zkOracle"
	"os"
	"strconv"
	"time"
)

func main() {

	//runs := flag.Int("r", 10, "filename of the config file")
	//csvFileName := flag.String("d", "./data.csv", "filename of the csv file to store the data")
	//circuit := flag.Bool("c", true, "aggregation or slashing circuit")
	flag.Parse()

	csvFile, err := os.OpenFile("data.csv", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}

	csvWriter := csv.NewWriter(csvFile)

	headerRow := []string{
		"ID", "function", "gas",
	}

	err = csvWriter.Write(headerRow)
	if err != nil {
		panic(err)
	}
	csvWriter.Flush()

	ethClient, err := ethclient.Dial("ws://127.0.0.1:8545/")
	if err != nil {
		panic(err)
	}

	contract, err := zkOracle.NewZKOracleContract(common.HexToAddress("0xCf7Ed3AccA5a467e9e704C703E8D87F634fB0Fc9"), ethClient)
	if err != nil {
		panic(err)
	}

	chainID, err := ethClient.ChainID(context.Background())
	if err != nil {
		panic(err)
	}

	ecdsaPrivateKey, err := crypto.HexToECDSA("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	if err != nil {
		panic(err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(ecdsaPrivateKey, chainID)
	if err != nil {
		panic(err)
	}
	auth.GasPrice = big.NewInt(20000000000)

	tx, err := contract.GetBlockByNumber(auth, big.NewInt(0))
	if err != nil {
		panic(err)
	}

	receipt, err := bind.WaitMined(context.Background(), ethClient, tx)
	if err != nil {
		panic(err)
	}

	row := []string{
		"ID", "getBlockByNumber", strconv.FormatInt(int64(receipt.CumulativeGasUsed), 10),
	}

	err = csvWriter.Write(row)
	if err != nil {
		panic(err)
	}
	csvWriter.Flush()

	keys, _ := GeneratePrivateKeys(4)
	accounts, _ := CreateAccounts(keys)
	state, _ := zkOracle.NewState(mimc.NewMiMC(), accounts)

	for i := 0; i < len(keys); i++ {
		x := new(big.Int)
		y := new(big.Int)

		keys[i].PublicKey.A.X.ToBigIntRegular(x)
		keys[i].PublicKey.A.Y.ToBigIntRegular(y)

		tx, err := contract.Register(auth, zkOracle.ZKOraclePublicKey{
			X: x,
			Y: y,
		}, "localhost:25565")
		if err != nil {
			panic(err)
		}

		receipt, err := bind.WaitMined(context.Background(), ethClient, tx)
		if err != nil {
			panic(err)
		}

		row := []string{
			"ID", "register", strconv.FormatInt(int64(receipt.CumulativeGasUsed), 10),
		}

		err = csvWriter.Write(row)
		if err != nil {
			panic(err)
		}
		csvWriter.Flush()
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	sk, err := eddsa.GenerateKey(r)
	if err != nil {
		panic(err)
	}
	account := &zkOracle.Account{
		Index:     big.NewInt(int64(0)),
		PublicKey: &sk.PublicKey,
		Balance:   big.NewInt(0),
	}

	xNew := new(big.Int)
	yNew := new(big.Int)

	account.PublicKey.A.X.ToBigIntRegular(xNew)
	account.PublicKey.A.Y.ToBigIntRegular(yNew)

	_, path, helper, err := state.MerkleProofTest(uint64(0))
	if err != nil {
		panic(err)
	}

	x := new(big.Int)
	y := new(big.Int)

	accounts[0].PublicKey.A.X.ToBigIntRegular(x)
	accounts[0].PublicKey.A.Y.ToBigIntRegular(y)

	tx, err = contract.Replace(auth, zkOracle.ZKOraclePublicKey{
		X: xNew,
		Y: yNew,
	}, zkOracle.ZKOracleAccount{
		Index: accounts[0].Index,
		PubKey: zkOracle.ZKOraclePublicKey{
			X: x,
			Y: y,
		},
		Balance: accounts[0].Balance,
	}, path[:], helper[:])
	if err != nil {
		panic(err)
	}

	receipt, err = bind.WaitMined(context.Background(), ethClient, tx)
	if err != nil {
		panic(err)
	}

	accounts[0] = account
	err = state.WriteAccount(*account)
	if err != nil {
		panic(err)
	}

	row = []string{
		"ID", "replace", strconv.FormatInt(int64(receipt.CumulativeGasUsed), 10),
	}

	err = csvWriter.Write(row)
	if err != nil {
		panic(err)
	}
	csvWriter.Flush()

	for i := 0; i < len(accounts); i++ {

		x := new(big.Int)
		y := new(big.Int)

		accounts[i].PublicKey.A.X.ToBigIntRegular(x)
		accounts[i].PublicKey.A.Y.ToBigIntRegular(y)

		_, path, helper, err := state.MerkleProofTest(uint64(i))
		if err != nil {
			panic(err)
		}

		tx, err = contract.Exit(auth, zkOracle.ZKOracleAccount{
			Index: accounts[i].Index,
			PubKey: zkOracle.ZKOraclePublicKey{
				X: x,
				Y: y,
			},
			Balance: accounts[i].Balance,
		}, path[:], helper[:])
		if err != nil {
			panic(err)
		}

		receipt, err = bind.WaitMined(context.Background(), ethClient, tx)
		if err != nil {
			panic(err)
		}

		row = []string{
			"ID", "exit", strconv.FormatInt(int64(receipt.CumulativeGasUsed), 10),
		}

		err = csvWriter.Write(row)
		if err != nil {
			panic(err)
		}
		csvWriter.Flush()
	}

	for i := 0; i < len(accounts); i++ {

		x := new(big.Int)
		y := new(big.Int)

		accounts[i].PublicKey.A.X.ToBigIntRegular(x)
		accounts[i].PublicKey.A.Y.ToBigIntRegular(y)

		_, path, helper, err := state.MerkleProofTest(uint64(i))
		if err != nil {
			panic(err)
		}

		tx, err = contract.Withdraw(auth, zkOracle.ZKOracleAccount{
			Index: accounts[i].Index,
			PubKey: zkOracle.ZKOraclePublicKey{
				X: x,
				Y: y,
			},
			Balance: accounts[i].Balance,
		}, path[:], helper[:])
		if err != nil {
			panic(err)
		}

		receipt, err = bind.WaitMined(context.Background(), ethClient, tx)
		if err != nil {
			panic(err)
		}

		row = []string{
			"ID", "withdraw", strconv.FormatInt(int64(receipt.CumulativeGasUsed), 10),
		}

		err = csvWriter.Write(row)
		if err != nil {
			panic(err)
		}
		csvWriter.Flush()
	}

	/*if *circuit {
		aggregationCircuitAnalyzer := AggregationCircuitAnalyzer{
			runs: *runs,
			dst:  *csvFileName,
		}
		aggregationCircuitAnalyzer.Analyze()
	} else {
		slashingCircuitAnalyzer := SlashingCircuitAnalyzer{
			runs: *runs,
			dst:  *csvFileName,
		}
		slashingCircuitAnalyzer.Analyze()
	}*/

}
