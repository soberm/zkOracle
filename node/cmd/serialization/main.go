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

	b := flag.String("b", "./build/", "filename of the config file")
	c := flag.Bool("c", true, "aggregation or slashing circuit")

	flag.Parse()

	if _, err := os.Stat(*b); os.IsNotExist(err) {
		err := os.MkdirAll(*b, 0700)
		if err != nil {
			panic(err)
		}
	}

	var circuit frontend.Circuit

	if *c {
		circuit = &zkOracle.AggregationCircuit{}
	} else {
		circuit = &zkOracle.SlashingCircuit{}
	}

	_r1cs, err := frontend.Compile(ecc.BN254, r1cs.NewBuilder, circuit)
	if err != nil {
		panic(err)
	}

	file, err := os.Create(path.Join(*b, "r1cs"))
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

	file, err = os.Create(path.Join(*b, "pk"))
	if err != nil {
		panic(err)
	}

	_, err = pk.WriteRawTo(file)
	if err != nil {
		panic(err)
	}
	_ = file.Close()

	file, err = os.Create(path.Join(*b, "vk"))
	if err != nil {
		panic(err)
	}
	_, err = vk.WriteRawTo(file)
	if err != nil {
		panic(err)
	}
	_ = file.Close()

	file, err = os.Create(path.Join(*b, "Verifier.sol"))
	if err != nil {
		panic(err)
	}
	err = vk.ExportSolidity(file)
	if err != nil {
		panic(err)
	}

}
