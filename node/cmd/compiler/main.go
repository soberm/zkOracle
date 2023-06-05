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
	flag.Parse()

	err := compile(new(zkOracle.AggregationCircuit), path.Join(*b, "aggregation"))
	if err != nil {
		panic(err)
	}

	err = compile(new(zkOracle.SlashingCircuit), path.Join(*b, "slashing"))
	if err != nil {
		panic(err)
	}

}

func compile(circuit frontend.Circuit, dst string) error {

	if _, err := os.Stat(dst); os.IsNotExist(err) {
		err := os.MkdirAll(dst, 0700)
		if err != nil {
			panic(err)
		}
	}

	_r1cs, err := frontend.Compile(ecc.BN254, r1cs.NewBuilder, circuit)
	if err != nil {
		return fmt.Errorf("compile: %w", err)
	}

	file, err := os.Create(path.Join(dst, "r1cs"))
	if err != nil {
		return fmt.Errorf("create r1cs file: %w", err)
	}

	_, err = _r1cs.WriteTo(file)
	if err != nil {
		return fmt.Errorf("write r1cs file: %w", err)
	}

	pk, vk, err := groth16.Setup(_r1cs)
	if err != nil {
		return fmt.Errorf("setup: %w", err)
	}

	file, err = os.Create(path.Join(dst, "pk"))
	if err != nil {
		return fmt.Errorf("create pk file: %w", err)
	}

	_, err = pk.WriteRawTo(file)
	if err != nil {
		return fmt.Errorf("write pk file: %w", err)
	}
	_ = file.Close()

	file, err = os.Create(path.Join(dst, "vk"))
	if err != nil {
		return fmt.Errorf("create vk file: %w", err)
	}
	_, err = vk.WriteRawTo(file)
	if err != nil {
		return fmt.Errorf("write vk file: %w", err)
	}
	_ = file.Close()

	file, err = os.Create(path.Join(dst, "Verifier.sol"))
	if err != nil {
		panic(err)
	}
	err = vk.ExportSolidity(file)
	if err != nil {
		panic(err)
	}

	return nil
}
