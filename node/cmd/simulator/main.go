package main

import (
	"flag"
	"node/pkg/simulator"
)

func main() {

	runs := flag.Int("r", 10, "filename of the config file")
	dst := flag.String("d", "./data.csv", "filename of the csv file to store the data")
	ethURL := flag.String("e", "ws://127.0.0.1:8545/", "eth client url")
	contract := flag.String("c", "0x40918Ba7f132E0aCba2CE4de4c4baF9BD2D7D849", "oracle contract address")
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
