package main

import (
	"flag"
	"fmt"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"node/pkg/zkOracle"
	"os"
	"path"
)

func main() {

	buildPath := flag.String("c", "./build/", "filename of the config file")
	flag.Parse()

	if _, err := os.Stat(*buildPath); os.IsNotExist(err) {
		err := os.MkdirAll(*buildPath, 0700)
		if err != nil {
			panic(err)
		}
	}

	var circuit zkOracle.AggregationCircuit

	_r1cs, err := frontend.Compile(ecc.BN254, r1cs.NewBuilder, &circuit)
	if err != nil {
		panic(err)
	}

	file, err := os.Create(path.Join(*buildPath, "r1cs"))
	if err != nil {
		panic(err)
	}

	_, err = _r1cs.WriteTo(file)
	if err != nil {
		panic(err)
	}

	pk, vk, err := groth16.Setup(_r1cs)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	file, err = os.Create(path.Join(*buildPath, "pk"))
	if err != nil {
		panic(err)
	}

	_, err = pk.WriteRawTo(file)
	if err != nil {
		panic(err)
	}

	file, err = os.Create(path.Join(*buildPath, "vk"))
	if err != nil {
		panic(err)
	}
	_, err = vk.WriteRawTo(file)
	if err != nil {
		panic(err)
	}

	file, err = os.Create(path.Join(*buildPath, "Verifier.sol"))
	if err != nil {
		panic(err)
	}
	err = vk.ExportSolidity(file)
	if err != nil {
		panic(err)
	}

}
