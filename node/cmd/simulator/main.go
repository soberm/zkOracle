package main

import (
	"flag"
	"node/pkg/simulator"
)

func main() {

	runs := flag.Int("r", 10, "filename of the config file")
	dst := flag.String("d", "./data.csv", "filename of the csv file to store the data")
	ethURL := flag.String("e", "ws://127.0.0.1:8545/", "eth client url")
	contract := flag.String("c", "0xDc64a140Aa3E981100a9becA4E685f962f0cF6C90", "oracle contract address")
	privateKey := flag.String("k", "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80", "private key")
	mode := flag.Int("m", 0, "simulation mode")
	flag.Parse()

	aggregationCircuitAnalyzer, err := simulator.NewSimulator(
		*runs,
		*dst,
		simulator.SimulationMode(*mode),
		*ethURL,
		*contract,
		*privateKey,
	)
	if err != nil {
		panic(err)
	}

	err = aggregationCircuitAnalyzer.Simulate()
	if err != nil {
		panic(err)
	}
}
