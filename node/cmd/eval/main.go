package main

import (
	"flag"
)

func main() {

	runs := flag.Int("r", 10, "filename of the config file")
	csvFileName := flag.String("d", "./data.csv", "filename of the csv file to store the data")
	circuit := flag.Bool("c", true, "aggregation or slashing circuit")
	flag.Parse()

	if *circuit {
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
	}

}
